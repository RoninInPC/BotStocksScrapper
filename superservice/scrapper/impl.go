package scrapper

import (
	"time"

	"scrapper-bot/config"
	"scrapper-bot/entity"
	dr "scrapper-bot/superservice/scrapper/apiDriver"
)

type scrapper struct {
	config        config.Config
	driver        *dr.ApiDriver
	trackedStocks []entity.TrackedStock
	StockChannel  chan entity.StockInfo
	stopScrapping chan bool
}

func NewScrapper(config config.Config) Scrapper {
	s := scrapper{
		StockChannel:  make(chan entity.StockInfo, 100),
		stopScrapping: make(chan bool),
		config:        config,
		driver:        dr.NewApiDriver(config.TinkoffToken),
		trackedStocks: config.StocksList,
	}

	return s
}

// Запускает горутину скраппера.
// Возвращает канал в который приходят отловленные аномалии
func (s scrapper) Scrape(sleepTimeMs time.Duration) <-chan entity.StockInfo {
	go func() {
		for {
			select {
			case <-s.stopScrapping:
				return

			default:
				for _, stock := range s.trackedStocks {
					stockInfo := s.driver.GetStockInfo(stock.StockTag)
					if stockInfo != nil {
						s.StockChannel <- *stockInfo
					}
				}
			}
			time.Sleep(time.Millisecond * sleepTimeMs)
		}
	}()

	return s.StockChannel
}

func (s *scrapper) StopScrape() {
	s.stopScrapping <- true
}
