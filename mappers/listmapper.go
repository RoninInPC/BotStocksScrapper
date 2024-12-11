package mappers

import (
	"BotStocksScrapper/list"
	"BotStocksScrapper/markdown"
	"fmt"
)

var helloString = "Я бот, который просматривает аномальные объёмы по акциям:"

func ToTelegramFormat(list []list.StockScrapeInfo) string {
	l := ""
	for _, v := range list {
		count, name := getNumStr(v.AnomalySize)
		l += fmt.Sprintf("%s превышающие %d %s,\n", markdown.ToItalic(v.StockTag), count, name)
	}
	return markdown.ToBold(helloString) + "\n" + l
}

func getNumStr(a uint64) (int, string) {
	if a > 1000000000 && a < 1000000000000 {
		return int(a / 1000000000), "млрд"
	}
	if a > 1000000 && a < 1000000000 {
		return int(a / 1000000), "млн"
	}
	if a > 1000 && a < 1000000 {
		return int(a / 1000), "тыс"
	}
	return int(a), ""
}
