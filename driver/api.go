package driver

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os/signal"
	"syscall"
	"time"

	"BotStocksScrapper/entity"
	tsdk "github.com/tinkoff/invest-api-go-sdk/investgo"
	investapi "github.com/tinkoff/invest-api-go-sdk/proto"
)

type ApiDriver struct {
	config             tsdk.Config
	client             *tsdk.Client
	instrumentsClient  *tsdk.InstrumentsServiceClient
	marketClient       *tsdk.MarketDataServiceClient
	marketStreamClient *tsdk.MarketDataStreamClient
	operationsClient   *tsdk.OperationsServiceClient
	ctx                context.Context
	logger             entity.Logger
}

func NewApiDriver(cfg tsdk.Config, lg entity.Logger) (*ApiDriver, error) {
	driver := &ApiDriver{
		config: cfg,
		logger: lg,
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	driver.ctx = ctx

	var err error
	driver.client, err = tsdk.NewClient(ctx, driver.config, driver.logger)
	if err != nil {
		return &ApiDriver{}, errors.New(fmt.Sprintf("client creating error %v", err.Error()))
	}

	driver.instrumentsClient = driver.client.NewInstrumentsServiceClient()
	driver.marketClient = driver.client.NewMarketDataServiceClient()
	driver.operationsClient = driver.client.NewOperationsServiceClient()
	driver.marketStreamClient = driver.client.NewMarketDataStreamClient()

	return driver, nil
}

// Инициализирует отслеживаемые акции.
// Заполняет поля структур entity.Stock для дальнейшего использования внутри драйвера
func (d *ApiDriver) InitStocks(trackedStocks []entity.TrackedStock) ([]entity.Stock, error) {
	stocks := []entity.Stock{}

	for _, stock := range trackedStocks {
		response, err := d.instrumentsClient.ShareByFigi(stock.FIGI)
		if err != nil || response == nil {
			d.logger.Errorf("не удалось получить ID инструмента по FIGI <%s>: %s", stock.FIGI, err.Error())
			return nil, err
		}

		stocks = append(stocks, entity.Stock{
			Name:         response.GetInstrument().Name,
			Ticker:       response.GetInstrument().Ticker,
			FIGI:         response.GetInstrument().Figi,
			UID:          response.GetInstrument().Uid,
			MinLotCount:  int(response.GetInstrument().Lot),
			AnomalySize:  stock.AnomalySize,
			RealExchange: response.GetInstrument().RealExchange.String(),
			Exchange:     response.GetInstrument().Exchange,
			Price:        0,
		})
	}

	return stocks, nil
}

// Запускает стрим подписки на обезличенные сделки
// Принимает список акций за обезличенными сделками которых необходимо следить
// Возвращает сущность entity.TradeStream из которой достается через канал сделки по акциям
func (d *ApiDriver) GetTradeCh(stocks []entity.Stock) (*entity.TradeStream, error) {

	var instrumentIDs []string
	for _, stock := range stocks {
		instrumentIDs = append(instrumentIDs, stock.UID)
	}

	tradeStream := &entity.TradeStream{}

	// Подключаемся к потоку обезличенных сделок
	var err error
	tradeStream.Stream, err = d.marketStreamClient.MarketDataStream()
	if err != nil {
		d.logger.Errorf("ошибка создания стрима биржевой информации: %s", err.Error())
		return nil, err
	}

	tradeStream.Channel, err = tradeStream.Stream.SubscribeTrade(instrumentIDs)
	if err != nil {
		d.logger.Errorf("ошибка подписки на стрим обезличенных сделок: %s", err.Error())
		return nil, err
	}
	tradeStream.IsListen = true

	go func() {
		err := tradeStream.Stream.Listen()
		if err != nil {
			d.logger.Errorf("не удалось запустить стрим: %s", err.Error())
			tradeStream.Stream = nil
			tradeStream.IsListen = false
		}
	}()

	return tradeStream, nil
}

// Принимает слайс структур entity.StockInfo и дозаполняет информацию об акциях
// Заполняет данные статистики за день и данные изменения на основе аномалии
// Сущность entity.StockInfo должна быть предзаполнена
func (d *ApiDriver) GetPerDayStatistics(stock *entity.StockInfo) error {

	// Получаем расписание текущей биржи инструмента
	rsp, err := d.instrumentsClient.TradingSchedules(stock.Stock.Exchange, time.Now(), time.Now())
	if err != nil {
		d.logger.Errorf("Ошибка получения расписания для биржи %s: %s", stock.Stock.Exchange, err.Error())
		return errors.New("Ошибка получения расписания биржи.")
	}

	if rsp.Exchanges == nil {
		d.logger.Errorf("Отсутствует расписание для биржи %s в ответе API", stock.Stock.Exchange)
		return errors.New("Отсутствует расписание для биржи.")
	}
	exch := rsp.GetExchanges()[0]

	var exchStartTime, exchEndTime time.Time
	for _, day := range exch.GetDays() {
		if day.Date.AsTime().Year() == time.Now().Year() && day.Date.AsTime().YearDay() == time.Now().YearDay() {
			if !day.IsTradingDay {
				d.logger.Warnf("Биржа %s закрыта. Данные не заполнены", stock.Stock.Exchange)
				return errors.New("Биржа закрыта. Невозможно получить данные")
			} else {
				exchStartTime = day.StartTime.AsTime()
				exchEndTime = day.EndTime.AsTime()
			}
		}
	}

	// Получаем свечи инструмента с момента открытия торгового дня
	response, err := d.marketClient.GetCandles(stock.Stock.UID, investapi.CandleInterval_CANDLE_INTERVAL_30_MIN, exchStartTime, exchEndTime)
	if err != nil {
		d.logger.Errorf("Ошибка получения свечей за день по акции %s-%s: %s", stock.Stock.Name, stock.Stock.Ticker, err.Error())
		return err
	}
	if response.Candles == nil {
		d.logger.Errorf("Для инструмента %s:%s нет свечей", stock.Stock.Name, stock.Stock.Ticker)
		return errors.New("Пустой массив свечей. Информация об акции не заполнена.")
	}

	candle := response.GetCandles()[0]
	stock.Volume = float64(candle.Volume)
	openPrice := float64(candle.Open.Units) + float64(candle.Open.Nano)/1e9

	// Получаем последнюю сделку по инструменту
	lpre, err := d.marketClient.GetLastPrices([]string{stock.Stock.UID})
	if err != nil {
		d.logger.Errorf("Ошибка получения последних сделок для инструмента %s: %s", stock.Stock.Ticker, err.Error())
		return err
	}

	lastPrice := float64(lpre.LastPrices[0].Price.Units) + float64(lpre.LastPrices[0].Price.Nano)/1e9

	// Если сегодня выходной, то метод GetLastPrices вернет последнюю сделку на бирже в последний рабочий день
	//  Тогда цена первой сделки считается от цены последней сделки в рабочий день
	today := time.Now()
	if today.Weekday() == time.Saturday || today.Weekday() == time.Sunday {
		stock.PerDayPriceChange = math.Round(((stock.Stock.Price-lastPrice)/lastPrice)*10000) / 100
		d.logger.Info("Выходной день. Процент изменения цены вычисляется по формуле выходного дня")
	} else {
		d.logger.Info("Рабочий день. Процент изменения цены вычисляется по формуле рабочего дня")
		stock.PerDayPriceChange = math.Round(((stock.Stock.Price-openPrice)/openPrice)*10000) / 100
	}

	var i int
	for i, candle = range response.GetCandles() {
		stock.PerDayVolume += candle.Volume
	}
	d.logger.Infof("Обработано %d свечей для акции %s:%s", i, stock.Stock.Name, stock.Stock.Ticker)

	return nil
}
