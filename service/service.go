package service

import (
	"gorm.io/gorm"

	"github.com/fish/ai-tools/config"
	_ "github.com/fish/ai-tools/db"
	"github.com/fish/ai-tools/service/openai"
)

var (
	DB     *gorm.DB
	OpenAI *openai.Client
)

func Init(conf *config.Config) (err error) {
	OpenAI = openai.NewClient(conf.OpenAI)
	if err != nil {
		return err
	}

	// TODO: 添加 DB 用于存储调整 prompt
	// DB, err = db.InitDatabase(&conf.DB)
	// if err != nil {
	//     return err
	// }

	return nil
}
