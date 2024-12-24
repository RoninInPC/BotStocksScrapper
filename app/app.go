package app

import (
	"BotStocksScrapper/hash"
	chb "BotStocksScrapper/repository/changeBase/impl"
	"BotStocksScrapper/repository/logBase/implLogBase"
	"BotStocksScrapper/usecase"
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
	app.ScrapperService, err = scrapper.NewScrapperService(app.Config, app.TgService.TelegramBot.BotApi, app.Config.ChatID)
	if err != nil {
		return nil, err
	}
	// Инициализируем крон
	cleaner := usecase.NewDatabaseCleaner(
		implLogBase.NewRedisRepository(
			implLogBase.NewRedisClient(app.Config.RedisLog), hash.Nothing),
		chb.NewChangeBaseRedisRepository(chb.NewChangeBaseClient(app.Config.RedisChange)))

	app.CronService = cron.Service{[]cron.Task{{
		Schedule: "0 3 * * *", // Каждый день в 3:00 утра
		Action:   cleaner.Clean},
	}}
	return app, nil
}

// Запускает скраппер-бота
func (a *App) Work() {
	go a.ScrapperService.Work()
	go a.CronService.Work()
	a.TgService.Work()
}
