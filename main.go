package main

import (
	"time"

	"BotStocksScrapper/config"
	"BotStocksScrapper/entity"
	"BotStocksScrapper/service/scrapper"
)

func main() {
	cfg, err := config.LoadConfig("./config/config.yaml")
	if err != nil {
		panic(err)
	}

	// Инициализируем логер
	cfg.Logger = entity.NewLogger()

	// Запускаем сервис скраппера
	// Указать вторым параметром клиент телеги, третьим ID чата
	scrservice, err := scrapper.NewScrapperService(cfg, nil, 0)
	if err != nil {
		panic(err)
	}

	go scrservice.Scrap()

	time.Sleep(1 * time.Minute)
	scrservice.Stop()

}
