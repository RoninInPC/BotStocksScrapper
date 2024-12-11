package scrapper

import (
	"time"

	"github.com/sirupsen/logrus"
	"scrapper-bot/config"
	dr "scrapper-bot/driver"
	"scrapper-bot/entity"
)

type scrapper struct {
	config        config.Config
	driver        *dr.ApiDriver
	trackedStocks []entity.TrackedStock
	StockChannel  chan entity.StockInfo
	stopScrapping chan bool
	logger        *logrus.Logger
}

func NewScrapper(config config.Config, lg *logrus.Logger) (Scrapper, error) {
	s := scrapper{
		StockChannel:  make(chan entity.StockInfo, 100),
		stopScrapping: make(chan bool),
		config:        config,
		trackedStocks: config.StocksList,
		logger:        lg,
	}

	var err error
	s.driver, err = dr.NewApiDriver(config.TinkoffApiConfig, lg)
	if err != nil {
		s.logger.Errorf("не удалось создать драйвер tinkoff api: %s", err.Error())
		return scrapper{}, err
	}

	return s, nil
}

// Запускает горутину скраппера.
// Возвращает канал в который приходят отловленные аномалии
func (s scrapper) Scrape(sleepTimeMs time.Duration) <-chan entity.StockInfo {
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
			time.Sleep(time.Millisecond * sleepTimeMs)
		}
	}()

	return s.StockChannel
}

// Посылает сигнал для остановки скраппинга
func (s *scrapper) StopScrape() {
	s.stopScrapping <- true
}
