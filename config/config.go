package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Mention struct {
	All     bool     `yaml:"all"`
	Emails  []string `yaml:"emails"`
	OpenIDs []string `yaml:"open_ids"`
}

type Bot struct {
	// Bot Webhook URL
	Webhook string   `yaml:"url"`
	Mention *Mention `yaml:"mention"`
}

type App struct {
	ID     string `yaml:"id"`
	Secret string `yaml:"secret"`
}

type Config struct {
	Bots map[string]*Bot `yaml:"bots"`
	App  *App            `yaml:"app"`
}

func Load(filename string) (*Config, error) {
	var conf Config
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bs, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
