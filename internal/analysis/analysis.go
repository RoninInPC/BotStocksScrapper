package analysis

import "BotStocksScrapper/internal/entity"

type Analysis interface {
	Add(entity.StockInfo) bool
	GetInfoByTicker(string) entity.StocksByMoveType
	RemoveTicker(string) bool
	RemoveAll() bool
	GetAnomaly(entity.Stock) (entity.StocksByMoveType, bool)
}
