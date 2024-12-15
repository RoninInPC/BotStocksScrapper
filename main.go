package main

import (
	"io"
	"os"
	"time"

	"BotStocksScrapper/config"
	"BotStocksScrapper/service/scrapper"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.LoadConfig("./BotStocksScrapper/config/config.yaml")
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	writer := io.MultiWriter(os.Stdout, file)
	lg := logrus.Logger{}
	lg.SetOutput(writer)
	lg.SetLevel(logrus.DebugLevel)
	lg.SetReportCaller(true)
	lg.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})

	cfg.Logger = &lg

	// Указать вторым параметром клиент телеги, третьим ID чата
	scrservice, err := scrapper.NewScrapperService(cfg, nil, 0)
	if err != nil {
		panic(err)
	}

	go scrservice.Scrap()

	time.Sleep(1 * time.Minute)
	scrservice.Stop()

}
