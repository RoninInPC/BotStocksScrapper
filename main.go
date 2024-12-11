package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"scrapper-bot/config"
	"scrapper-bot/scrapper"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	file, err := os.Open("logs.txt")
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
	scr, err := scrapper.NewScrapper(cfg, &lg)
	if err != nil {
		panic(err)
	}

	stockChannel := scr.Scrape(10 * time.Second)

	for stock := range stockChannel {
		fmt.Printf("\n################################\n\tПолучена аномалия\n")
		fmt.Printf("NAME: %s\n", stock.Stock.Name)
		fmt.Printf("TICKER: %s\n", stock.Stock.Tag)
		fmt.Printf("\tДанные аномалии:\n")
		fmt.Printf("Цена: %f\n", stock.Price)
		fmt.Printf("Объем: %f\n", stock.Volume)
		fmt.Printf("Изменение на объеме: %s\n", stock.VolumeChange)
		fmt.Printf("Движение: %s\n", stock.StockMove)
		fmt.Printf("Итого за день:\n")
		fmt.Printf("Изменение цены: %f\n", stock.PriceChangePerDay)
		fmt.Printf("Покупки: %f %f\n", stock.PurchasesPerDay, stock.PurchasesPerDayVolume)
		fmt.Printf("Продажи: %f %f\n", stock.SalesPerDay, stock.SalesPerDayVolume)
	}
}
