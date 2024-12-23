package main

import (
	"time"

	"BotStocksScrapper/app"
)

func main() {
	superGigaUltimateBot, err := app.NewApp()
	if err != nil {
		panic(err)
	}

	err = superGigaUltimateBot.Work()
	if err != nil {
		panic(err)
	}

	time.Sleep(10 * time.Second)
	err = superGigaUltimateBot.Stop()
	if err != nil {
		panic(err)
	}
}
