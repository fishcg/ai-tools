package service

import (
	"github.com/jinzhu/gorm"

	"github.com/fish/ai-tools/config"
	"github.com/fish/ai-tools/db"
	"github.com/fish/ai-tools/service/openai"
)

var (
	DB     *gorm.DB
	OpenAI *openai.Client
)

func Init(conf *config.Config) (err error) {
	DB, err = db.InitDatabase((*db.Config)(&conf.DB))
	if err != nil {
		return err
	}

	OpenAI = openai.NewClient(conf.OpenAI)
	if err != nil {
		return err
	}
	return nil
}
