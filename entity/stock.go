package entity

import (
	"fmt"
	"sync"

	tsdk "github.com/tinkoff/invest-api-go-sdk/investgo"
	investapi "github.com/tinkoff/invest-api-go-sdk/proto"
)

type StockMoveType string

const (
	Sale StockMoveType = "Sale"
	Buy  StockMoveType = "Buy"
)

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
// AnomalySize - Объем от которого сделка считается аномальной
type Stock struct {
	Name         string
	Ticker       string
	FIGI         string
	UID          string
	MinLotCount  int
	RealExchange string
	Exchange     string
	Price        float64
	AnomalySize  float64
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
	PerDayVolume       int64
	PerDayPriceChange  float64
	PerDaySalesVolume  float64
	PerDaySalesPercent float64
	PerDayBuysVolume   float64
	PerDayBuysPercent  float64
}

// Возвращает строку от StockInfo
func (s *StockInfo) String() string {
	return fmt.Sprintf("Name:%s;Ticker:%s;FIGI:%s;UID:%s;Price:%d;Volume:%d;LotsCount:%d;StockMove:%s",
		s.Stock.Name, s.Stock.Ticker, s.Stock.FIGI, s.Stock.UID, s.Stock.Price, s.Volume, s.LotsCount, s.StockMove)
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

var lock sync.Mutex
var trackedStocks = []TrackedStock{
	{"LKOH", "BBG004731032", 80000000},
	{"SBER", "BBG004730N88", 80000000},
	{"GAZP", "BBG004730RP0", 80000000},
	{"ROSN", "BBG004731354", 70000000},
	{"NVTK", "BBG00475KKY8", 800},
	{"T", "TCS80A107UL4", 80000000},
}

// шиза на счёт гонки данных scrapper-ом и командой бота
func GetStocksInfoList() []TrackedStock {
	lock.Lock()
	defer lock.Unlock()
	answer := make([]TrackedStock, len(trackedStocks))
	copy(answer, trackedStocks)
	return answer
}
