package main

import (
	"io"
	"os"
	"time"

	"BotStocksScrapper/config"
	"BotStocksScrapper/entity"
	"BotStocksScrapper/service/scrapper"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.LoadConfig("./config/config.yaml")
	if err != nil {
		panic(err)
	}

	// Инициализируем логер
	file, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	writer := io.MultiWriter(os.Stdout, file)
	cfg.Logger = entity.NewLogger(writer, logrus.DebugLevel)

	// Запускаем сервис скраппера
	// Указать вторым параметром клиент телеги, третьим ID чата
	scrservice, err := scrapper.NewScrapperService(cfg, nil, 0)
	if err != nil {
		panic(err)
	}

	go scrservice.Scrap()

	time.Sleep(10 * time.Second)
	scrservice.Stop()
}
