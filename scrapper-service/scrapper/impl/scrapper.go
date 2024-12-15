package impl

import (
	"scrapper-bot/scrapper-service/entity"
)

type Scrapper interface {
	Scrape() (<-chan entity.StockInfo, error)
	StopScrape()
}
