package scrapper

import (
	"time"

	"scrapper-bot/entity"
)

type Scrapper interface {
	Scrape(sleepTimeMs time.Duration) <-chan entity.StockInfo
}
