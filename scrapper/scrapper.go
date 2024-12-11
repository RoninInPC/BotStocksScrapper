package scrapper

import (
	"time"

	"BotStocksScrapper/entity"
)

type Scrapper interface {
	Scrape(sleepTimeMs time.Duration) <-chan entity.StockInfo
}
