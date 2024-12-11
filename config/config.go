package config

import (
	"fmt"
	"os"

	"github.com/tinkoff/invest-api-go-sdk/investgo"
	"gopkg.in/yaml.v3"
	"scrapper-bot/entity"
)

type Config struct {
	TinkoffApiConfig investgo.Config       `yaml:"tinkoff-parameters"`
	StocksList       []entity.TrackedStock `yaml:"tracked-stocks"`
}

func LoadConfig() (Config, error) {
	file, err := os.Open("config.yaml")
	if err != nil {
		return Config{}, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return Config{}, fmt.Errorf("failed to decode YAML: %w", err)
	}

	return config, nil
}
