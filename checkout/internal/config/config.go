package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const pathToConfig = "config.yaml"

type Config struct {
	Token    string `yaml:"token"`
	Services struct {
		Loms           string `yaml:"loms"`
		ProductService string `yaml:"productservice"`
	} `yaml:"services"`

	DB struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Server   string `yaml:"server"`
		Name     string `yaml:"name"`
	} `yaml:"DB"`
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
