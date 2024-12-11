package driver

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"BotStocksScrapper/entity"
	"github.com/sirupsen/logrus"
	tsdk "github.com/tinkoff/invest-api-go-sdk/investgo"
	investapi "github.com/tinkoff/invest-api-go-sdk/proto"
)

type ApiDriver struct {
	config            tsdk.Config
	client            *tsdk.Client
	instrumentsClient *tsdk.InstrumentsServiceClient
	marketClient      *tsdk.MarketDataServiceClient
	operationsClient  *tsdk.OperationsServiceClient
	ctx               context.Context
	logger            *logrus.Logger
}

func NewApiDriver(cfg tsdk.Config, lg *logrus.Logger) (*ApiDriver, error) {
	driver := &ApiDriver{
		config: cfg,
		logger: lg,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()
	driver.ctx = ctx

	var err error
	driver.client, err = tsdk.NewClient(ctx, driver.config, driver.logger)
	if err != nil {
		return &ApiDriver{}, errors.New(fmt.Sprintf("client creating error %v", err.Error()))
	}

	driver.instrumentsClient = driver.client.NewInstrumentsServiceClient()
	driver.marketClient = driver.client.NewMarketDataServiceClient()
	driver.operationsClient = driver.client.NewOperationsServiceClient()

	return driver, nil
}

func (d *ApiDriver) GetStockInfo(figi string) (entity.StockInfo, error) {

	shareRsp, err := d.instrumentsClient.ShareByFigi(figi)
	if err != nil {
		d.logger.Errorf("ошибка получения информации об акции: %s", err.Error())
		return entity.StockInfo{}, err
	}

	stock := entity.Stock{
		Name: shareRsp.GetInstrument().Name,
		Tag:  shareRsp.GetInstrument().Ticker,
		FIGI: figi,
	}

	// Получение текущей цены
	lastPricesResp, err := d.marketClient.GetLastPrices([]string{figi})
	if err != nil {
		d.logger.Errorf("ошибка получения текущей цены акции: %s", err.Error())
		return entity.StockInfo{}, fmt.Errorf("failed to get last prices: %v", err)
	}

	var currentPrice float64
	for _, price := range lastPricesResp.LastPrices {
		if price.Figi == figi {
			currentPrice = float64(price.Price.Units) + float64(price.Price.Nano)/1e9
			break
		}
	}

	if currentPrice == 0 {
		d.logger.Errorf("инструмент с FIGI %s не найден в последних ценах", figi)
		return entity.StockInfo{}, fmt.Errorf("instrument with FIGI %s not found in last prices", figi)
	}

	// Получение сделок за день
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	operationsResp, err := d.operationsClient.GetOperations(&tsdk.GetOperationsRequest{
		AccountId: d.config.AccountId,
		From:      startOfDay,
		To:        now,
		Figi:      figi,
	})
	if err != nil {
		d.logger.Errorf("не удалось получить сделки за день для FIGI %s: %s", figi, err.Error())
		return entity.StockInfo{}, fmt.Errorf("failed to get operations: %v", err)
	}

	var totalVolumeRub, totalVolumeLots, totalBuysRub, totalSellsRub float64
	var totalBuys, totalSells int64
	for _, op := range operationsResp.Operations {
		amount := float64(op.Payment.Units) + float64(op.Payment.Nano)/1e9
		if op.OperationType == investapi.OperationType_OPERATION_TYPE_BUY {
			totalBuysRub += amount
			totalBuys++
		} else if op.OperationType == investapi.OperationType_OPERATION_TYPE_SELL {
			totalSellsRub += amount
			totalSells++
		}
		totalVolumeRub += amount
		totalVolumeLots += float64(op.Quantity)
	}

	// Получение свечей для изменения цены
	candlesResp, err := d.marketClient.GetCandles(figi, investapi.CandleInterval_CANDLE_INTERVAL_DAY, startOfDay, now)
	if err != nil {
		d.logger.Errorf("не удалось получить свечи для FIGI %s: %s", figi, err.Error())
		return entity.StockInfo{}, fmt.Errorf("failed to get candles: %v", err)
	}

	var priceChangePct float64
	if len(candlesResp.Candles) > 0 {
		candle := candlesResp.Candles[0]
		openPrice := float64(candle.Open.Units) + float64(candle.Open.Nano)/1e9
		priceChangePct = ((currentPrice - openPrice) / openPrice) * 100
	}

	// Определение "покупка или продажа" и итоговые данные
	var buyOrSell entity.StockMoveType
	if totalBuys > totalSells {
		buyOrSell = entity.Purchase
	} else if totalSells > totalBuys {
		buyOrSell = entity.Sale
	}

	d.logger.Infof("успешно получена информация для иснтрумента FIGI %s", figi)

	return entity.StockInfo{
		Stock:                 stock,
		Price:                 currentPrice,
		Volume:                totalVolumeRub,
		VolumeChange:          priceChangePct,
		StockMove:             buyOrSell,
		PriceChangePerDay:     priceChangePct,
		PurchasesPerDay:       (totalBuysRub / totalVolumeRub) * 100,
		PurchasesPerDayVolume: totalBuysRub,
		SalesPerDay:           (totalSellsRub / totalVolumeRub) * 100,
		SalesPerDayVolume:     totalSellsRub,
	}, nil
}
