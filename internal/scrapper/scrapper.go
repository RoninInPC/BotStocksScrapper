package scrapper

import (
	"BotStocksScrapper/internal/entity"
)

type Scrapper interface {
	Scrape() (<-chan entity.StockInfo, error)
	StopScrape()
}
