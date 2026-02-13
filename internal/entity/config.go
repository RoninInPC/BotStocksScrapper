package entity

import (
	"github.com/tinkoff/invest-api-go-sdk/investgo"
)

// Структура объединяющая все конфигурационные параметры для бота
type Config struct {
	TinkoffApiConfig   investgo.Config `yaml:"tinkoff-parameters"`
	TgToken            string          `yaml:"tg-token"`
	ChatID             int64           `yaml:"chat-id"`
	AdminId            int64           `yaml:"admin-id"`
	RedisChange        RedisDBConfig   `yaml:"redis-change"`
	RedisLog           RedisDBConfig   `yaml:"redis-log"`
	CandleScrapperMode bool            `yaml:"candle-scrapper-mode"`
	CandleDuration     int64           `yaml:"candle-duration-ns"`
	Logger             Logger
}

type RedisDBConfig struct {
	Addr     string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}
