package bot

import (
	"BotStocksScrapper/internal/telegram"
)

type BotService struct {
	TelegramBot *telegram.TelegramBot
}

func (botService *BotService) Work() {
	botService.TelegramBot.Work()
}
