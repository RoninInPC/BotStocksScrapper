package mappers

import (
	"fmt"

	"BotStocksScrapper/internal/entity"
	"BotStocksScrapper/internal/markdown"
)

func StockInfoToStringMsg(stock entity.StockInfo) string {
	volume, rng := VolumeFormater(stock.Volume)
	answerPreview := fmt.Sprintf("#%s %.2f%% %.2f%s\n", stock.Stock.Ticker, stock.VolumeChange, volume, rng)
	stockName := markdown.ToBold(stock.Stock.Name) + "\n"

	var answerDescr string
	if stock.StockMove != entity.None {
		answerDescr = "❗️ИНСАЙДЕРСКАЯ СДЕЛКА на "
		if stock.StockMove == entity.Sale {
			answerDescr = "🔴" + answerDescr
			answerDescr += fmt.Sprintf("продажу:\n")
		} else {
			answerDescr = "🟢" + answerDescr
			answerDescr += fmt.Sprintf("покупку:\n")
		}
		answerDescr = "#единичныесделки\n" + answerDescr
		answerDescr, answerPreview = answerPreview, answerDescr
		answerDescr = fmt.Sprintf("%s\n\n\n", markdown.ToBold(answerDescr))
	} else {
		answerDescr = "Аномальный объем:"
		answerDescr = fmt.Sprintf("%s\n", markdown.ToBold(answerDescr))
	}

	answerBody := fmt.Sprintf("%s %.2f₽\n",
		markdown.ToItalic("Цена:"),
		stock.Stock.Price)
	answerBody += fmt.Sprintf("%s %.1f \\[%d %s]\n",
		markdown.ToItalic("Объем:"),
		stock.Volume,
		stock.LotsCount,
		LotWordEnding(stock.LotsCount))
	answerBody += fmt.Sprintf("%s %.2f%%\n\n",
		markdown.ToItalic("Изменение на объеме:"),
		stock.VolumeChange)

	statisticFiveMin := ""
	if stock.StockMove == entity.None {
		slice := stock.InfoByFiveMin.ToSlice()

		statisticFiveMin += fmt.Sprintf("%s\n", markdown.ToBold("Статистика:"))

		statisticFiveMin += fmt.Sprintf("%s %d\n",
			markdown.ToItalic("Количество сделок:"), slice.Len())

		avgLots := int64(slice.AvgLots())
		statisticFiveMin += fmt.Sprintf("%s %.1f \\[%d %s]\n",
			markdown.ToItalic("Средний объём:"),
			slice.AvgVolume(),
			avgLots,
			LotWordEnding(avgLots))

		medLots := slice.MedianLots()
		statisticFiveMin += fmt.Sprintf("%s %.1f \\[%d %s]\n",
			markdown.ToItalic("Медианный объём:"),
			slice.MedianVolume(),
			medLots,
			LotWordEnding(medLots))

		minLots := slice.MinLots()
		statisticFiveMin += fmt.Sprintf("%s %.1f \\[%d %s]\n",
			markdown.ToItalic("Минимальный объём:"),
			slice.MinVolume(),
			minLots,
			LotWordEnding(minLots))

		maxLots := slice.MaxLots()
		statisticFiveMin += fmt.Sprintf("%s %.1f \\[%d %s]\n",
			markdown.ToItalic("Максимальный объём:"),
			slice.MaxVolume(),
			maxLots,
			LotWordEnding(maxLots))

		bought := stock.InfoByFiveMin[entity.Buy]
		saled := stock.InfoByFiveMin[entity.Sale]
		if bought.SumVolume() < saled.SumVolume() {
			answerPreview = "🔴" + answerPreview
		} else {
			answerPreview = "🟢" + answerPreview
		}
		statisticFiveMin += fmt.Sprintf("%s %.2f%%\n",
			markdown.ToItalic("Покупки:"),
			bought.SumVolume()/stock.Volume*100)
		statisticFiveMin += fmt.Sprintf("%s %.2f%%\n\n",
			markdown.ToItalic("Продажи:"),
			saled.SumVolume()/stock.Volume*100)
	} else {
		answerBody += "\n"
	}

	dayStatistic := fmt.Sprintf("%s\n", markdown.ToBold("Итого за день:"))
	dayStatistic += fmt.Sprintf("%s %.2f%%\n", markdown.ToItalic("Изменение цены:"), stock.PerDayPriceChange)
	dayStatistic += fmt.Sprintf("%s %.2f%%, %.4f ₽\n", markdown.ToItalic("Покупки:"), stock.PerDaySalesPercent, stock.PerDaySalesVolume)
	dayStatistic += fmt.Sprintf("%s %.2f%%, %.4f ₽\n", markdown.ToItalic("Продажи:"), stock.PerDayBuysPercent, stock.PerDayBuysVolume)

	answer := answerPreview + stockName + answerDescr + answerBody + statisticFiveMin + dayStatistic
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
