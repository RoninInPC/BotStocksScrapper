package scrapper

import (
	"BotStocksScrapper/entity"
)

type Scrapper interface {
	Scrape() (<-chan entity.StockInfo, error)
	StopScrape()
}
