package entity

import (
	"fmt"

	tsdk "github.com/tinkoff/invest-api-go-sdk/investgo"
	investapi "github.com/tinkoff/invest-api-go-sdk/proto"
)

type StockMoveType string

const (
	Sale StockMoveType = "Sale"
	Buy  StockMoveType = "Buy"
	None StockMoveType = "None"
)

func (s StockMoveType) String() string {
	return string(s)
}

type MarketType string

const (
	Tinkoff MarketType = "TQBR"
	Moscow  MarketType = "MOEX"
)

func (m MarketType) String() string {
	return string(m)
}

// Stock Структура описывающая акцию. Содержит базовую информацию об акции
//
// Name - Имя акции (Например 'Тинькофф акции')
// Ticker - Тег акции (Например TCSG)
// FIGI - Уникальный идентификатор инструмента(акции)
// UID - Уникальный идентификкатор инструмента внутри Tinkoff API
// MinLotCount - Минимальное количество акций для торгов (количество акций в лоте)
// Price - Цена одной акции
// AnomalySizeVolume - Объем от которого сделка считается аномальной
type Stock struct {
	Name              string
	Ticker            string
	FIGI              string
	UID               string
	MinLotCount       int
	RealExchange      string
	Exchange          string
	Price             float64
	AnomalySizeVolume float64
	AnomalySizeSolo   float64
}

// Структура описывающая аномалию.
// Содержит информацию об акции и параметрах аномалии.
//
// Stock - Базовая информация об акции
// Volume - Объем аномалии
// VolumeChange - Изменение на объеме в процентах
// LotsCount - Количесво лотов в аномалии
// StockMove - Тип аномалии покупка/продажа
// PerDayVolume - Общий объем продаж и покупок за день. Выражается в количестве лотов.
// PerDayPriceChange - Изменение цены за день
// PerDaySalesVolume - Объем продаж за день (в кол-ве лотов)
// PerDaySalesPercent - Процент продаж за день от общего объема
// PerDayBuysVolume - Объем покупок за день (в кол-ве лотов)
// PerDayBuysPercent - Процент покупок за день от общего объема
type StockInfo struct {
	Stock              Stock
	Volume             float64
	VolumeChange       float64
	LotsCount          int64
	StockMove          StockMoveType
	IsAnomaly          bool
	PerDayVolume       float64
	PerDayPriceChange  float64
	PerDaySalesVolume  float64
	PerDaySalesPercent float64
	PerDayBuysVolume   float64
	PerDayBuysPercent  float64
	InfoByFiveMin      StocksByMoveType
}

// Возвращает строку от StockInfo
func (s *StockInfo) StringFull() string {
	return fmt.Sprintf("Name:%s;Ticker:%s;FIGI:%s;UID:%s;Price:%d;Volume:%d;LotsCount:%d;StockMove:%s",
		s.Stock.Name, s.Stock.Ticker, s.Stock.FIGI, s.Stock.UID, s.Stock.Price, s.Volume, s.LotsCount, s.StockMove)
}
func (s *StockInfo) StringShort() string {
	return fmt.Sprintf("Price: %d; Volume: %d; LotsCount: %d;", s.Stock.Price, s.Volume, s.LotsCount)
}

// Отслеживаемые акции
type TrackedStock struct {
	Ticker      string  // Тег акции (например TCSG)
	FIGI        string  // Уникальный идентификатор инструмента
	AnomalySize float64 // Граница аномалии
}

// Сущность описывающая подписку на стакан обезличенных сделок
type TradeStream struct {
	Stream   *tsdk.MarketDataStream  // Стрим данных
	Channel  <-chan *investapi.Trade // Выходной канал в который отправляются сделки
	IsListen bool                    // Флаг успешного запуска. Необходимо проверять перед началом прослушивания канала
}

type StockAdd struct {
	StockName string
	Type      string // SALE или BUY
	NumPrice  int64
}

type CandleStockInfo struct {
	Stock       Stock
	Volume      float64
	TradesCount int64
	SalesCount  int64
	BuysCount   int64

	LotsCount int64
	StockMove StockMoveType
	IsAnomaly bool
}
