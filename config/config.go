package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
	"scrapper-bot/entity"
)

type Config struct {
	BotToken     string                `yaml:"telegram"`
	TinkoffToken string                `yaml:"tinkoff-api"`
	StocksList   []entity.TrackedStock `yaml:"tracked-stocks"`
}

func LoadConfig() (*Config, error) {
	// Открываем файл
	file, err := os.Open("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	// Декодируем YAML в структуру Config
	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %w", err)
	}

	return &config, nil
}
