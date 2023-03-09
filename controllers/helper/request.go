package helper

import (
	"errors"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

var test int

type Context struct {
	C          *gin.Context
	CreateOnce sync.Once
}

type RespInfo interface{}

func (c *Context) GetQuery(key string) string {
	return GetQuery(c.C, key)
}

func (c *Context) GetQueryInt64(key string) (int64, error) {
	return GetQueryInt64(c.C, key)
}

func (c *Context) GetPostFormUint64(key string) (uint64, error) {
	return GetPostFormUint64(c.C, key)
}

func GetQuery(c *gin.Context, key string) string {
	val, _ := c.GetQuery(key)
	return val
}

// GetQueryInt64 It returns the keyed url query value with uint64 type
func GetQueryInt64(c *gin.Context, key string) (int64, error) {
	val, ok := c.GetQuery(key)
	if !ok {
		return 0, errors.New("get params err: key")
	}
	return strconv.ParseInt(val, 10, 64)
}

// GetPostFormUint64 It returns the specified key from a POST urlencoded form or multipart with uint64 type
func GetPostFormUint64(c *gin.Context, key string) (uint64, error) {
	val, ok := c.GetPostForm(key)
	if !ok {
		return 0, errors.New("get params err: key")
	}
	return strconv.ParseUint(val, 10, 64)
}
