package mappers

import "BotStocksScrapper/internal/entity"

func StockByMovementToStock(stock entity.Stock, moveType entity.StocksByMoveType) entity.StockInfo {
	var stockInfo entity.StockInfo
	sliceStock := moveType.ToSlice()
	stockInfo.Stock.Ticker = stock.Ticker
	stockInfo.Stock.FIGI = stock.FIGI
	stockInfo.Stock.UID = stock.UID
	stockInfo.Stock.Price = sliceStock.AvgPrice()
	stockInfo.Stock.Name = stock.Name
	stockInfo.Stock.MinLotCount = stock.MinLotCount
	stockInfo.Stock.Exchange = stock.Exchange
	stockInfo.Stock.RealExchange = stock.RealExchange
	stockInfo.Volume = sliceStock.SumVolume()
	stockInfo.LotsCount = sliceStock.SumLots()
	stockInfo.StockMove = entity.None
	stockInfo.InfoByFiveMin = moveType.Copy()
	return stockInfo
}
