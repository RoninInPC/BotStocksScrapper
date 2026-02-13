package app

import (
	"BotStocksScrapper/internal/list"
	"BotStocksScrapper/internal/service/application"
)

func App() {
	list.InitFile()
	superGigaUltimateBot, err := application.NewApp()
	if err != nil {
		panic(err)
	}

	superGigaUltimateBot.Work()
}
