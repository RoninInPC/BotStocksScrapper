package app

import (
	"io"
	"os"

	"BotStocksScrapper/config"
	"BotStocksScrapper/entity"
	"BotStocksScrapper/service/bot"
	"BotStocksScrapper/service/cron"
	"BotStocksScrapper/service/scrapper"
	"BotStocksScrapper/telegram"
	"github.com/sirupsen/logrus"
)

// Создает главную сущность приложения
//
//	Заполняет конфиг приложения
//	Инициализирует логгер
//	Инициализирует вложенные сервисы
func NewApp() (*App, error) {
	app := &App{}

	var err error
	app.Config, err = config.LoadConfig("./config/config.yaml")
	if err != nil {
		return nil, err
	}

	// Инициализируем логер
	file, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	writer := io.MultiWriter(os.Stdout, file)
	app.Config.Logger = entity.NewLogger(writer, logrus.DebugLevel)

	app.logger = app.Config.Logger

	// Инициализируем сервис телеграм бота
	app.TgService = bot.BotService{}
	bot, err := telegram.InitBot(app.Config.TgToken)
	if err != nil {
		return nil, err
	}
	app.TgService.TelegramBot = bot

	// Инициализируем скраппер
	// TODO необходимо передать chatID
	app.ScrapperService, err = scrapper.NewScrapperService(app.Config, app.TgService.TelegramBot.BotApi, 0)
	if err != nil {
		return nil, err
	}

	// Инициализируем крон
	app.CronService = cron.Service{}
	// TODO возможно необходимо еще заполнить таски??

	return app, nil
}

// Запускает скраппер-бота
func (a *App) Work() error {
	// TODO implement this method
	return nil
}

// Останавливает все сервисы и бота скраппера
func (a *App) Stop() error {
	// TODO implement this method
	return nil
}
