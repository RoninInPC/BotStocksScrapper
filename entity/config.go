package entity

import (
	"github.com/tinkoff/invest-api-go-sdk/investgo"
)

// Структура объединяющая все конфигурационные параметры для бота
type Config struct {
	TinkoffApiConfig investgo.Config `yaml:"tinkoff-parameters"`
	TgToken          string          `yaml:"tg-token"`
	RedisCfg         RedisConfig     `yaml:"redis"`
	Logger           Logger
}

// Структура конфига для БД Redis
type RedisConfig struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	DbNumber int    `json:"db_number"`
}
