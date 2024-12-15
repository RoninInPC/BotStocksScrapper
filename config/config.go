package config

import (
	"fmt"
	"os"

	"BotStocksScrapper/entity"
	"gopkg.in/yaml.v3"
)

func LoadConfig(path string) (entity.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return entity.Config{}, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config entity.Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return entity.Config{}, fmt.Errorf("failed to decode YAML: %w", err)
	}

	return config, nil
}
