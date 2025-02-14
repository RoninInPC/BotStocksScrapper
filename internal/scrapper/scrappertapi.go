package scrapper

import (
	"BotStocksScrapper/internal/analysis"
	"BotStocksScrapper/internal/analysis/simple"
	"BotStocksScrapper/internal/mappers"
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
	analysis      analysis.Analysis
	candleMod     bool
}

// Инициализирует и создает сущность скраппера
// Инициализирует вложенные сущности: драйвер, БД редис
func InitScrapper(config entity.Config) (Scrapper, error) {
	rClient := chb.NewChangeBaseClient(config.RedisChange)
	if rClient == nil {
		config.Logger.Panicf("Не удалось создать клиент REDIS!\n")
	}

	r := chb.NewChangeBaseRedisRepository(rClient)
	s := ScrapperTAPI{
		StockChannel:  make(chan entity.StockInfo, 1000000),
		stopScrapping: make(chan bool),
		config:        config,
		trackedStocks: list.GetStocksInfoList(),
		logger:        config.Logger,
		redis:         r,
		analysis:      simple.Init(),
	}

	var err error
	s.driver, err = dr.NewApiDriver(config.TinkoffApiConfig, config.Logger, config.CandleDuration)
	if err != nil {
		s.logger.Errorf("не удалось создать драйвер tinkoff api: %s", err.Error())
		return &ScrapperTAPI{}, err
	}

	s.candleMod = config.CandleScrapperMode

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
			time.Sleep(5 * time.Minute)
			for _, stock := range stocks {
				info, isAnomaly := s.analysis.GetAnomaly(stock)

				if isAnomaly {
					infoStock, b := s.IsAnomaly(mappers.StockByMovementToStock(stock, info), stock, true)
					if b {
						s.StockChannel <- infoStock
					}
				}
				s.analysis.RemoveTicker(stock.Ticker)
			}
		}
	}()
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
				s.processingTrade(trade, stocks)
				break
			}
		}
	}()

	return s.StockChannel, nil
}

func (s *ScrapperTAPI) processingTrade(trade *investapi.Trade, stocks []entity.Stock) {
	var currentStock entity.Stock
	for _, stock := range stocks {
		if stock.FIGI == trade.Figi {
			currentStock = stock
			break
		}
	}
	stockInfo := entity.StockInfo{}

	stockInfo.Stock.Price = float64(trade.Price.GetUnits()) + float64(trade.Price.Nano)/1e9
	totalVolume := stockInfo.Stock.Price * float64(trade.Quantity) * float64(currentStock.MinLotCount)

	if trade.Direction == investapi.TradeDirection_TRADE_DIRECTION_BUY {
		stockInfo.StockMove = entity.Buy
	} else {
		if trade.Direction == investapi.TradeDirection_TRADE_DIRECTION_SELL {
			stockInfo.StockMove = entity.Sale
		} else {
			stockInfo.StockMove = entity.None
		}
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
		stockInfo.Stock.Name,
		stockInfo.Stock.Ticker,
		stockInfo.Stock.Price,
		stockInfo.LotsCount,
		stockInfo.StockMove)

	stockInfo, _ = s.IsAnomaly(stockInfo, currentStock, false)

	if stockInfo.StockMove != entity.None {
		ok := s.analysis.Add(stockInfo)
		if !ok {
			s.logger.Errorf("Не удалось добавить в систему анализа запись о сделке")
		}
		ok = s.redis.Add(entity.StockAdd{
			StockName: stockInfo.Stock.UID,
			Type:      stockInfo.StockMove.String(),
			NumPrice:  int64(stockInfo.Volume),
		})
		if !ok {
			s.logger.Errorf("Не удалось добавить в редис запись о сделке")
		}
	}

	s.StockChannel <- stockInfo
}

func (s *ScrapperTAPI) IsAnomaly(info entity.StockInfo, stock entity.Stock, doubleVolume bool) (entity.StockInfo, bool) {
	if info.Volume >= stock.AnomalySize {
		info.IsAnomaly = true
		s.logger.Warnf("Обнаружена аномалия: NAME:%s PRICE: %f ANOMALY SIZE: %f LOT COUNT: %d STOCK MOVE: %s",
			info.Stock.Name, info.Stock.Price, info.Volume, info.LotsCount, info.StockMove)
		err := s.driver.GetPerDayStatistics(&info)
		if err != nil {
			s.logger.Errorf("Ошибка получения доп.информации об акции: %s", err.Error())
			s.logger.Errorf("Информация об акции: NAME:%s PRICE: %f ANOMALY SIZE: %f LOT COUNT: %d STOCK MOVE: %s",
				info.Stock.Name, info.Stock.Price, info.Volume, info.LotsCount, info.StockMove)

		} else {
			s.logger.Info("Аномалия успешно обработана")
		}

		info.PerDaySalesVolume = float64(s.redis.Get(info.Stock.UID, entity.Sale.String()))
		if info.PerDaySalesVolume == 0 {
			s.logger.Warn("Аномальное значение суммы продаж: 0 !!!")
		}
		info.PerDayBuysVolume = float64(s.redis.Get(info.Stock.UID, entity.Buy.String()))
		if info.PerDayBuysVolume == 0 {
			s.logger.Warn("Аномальное значение суммы покупок: 0 !!!")
		}
		if doubleVolume {
			info.PerDayVolume = info.PerDaySalesVolume + info.PerDayBuysVolume
		} else {
			info.PerDayVolume = info.PerDaySalesVolume + info.PerDayBuysVolume + info.Volume
		}

		info.PerDaySalesPercent = (info.PerDaySalesVolume / info.PerDayVolume) * float64(100)
		info.PerDayBuysPercent = float64(100) - info.PerDaySalesPercent

		info.VolumeChange = ((float64(100) / (info.PerDaySalesVolume + info.PerDayBuysVolume)) * info.PerDayVolume) - float64(100)
		return info, true
	}
	return info, false
}

// Посылает сигнал для остановки скраппинга
func (s *ScrapperTAPI) StopScrape() {
	s.logger.Info("Скраппером получен сигнал на остановку")
	s.stopScrapping <- true
}

func (s ScrapperTAPI) skipTime() bool {
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
