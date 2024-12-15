package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"scrapper-bot/config"
	"scrapper-bot/service/scrapper"
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

	cfg.
		scr, err := impl.NewScrapper(cfg, &lg)
	if err != nil {
		panic(err)
	}

	stockChannel, err := scr.Scrape(10 * time.Second)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	go func() {
		time.Sleep(1 * time.Minute)
		scr.StopScrape()
		return
	}()

	for stock := range stockChannel {
		fmt.Printf("\n################################\n\tПолучена аномалия\n")
		fmt.Printf("TICKER: 			%s\n", stock.Stock.Ticker)
		fmt.Printf("FIGI: 				%s\n", stock.Stock.FIGI)
		fmt.Printf("UID: 				%s\n", stock.Stock.UID)
		fmt.Printf("Name: 				%s\n", stock.Stock.Name)
		fmt.Printf("\tДанные аномалии:\n")
		fmt.Printf("Цена: 				%f\n", stock.Stock.Price)
		fmt.Printf("Объем: 				%f\n", stock.Volume)
		fmt.Printf("Количество лотов: 	%d\n", stock.NumberLots)
		fmt.Printf("Движение: 			%s\n", stock.StockMove)
		fmt.Printf("Временная метка: 	%s\n", time.Now().String())
		/*fmt.Printf("Итого за день:\n")
		fmt.Printf("Изменение цены: %f\n", stock.PriceChangePerDay)
		fmt.Printf("Покупки: %f %f\n", stock.PurchasesPerDay, stock.PurchasesPerDayVolume)
		fmt.Printf("Продажи: %f %f\n", stock.SalesPerDay, stock.SalesPerDayVolume)*/
	}
}
