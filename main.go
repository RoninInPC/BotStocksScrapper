package main

import (
	"BotStocksScrapper/app"
)

func main() {
	superGigaUltimateBot, err := app.NewApp()
	if err != nil {
		panic(err)
	}

	superGigaUltimateBot.Work()
}
