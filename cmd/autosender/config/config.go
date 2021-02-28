package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	Message struct {
		ServerID  string `yaml:"server_id"`
		ChannelID string `yaml:"channel_id"`
		Text      string `yaml:"text"`
	}
	Settings struct {
		MinRandMinute int `yaml:"min_rand_minute"`
		MaxRandMinute int `yaml:"max_rand_minute"`
	}
	Planned []int `yaml:"planned"`

	IsDebug bool `yaml:"debug"`
}

func (c *Conf) GetConf() *Conf {
	file, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Printf("Cant get config #%v", err)
	}

	err = yaml.Unmarshal(file, c)

	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
