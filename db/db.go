package db

import (
	"database/sql"
	_ "database/sql" // used by gorm
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // used by gorm
	"github.com/jinzhu/gorm"
)

// Drivers
const (
	DriverMysql  = "mysql"
	DriverSqlite = "sqlite"

	defaultDriver = DriverMysql
)

// Driver 当前使用的 driver
var Driver = defaultDriver

var errParams = "Incorrect parameter count"

// ErrNoRowsAffected is the error of did not affect any rows in one operation.
var ErrNoRowsAffected = errors.New("no rows affected")

const (
	// https://github.com/go-sql-driver/mysql#charset 对于旧版本 MySQL 会自动使用 utf8
	defaultCharset   = "utf8mb4,utf8"
	defaultParseTime = true
	defaultLoc       = "Local"
)

// IsErrNoRows 是否是未查询到数据
func IsErrNoRows(err error) bool {
	return err == sql.ErrNoRows || gorm.IsRecordNotFoundError(err) //nolint:bannedfunc
}

// Config for database config
type Config struct {
	URL          string `yaml:"url"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Name         string `yaml:"name"`
	User         string `yaml:"user"`
	Pass         string `yaml:"pass"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxLifeTime  string `yaml:"max_life_time"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	Debug        bool   `yaml:"debug"`
}

// InitDatabase init database connection
func InitDatabase(conf *Config) (*gorm.DB, error) {
	maxLifeTime, err := time.ParseDuration(conf.MaxLifeTime)
	if err != nil {
		return nil, err
	}
	driver, dsn, err := conf.getDriverAndDSN()
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	db = db.BlockGlobalUpdate(true)
	if conf.Debug {
		db = db.Debug()
	}

	sqldb := db.DB()
	err = sqldb.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	sqldb.SetMaxIdleConns(conf.MaxIdleConns)
	// Fix db invalid connection after EOF
	sqldb.SetConnMaxLifetime(maxLifeTime)
	if conf.MaxOpenConns != 0 {
		sqldb.SetMaxOpenConns(conf.MaxOpenConns)
	}

	// connect success
	return db, nil
}

// getDriverAndDsn 获取驱动及 DSN
// https://github.com/go-sql-driver/mysql
func (conf *Config) getDriverAndDSN() (driver, dsn string, err error) {
	if conf.URL == "" {
		if conf.Host == "" || conf.Port == 0 || conf.Name == "" || conf.User == "" || conf.Pass == "" {
			err = errors.New("db config should not be empty")
			return
		}

		driver = defaultDriver
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s", conf.User, conf.Pass, conf.Host,
			conf.Port, conf.Name, defaultCharset, defaultParseTime, defaultLoc)
	} else {
		items := strings.SplitN(conf.URL, "://", 2)
		if len(items) != 2 {
			err = errors.New("db url should be configured such as 'mysql://user:password@tcp(host:port)/dbname?param1=value1'")
			return
		}
		driver, dsn = items[0], items[1]
		if driver != defaultDriver {
			err = errors.New("db driver only support mysql currently")
			return
		}
		if dsn == "" {
			err = errors.New("db dsn is empty")
		}
	}

	return
}

// TODO: Expr 方法返回 gorm.Expr() 的值

// IFExpr 实现 if else 的情况
// NOTICE: 有一定风险，不要传入过于复杂的情况
// mysql: IF(condition, valueTrue, valueFalse)
// sqlite: IIF(condition, valueTrue, valueFalse)
// 通用: CASE WHEN condition THEN valueTrue ELSE valueFalse END
func IFExpr(condition, valueTrue, valueFalse string) (query string) {
	switch Driver {
	case DriverSqlite:
		query = fmt.Sprintf("IIF(%s, %s, %s)", condition, valueTrue, valueFalse)
	case DriverMysql:
		query = fmt.Sprintf("IF(%s, %s, %s)", condition, valueTrue, valueFalse)
	default:
		query = fmt.Sprintf("CASE WHEN (%s) THEN (%s) ELSE (%s) END", condition, valueTrue, valueFalse)
	}
	return
}

// GreatestExpr 获得最大值
// 注意如果对于 mysql，如果某个表达式的值超过数据库中定义的范围，不要用 GreatestExpr
// mysql: GREATEST(value[0], value[1]...)
// sqlite: MAX(value[0], value[1]...)
func GreatestExpr(value ...string) (query string) {
	if len(value) < 2 {
		panic(errParams)
	}
	switch Driver {
	case DriverSqlite:
		query = fmt.Sprintf("MAX(%s)", strings.Join(value, ", "))
	case DriverMysql:
		fallthrough
	default:
		query = fmt.Sprintf("GREATEST(%s)", strings.Join(value, ", "))
	}
	return
}

// LeastExpr 获得最小值
// 注意如果对于 mysql，如果某个表达式的值超过数据库中定义的范围，不要用 LeastExpr，比如 LeastExpr("vitality - 20", 0)，数据库 vitality 定义为非负，而 vitality - 20 小于 0
// mysql: LEAST(value[0], value[1]...)
// sqlite: MIN(value[0], value[1]...)
func LeastExpr(value ...string) (query string) {
	if len(value) < 2 {
		panic(errParams)
	}
	switch Driver {
	case DriverSqlite:
		query = fmt.Sprintf("MIN(%s)", strings.Join(value, ", "))
	case DriverMysql:
		fallthrough
	default:
		query = fmt.Sprintf("LEAST(%s)", strings.Join(value, ", "))
	}
	return
}

// UpdateOnDuplicateKeyExpr 插入数据时如果有字段冲突时进行的更新操作
// 只有在 driver 为 sqlite 的情况下需要传 duplicateKey，duplicateKey 为冲突的字段
// mysql: ON DUPLICATE KEY UPDATE updates
// sqlite: ON CONFLICT(duplicateKeys) DO UPDATE SET updates
func UpdateOnDuplicateKeyExpr(updates string, duplicateKeys string) (query string) {
	switch Driver {
	case DriverSqlite:
		query = fmt.Sprintf("ON CONFLICT(%s) DO UPDATE SET %s", duplicateKeys, updates)
	case DriverMysql:
		fallthrough
	default:
		query = fmt.Sprintf("ON DUPLICATE KEY UPDATE %s", updates)
	}
	return
}

// SubSatExpr 饱和减法
// 参数为被减数的字段名称和减数，最后的差值结果不会为负数
func SubSatExpr(columnName string, subtrahend int) interface{} {
	return gorm.Expr(fmt.Sprintf("%s - %d", GreatestExpr(columnName, strconv.Itoa(subtrahend)), subtrahend))
}
