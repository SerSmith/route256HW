package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const pathToConfig = "config.yaml"

type Config struct {
	Telegram struct {
		Token            string `yaml:"token"`
		Reciever_chat_id string `yaml:"reciever_chat_id"`
	} `yaml:"telegram"`
	Kafka struct {
		TopicStatus string   `yaml:"topicStatus"`
		GroupName   string   `yaml:"groupName"`
		Brokers     []string `yaml:"brokers"`
	} `yaml:"kafka"`
}

var AppConfig = Config{}

func Init() error {
	rawYaml, err := os.ReadFile(pathToConfig)
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	err = yaml.Unmarshal(rawYaml, &AppConfig)
	if err != nil {
		return fmt.Errorf("parse config file: %w", err)
	}

	return nil
}
