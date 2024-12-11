package entity

import "fmt"

type StockMoveType string

const (
	Sale     StockMoveType = "Sale"
	Purchase StockMoveType = "Purchase"
)

// Структура описывающая акцию. Содержит базовую информацию
//
// Name - Имя акции (Например 'Тинькофф акции')
// Tag - Тег акции (Например TCSG)
// FIGI - Уникальный идентификатор инструмента(акции) в Tinkoff API
type Stock struct {
	Name string
	Tag  string
	FIGI string
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
	Stock                 Stock
	Price                 uint64
	Volume                uint64
	VolumeChange          float64
	StockMove             StockMoveType
	PriceChangePerDay     float64
	PurchasesPerDay       float64
	PurchasesPerDayVolume uint64
	SalesPerDay           float64
	SalesPerDayVolume     uint64
}

// Возвращает строку от StockInfo с информацией о полях Stock.Name, Stock.Tag, Stock.FIGI, Price, StockMove, Volume
func (s *StockInfo) String() string {
	return fmt.Sprintf("Name:%s;Tag:%s;FIGI:%s;Price:%d;StockMove:%s;Volume:%d",
		s.Stock.Name, s.Stock.Tag, s.Stock.FIGI, s.Price, s.StockMove, s.Volume)
}

// Отслеживаемые акции
type TrackedStock struct {
	StockTag    string `yaml:"ticker"`       // Тег акции (например TCSG)
	AnomalySize uint64 `yaml:"anomaly-size"` // Граница аномалии
}
