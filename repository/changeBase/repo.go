package changeBase

import (
	"BotStocksScrapper/entity"
)

type CBRepository interface {
	Add(stock entity.StockAdd) bool
	Get(stockName, operationType string) int64
	Free() bool
}
