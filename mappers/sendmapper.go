package mappers

import (
	"fmt"

	"BotStocksScrapper/entity"
)

func StockInfoToStringMsg(stock entity.StockInfo) string {
	answerPreview := fmt.Sprintf("🔴$%s %f%% %f\n", stock.Stock.Ticker, stock.VolumeChange, stock.Volume)
	answerDescr := "Аномальный объем на "
	if stock.StockMove == entity.Sale {
		answerDescr += "продажу"
	} else {
		answerDescr += "покупку"
	}

	answerBody := fmt.Sprintf("Цена: %f₽\n", stock.Stock.Price)
	answerBody += fmt.Sprintf("Объем: %f[%d лотов]\n", stock.Volume, stock.LotsCount)
	answerBody += fmt.Sprintf("Изменение на объеме: %f%%\n", stock.VolumeChange)
	if stock.StockMove == entity.Sale {
		answerBody += fmt.Sprintf("Тип: продажа\n")
	} else {
		answerDescr += fmt.Sprintf("Тип: покупка\n")
	}

	dayStatistic := fmt.Sprintf("Итого за день:\n")
	dayStatistic += fmt.Sprintf("Изменение цены: %s%%\n", stock.PerDayPriceChange)
	dayStatistic += fmt.Sprintf("Покупки: %f%%, %f ₽\n", stock.PerDaySalesPercent, stock.PerDaySalesVolume)
	dayStatistic += fmt.Sprintf("Продажи: %f%%, %f ₽\n", stock.PerDayBuysPercent, stock.PerDayBuysVolume)

	answer := answerPreview + answerDescr + answerBody + dayStatistic
	return answer
}
