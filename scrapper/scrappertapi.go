package scrapper

import (
	"errors"

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

				s.logger.Infof("Получена обезличенная сделка: NAME: %s; TICKER: %s; PRICE: %f; LOT_COUNT: %d; MOVE: %s",
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

					// TODO
					//  Получаем из бд все сделки к текущему моменту
					//  Дозаполняем StockInfo
					//   Дозаполняем поля
					//    PerDaySalesVolume  float64
					//	  PerDaySalesPercent float64
					//	  PerDayBuysVolume   float64
					//	  PerDayBuysPercent  float64
					// PerDayBuysVolume сумма всех лотов сделок на покупку
					// allStocks = db.GetAllStocksInfo(stockInfo.String()) -> []StockInfo
					// for _, stock := range allStocks
					//		if stock.StockMove == entity.Buy {
					//			stockInfo.PerDayBuysVolume += stock.LotsCount
					//		}
					//  stockInfo.PerDaySalesVolume = stockInfo.PerDayVolume - stockInfo.PerDayBuysVolume
					//  stockInfo.PerDaySalesPercent / stockInfo.PerDaySalesPercent - это просто процентное соотношение от общей суммы (stockInfo.PerDayVolume)
					//
				}
				// TODO Добавляем в бд информацию о сделке

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
