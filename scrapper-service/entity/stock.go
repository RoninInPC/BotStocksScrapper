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
)

// Stock Структура описывающая акцию. Содержит базовую информацию об акции
//
// Name - Имя акции (Например 'Тинькофф акции')
// Ticker - Тег акции (Например TCSG)
// FIGI - Уникальный идентификатор инструмента(акции)
// UID - Уникальный идентификкатор инструмента внутри Tinkoff API
// MinLotCount - Минимальное количество акций для торгов (количество акций в лоте)
// Price - Цена одной акции
type Stock struct {
	Name        string
	Ticker      string
	FIGI        string
	UID         string
	MinLotCount int
	Price       float64
	AnomalySize float64
}

// Структура описывающая аномалию.
// Содержит информацию об акции и параметрах аномалии.
//
// Stock - Базовая информация об акции
// Price - Цена акции в копейках
// Volume - Объем аномалии в копейках
// VolumeChange - Изменение на объеме в процентах
// StockMove - Тип аномалии покупка/продажа
// PriceChangePerDay - Изменение цены за день в процентах
// PurchasesPerDay - Покупок за день в процентах
// PurchasesPerDayVolume - Объем покупок за день в копейках
// SalesPerDay - Продаж за день в процентах
// SalesPerDayVolume - Объем продаж за день в копейках
type StockInfo struct {
	Stock      Stock
	Volume     float64
	NumberLots int64
	StockMove  StockMoveType
}

// Возвращает строку от StockInfo с информацией о полях Stock.Name, Stock.Tag, Stock.FIGI, Price, StockMove, Volume
func (s *StockInfo) String() string {
	return fmt.Sprintf("Name:%s;Ticker:%s;FIGI:%s;UID:%s;Price:%d;Volume:%d;NumberLots:%d;StockMove:%s",
		s.Stock.Name, s.Stock.Ticker, s.Stock.FIGI, s.Stock.UID, s.Stock.Price, s.Volume, s.NumberLots, s.StockMove)
}

// Отслеживаемые акции
type TrackedStock struct {
	Ticker      string  `yaml:"ticker"`       // Тег акции (например TCSG)
	FIGI        string  `yaml:"figi"`         // Уникальный идентификатор инструмента
	AnomalySize float64 `yaml:"anomaly-size"` // Граница аномалии
}

// Сущность описывающая подписку на стакан обезличенных сделок
type TradeStream struct {
	Stream   *tsdk.MarketDataStream  // Стрим данных
	Channel  <-chan *investapi.Trade // Выходной канал в который отправляются сделки
	IsListen bool                    // Флаг успешного запуска. Необходимо проверять перед началом прослушивания канала
}
