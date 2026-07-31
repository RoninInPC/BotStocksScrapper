package botservice

import (
	"BotStocksScrapper/telegram"
)

type BotService struct {
	TelegramBot *telegram.TelegramBot
}

func (botService *BotService) Work() {
	botService.TelegramBot.Work()
}
