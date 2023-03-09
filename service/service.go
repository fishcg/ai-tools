package service

import (
	"github.com/jinzhu/gorm"

	"github.com/fish/ai-tools/config"
	"github.com/fish/ai-tools/db"
	"github.com/fish/ai-tools/service/gpt"
)

var (
	DB  *gorm.DB
	Gpt *gpt.Client
)

func Init(conf *config.Config) (err error) {
	DB, err = db.InitDatabase((*db.Config)(&conf.DB))
	if err != nil {
		return err
	}

	Gpt = gpt.NewClient(conf.GPT)
	if err != nil {
		return err
	}
	return nil
}
