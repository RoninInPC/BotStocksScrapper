package app

import "BotStocksScrapper/internal/service/application"

func App() {
	superGigaUltimateBot, err := application.NewApp()
	if err != nil {
		panic(err)
	}

	superGigaUltimateBot.Work()
}
