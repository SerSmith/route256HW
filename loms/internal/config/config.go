package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const pathToConfig = "config.yaml"

type Config struct {
	DB struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Server   string `yaml:"server"`
		Name     string `yaml:"name"`
	} `yaml:"DB"`
	Kafka struct {
		Brokers     []string `yaml:"brokers"`
		TopicStatus string   `yaml:"topicStatus"`
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

func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		c.DB.User,
		c.DB.Password,
		c.DB.Server,
		c.DB.Name)
}
