package app

import (
	"BotStocksScrapper/entity"
	"BotStocksScrapper/service/bot"
	"BotStocksScrapper/service/cron"
	"BotStocksScrapper/service/scrapper"
)

type App struct {
	TgService       bot.BotService
	ScrapperService scrapper.ScrapperService
	CronService     cron.Service
	Config          entity.Config
	logger          entity.Logger
}
