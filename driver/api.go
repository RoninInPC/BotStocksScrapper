package driver

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"syscall"

	"BotStocksScrapper/entity"
	tsdk "github.com/tinkoff/invest-api-go-sdk/investgo"
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
			Name:        response.GetInstrument().Name,
			Ticker:      response.GetInstrument().Ticker,
			FIGI:        response.GetInstrument().Figi,
			UID:         response.GetInstrument().Uid,
			MinLotCount: int(response.GetInstrument().Lot),
			AnomalySize: stock.AnomalySize,
			Price:       0,
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

	go func() {
		err := tradeStream.Stream.Listen()
		if err != nil {
			d.logger.Errorf("не удалось запустить стрим: %s", err.Error())
			tradeStream.IsListen = false
		} else {
			tradeStream.IsListen = true
		}
	}()

	return tradeStream, nil
}
