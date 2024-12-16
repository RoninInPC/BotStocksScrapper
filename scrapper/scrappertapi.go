package scrapper

import (
	investapi "github.com/tinkoff/invest-api-go-sdk/proto"

	dr "BotStocksScrapper/driver"
	"BotStocksScrapper/entity"
)

type ScrapperTAPI struct {
	config        entity.Config
	driver        *dr.ApiDriver
	trackedStocks []entity.TrackedStock
	StockChannel  chan entity.StockInfo
	stopScrapping chan bool
	logger        entity.Logger
}

func InitScrapper(config entity.Config) (Scrapper, error) {
	s := ScrapperTAPI{
		StockChannel:  make(chan entity.StockInfo, 100),
		stopScrapping: make(chan bool),
		config:        config,
		trackedStocks: entity.GetStocksInfoList(),
		logger:        config.Logger,
	}

	var err error
	s.driver, err = dr.NewApiDriver(config.TinkoffApiConfig, config.Logger)
	if err != nil {
		s.logger.Errorf("не удалось создать драйвер tinkoff api: %s", err.Error())
		return &ScrapperTAPI{}, err
	}

	return &s, nil
}

// Запускает горутину скраппера.
// Возвращает канал в который приходят отловленные аномалии
func (s *ScrapperTAPI) Scrape() (<-chan entity.StockInfo, error) {

	stocks, err := s.driver.InitStocks(s.trackedStocks)
	if err != nil {
		s.logger.Errorf("Не удалось инициализировать акции: %s", err.Error())
		return s.StockChannel, err
	}

	tradeStream, err := s.driver.GetTradeCh(stocks)
	if err != nil {
		s.logger.Errorf("ошибка создания стрима обезличенных сделок: %s", err.Error())
		return s.StockChannel, err
	}

	go func() {
		for {
			select {
			case <-s.stopScrapping:
				close(s.StockChannel)
				close(s.stopScrapping)
				tradeStream.Stream.Stop()
				s.logger.Infof("Скраппер остановлен")
				return

			case trade := <-tradeStream.Channel:
				var currentStock entity.Stock
				for _, stock := range stocks {
					if stock.FIGI == trade.Figi {
						currentStock = stock
					}
				}
				stockInfo := entity.StockInfo{}

				stockInfo.Stock.Price = float64(trade.Price.GetUnits()) + float64(trade.Price.Nano)/1e9
				totalVolume := stockInfo.Stock.Price * float64(trade.Quantity) * float64(currentStock.MinLotCount)

				if trade.Direction == investapi.TradeDirection_TRADE_DIRECTION_BUY {
					stockInfo.StockMove = entity.Buy
				} else {
					stockInfo.StockMove = entity.Sale
				}
				stockInfo.Stock.Ticker = currentStock.Ticker
				stockInfo.Stock.FIGI = trade.Figi
				stockInfo.Stock.UID = trade.GetInstrumentUid()
				stockInfo.Stock.Price = float64(trade.Price.GetUnits()) + float64(trade.Price.Nano)/1e9
				stockInfo.Stock.Name = currentStock.Name
				stockInfo.Stock.MinLotCount = currentStock.MinLotCount
				stockInfo.Volume = totalVolume
				stockInfo.NumberLots = trade.Quantity

				if totalVolume >= currentStock.AnomalySize {
					stockInfo.IsAnomaly = true
					s.logger.Warnf("Обнаружена аномалия: NAME:%s PRICE: %f ANOMALY SIZE: %f STOCK MOVE: %s\n",
						stockInfo.Stock.Name, stockInfo.Stock.Price, stockInfo.Volume, stockInfo.StockMove)
				}

				s.StockChannel <- stockInfo
			}
		}
	}()

	return s.StockChannel, nil
}

// Посылает сигнал для остановки скраппинга
func (s *ScrapperTAPI) StopScrape() {
	s.stopScrapping <- true
}
