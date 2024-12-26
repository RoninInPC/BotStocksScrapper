package scrapper

import (
	"errors"
	"time"

	"BotStocksScrapper/internal/list"
	repo "BotStocksScrapper/internal/repository/changeBase"
	chb "BotStocksScrapper/internal/repository/changeBase/impl"
	investapi "github.com/tinkoff/invest-api-go-sdk/proto"

	dr "BotStocksScrapper/internal/driver"
	"BotStocksScrapper/internal/entity"
)

type ScrapperTAPI struct {
	config        entity.Config
	driver        *dr.ApiDriver
	trackedStocks []list.StockScrapeInfo
	StockChannel  chan entity.StockInfo
	stopScrapping chan bool
	logger        entity.Logger
	redis         repo.CBRepository
}

// Инициализирует и создает сущность скраппера
// Инициализирует вложенные сущности: драйвер, БД редис
func InitScrapper(config entity.Config) (Scrapper, error) {
	rClient := chb.NewChangeBaseClient(config.RedisChange)
	r := chb.NewChangeBaseRedisRepository(rClient)
	s := ScrapperTAPI{
		StockChannel:  make(chan entity.StockInfo, 100),
		stopScrapping: make(chan bool),
		config:        config,
		trackedStocks: list.GetStocksInfoList(),
		logger:        config.Logger,
		redis:         r,
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
	s.logger.Debug("Инициализация акций прошла успешно")

	tradeStream, err := s.driver.GetTradeCh(stocks)
	if err != nil {
		s.logger.Errorf("ошибка создания стрима обезличенных сделок: %s", err.Error())
		return s.StockChannel, err
	}
	if !tradeStream.IsListen {
		return nil, errors.New("не удалось запустить прослушивание стрима драйвера")
	}
	s.logger.Debug("Драйвер успешно подписался на обновления обезличенных сделок")

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
				if s.skipTime() {
					continue
				}
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
				stockInfo.Stock.Exchange = currentStock.Exchange
				stockInfo.Stock.RealExchange = currentStock.RealExchange
				stockInfo.Volume = totalVolume
				stockInfo.LotsCount = trade.Quantity

				s.logger.Debugf("Получена обезличенная сделка: NAME: %s; TICKER: %s; PRICE: %f; LOT_COUNT: %d; MOVE: %s",
					stockInfo.Stock.Name, stockInfo.Stock.Ticker, stockInfo.Stock.Price, stockInfo.LotsCount, stockInfo.StockMove)

				if totalVolume >= currentStock.AnomalySize {
					stockInfo.IsAnomaly = true
					s.logger.Warnf("Обнаружена аномалия: NAME:%s PRICE: %f ANOMALY SIZE: %f LOT COUNT: %d STOCK MOVE: %s",
						stockInfo.Stock.Name, stockInfo.Stock.Price, stockInfo.Volume, stockInfo.LotsCount, stockInfo.StockMove)
					err = s.driver.GetPerDayStatistics(&stockInfo)
					if err != nil {
						s.logger.Errorf("Ошибка получения доп.информации об акции: %s", err.Error())
						s.logger.Errorf("Информация об акции: NAME:%s PRICE: %f ANOMALY SIZE: %f LOT COUNT: %d STOCK MOVE: %s",
							stockInfo.Stock.Name, stockInfo.Stock.Price, stockInfo.Volume, stockInfo.LotsCount, stockInfo.StockMove)

					} else {
						s.logger.Info("Аномалия успешно обработана")
					}

					stockInfo.PerDaySalesVolume = float64(s.redis.Get(stockInfo.Stock.UID, entity.Sale.String()))
					if stockInfo.PerDaySalesVolume == 0 {
						s.logger.Warn("Аномальное значение суммы продаж: 0 !!!")
					}
					stockInfo.PerDayBuysVolume = float64(s.redis.Get(stockInfo.Stock.UID, entity.Buy.String()))
					if stockInfo.PerDayBuysVolume == 0 {
						s.logger.Warn("Аномальное значение суммы покупок: 0 !!!")
					}

					stockInfo.PerDayVolume = stockInfo.PerDaySalesVolume + stockInfo.PerDayBuysVolume + stockInfo.Volume
					stockInfo.PerDaySalesPercent = (stockInfo.PerDaySalesVolume / stockInfo.PerDayVolume) * float64(100)
					stockInfo.PerDayBuysPercent = float64(100) - stockInfo.PerDaySalesPercent

					stockInfo.VolumeChange = ((float64(100) / (stockInfo.PerDaySalesVolume + stockInfo.PerDayBuysVolume)) * stockInfo.PerDayVolume) - float64(100)
				}

				s.StockChannel <- stockInfo
				ok := s.redis.Add(entity.StockAdd{
					StockName: stockInfo.Stock.UID,
					Type:      stockInfo.StockMove.String(),
					NumPrice:  int64(stockInfo.Volume),
				})
				if !ok {
					s.logger.Errorf("Не удалось добавить в редис запись о сделке")
				}
			}
		}
	}()

	return s.StockChannel, nil
}

// Посылает сигнал для остановки скраппинга
func (s *ScrapperTAPI) StopScrape() {
	s.logger.Info("Скраппером получен сигнал на остановку")
	s.stopScrapping <- true
}

func (s ScrapperTAPI) skipTime() bool {
	return false
	today := time.Now()

	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		s.logger.Errorf("Ошибка установки локации: %s", err.Error())
		return false
	}

	today = today.In(location)

	if today.Weekday() == time.Saturday || today.Weekday() == time.Sunday {
		return true
	} else {
		startTime := time.Date(today.Year(), today.Month(), today.Day(), 9, 50, 0, 0, location)
		endTime := time.Date(today.Year(), today.Month(), today.Day(), 18, 50, 0, 0, location)

		if today.Before(startTime) || today.After(endTime) {
			return true
		}
		return false
	}

}
