package scrapper

import (
	"github.com/sirupsen/logrus"
	"time"

	"BotStocksScrapper/config"
	dr "BotStocksScrapper/driver"
	"BotStocksScrapper/entity"
)

type ScrapperTAPI struct {
	config        config.Config
	driver        *dr.ApiDriver
	trackedStocks []entity.TrackedStock
	StockChannel  chan entity.StockInfo
	stopScrapping chan bool
	logger        *logrus.Logger
}

func InitScrapper(config config.Config, lg *logrus.Logger) (ScrapperTAPI, error) {
	s := ScrapperTAPI{
		StockChannel:  make(chan entity.StockInfo, 100),
		stopScrapping: make(chan bool),
		config:        config,
		trackedStocks: config.StocksList,
	}

	var err error
	s.driver, err = dr.NewApiDriver(config.TinkoffApiConfig, lg)
	if err != nil {
		s.logger.Errorf("не удалось создать драйвер tinkoff api: %s", err.Error())
		return ScrapperTAPI{}, err
	}

	return s, nil
}

// Запускает горутину скраппера.
// Возвращает канал в который приходят отловленные аномалии
func (s ScrapperTAPI) Scrape(sleepTime time.Duration) <-chan entity.StockInfo {
	go func() {
		for {
			select {
			case <-s.stopScrapping:
				close(s.StockChannel)
				close(s.stopScrapping)
				return

			default:
				for _, stock := range s.trackedStocks {
					stockInfo, err := s.driver.GetStockInfo(stock.StockTag)
					if err != nil {
						s.logger.Errorf("Error getting stock info: %s", err.Error())
					} else {
						s.StockChannel <- stockInfo
					}
				}
			}
			time.Sleep(sleepTime)
		}
	}()

	return s.StockChannel
}

// Посылает сигнал для остановки скраппинга
func (s *ScrapperTAPI) StopScrape() {
	s.stopScrapping <- true
}
