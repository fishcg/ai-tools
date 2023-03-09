package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/fish/ai-tools/db"
	"github.com/fish/ai-tools/service/openai"
)

// Conf 全局配置文件
var Conf *Config

// Config config structure
type Config struct {
	OpenAI *openai.Config `yaml:"openai"`
	HTTP   SectionHTTP    `yaml:"http"`
	DB     db.Config      `yaml:"db"`
}

// SectionHTTP .
type SectionHTTP struct {
	Enabled bool   `yaml:"enabled"`
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
	Mode    string `yaml:"mode"`
}

// LoadFromYML load config from yml file
func LoadFromYML(configPath string) error {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, &Conf)
	if err != nil {
		return err
	}
	return err
}
