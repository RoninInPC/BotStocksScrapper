package mappers

import (
	"fmt"

	"BotStocksScrapper/internal/entity"
	"BotStocksScrapper/internal/markdown"
)

func StockInfoToStringMsg(stock entity.StockInfo) string {
	volume, rng := VolumeFormater(stock.Volume)
	answerPreview := fmt.Sprintf("$%s %.2f%% %.2f%s\n", stock.Stock.Ticker, stock.VolumeChange, volume, rng)
	stockName := markdown.ToBold(stock.Stock.Name) + "\n"

	var answerDescr string
	if stock.StockMove != entity.None {
		answerDescr = "Инсайдерская сделка на "
		if stock.StockMove == entity.Sale {
			answerPreview = "🔴" + answerPreview
			answerDescr += fmt.Sprintf("продажу")
		} else {
			answerPreview = "🟢" + answerPreview
			answerDescr += fmt.Sprintf("покупку")
		}
		answerDescr = fmt.Sprintf("%s\n\n\n", markdown.ToBold(answerDescr))
	} else {
		answerDescr = "Аномальный объем"
		answerDescr = fmt.Sprintf("%s\n\n\n", markdown.ToBold(answerDescr))
	}

	answerBody := fmt.Sprintf("Цена: %.2f₽\n", stock.Stock.Price)
	answerBody += fmt.Sprintf("Объем: %.1f [%d %s]\n", stock.Volume, stock.LotsCount, LotWordEnding(stock.LotsCount))
	answerBody += fmt.Sprintf("Изменение на объеме: %.2f%%\n\n\n", stock.VolumeChange)

	dayStatistic := fmt.Sprintf("%s\n", markdown.ToBold("Итого за день"))
	dayStatistic += fmt.Sprintf("Изменение цены: %.2f%%\n", stock.PerDayPriceChange)
	dayStatistic += fmt.Sprintf("Покупки: %.2f%%, %.4f ₽\n", stock.PerDaySalesPercent, stock.PerDaySalesVolume)
	dayStatistic += fmt.Sprintf("Продажи: %.2f%%, %.4f ₽\n", stock.PerDayBuysPercent, stock.PerDayBuysVolume)

	answer := answerPreview + stockName + answerDescr + answerBody + dayStatistic
	return answer
}

func VolumeFormater(volume float64) (float64, string) {
	if volume >= 1000000000 {
		return volume / float64(1000000000), "млрд."
	} else if volume >= 1000000 {
		return volume / float64(1000000), "млн."
	} else {
		return volume / float64(1000), "тыс."
	}
}

func LotWordEnding(lots int64) string {
	if lots < 0 {
		lots = -lots
	}

	lastDigit := lots % 10
	lastTwoDigits := lots % 100

	if lastTwoDigits >= 11 && lastTwoDigits <= 19 {
		return "лотов"
	}

	switch lastDigit {
	case 1:
		return "лот"
	case 2, 3, 4:
		return "лота"
	default:
		return "лотов"
	}
}
