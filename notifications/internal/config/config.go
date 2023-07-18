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
	DB struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Server   string `yaml:"server"`
		Name     string `yaml:"name"`
	} `yaml:"DB"`
	Redis struct {
		Host string `yaml:"host"`
		Pass string `yaml:"pass"`
		DB   int    `yaml:"int"`
	} `yaml:"Redis"`
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

func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		c.DB.User,
		c.DB.Password,
		c.DB.Server,
		c.DB.Name)
}
