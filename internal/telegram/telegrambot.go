package telegram

import (
	"BotStocksScrapper/internal/list"
	"BotStocksScrapper/internal/mappers"
	"strconv"
	"strings"

	"github.com/and3rson/telemux/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	Commands Commands
	BotApi   *tgbotapi.BotAPI
	adminId  int64
}

func InitBot(token string, adminId int64) (*TelegramBot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TelegramBot{Commands: make(Commands, 0), BotApi: api, adminId: adminId}, nil
}

func (telegramBot *TelegramBot) AddCommand(command Command) {
	telegramBot.Commands = append(telegramBot.Commands, command)
}

func (telegramBot *TelegramBot) initBotMenu() {
	var sliceArr []tgbotapi.BotCommand
	for _, action := range telegramBot.Commands {
		if len(action.Description) > 0 {
			sliceArr = append(sliceArr, tgbotapi.BotCommand{
				Command:     action.Name,
				Description: action.Description,
			})
		}
	}
	cmdCfg := tgbotapi.NewSetMyCommands(
		sliceArr...,
	)
	_, _ = telegramBot.BotApi.Send(cmdCfg)
}

func (telegramBot *TelegramBot) getUpdates(timeOut int) tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = timeOut
	return telegramBot.BotApi.GetUpdatesChan(u)
}

func (telegramBot *TelegramBot) dispatchUpdates() {
	mux := telemux.NewMux()

	for _, command := range telegramBot.Commands {
		mux.AddHandler(telemux.NewHandler(command.Filter, func(u *telemux.Update) {
			command.Action.Action(u)
		}))
	}
	for update := range telegramBot.getUpdates(40) {
		mux.Dispatch(telegramBot.BotApi, update)
	}
}

func (telegramBot *TelegramBot) Work() {
	telegramBot.AddCommand(
		MakeCommandByFilterDefault(
			"start",
			"Начнём?",
			SimpleActionStruct{
				Telegram: telegramBot,
				SimpleAction: func(telegramBot *TelegramBot, u *telemux.Update) {
					_, _ = telegramBot.BotApi.Send(
						tgbotapi.NewMessage(
							u.FromChat().ID,
							"Здравствуйте, я бот, собирающий данные об аномальных объёмах. Подробнее об объёмах в /help.",
						),
					)
				},
			}),
	)
	telegramBot.AddCommand(
		MakeCommandByFilterDefault(
			"help",
			"Справка.",
			SimpleActionStruct{
				Telegram: telegramBot,
				SimpleAction: func(telegramBot *TelegramBot, u *telemux.Update) {
					msg := tgbotapi.NewMessage(
						u.FromChat().ID,
						mappers.ToTelegramFormat(list.GetStocksInfoList()),
					)
					msg.ParseMode = tgbotapi.ModeMarkdown
					_, _ = telegramBot.BotApi.Send(msg)
				},
			}),
	)
	telegramBot.AddCommand(
		MakeCommandByFilterDefault(
			"set",
			"Инициализация акции",
			SimpleActionStruct{
				Telegram: telegramBot,
				SimpleAction: func(telegramBot *TelegramBot, u *telemux.Update) {
					if u.SentFrom().ID != telegramBot.adminId {
						msg := tgbotapi.NewMessage(
							u.FromChat().ID,
							"У вас нет доступа к этой функции бота",
						)
						msg.ParseMode = tgbotapi.ModeMarkdown
						_, _ = telegramBot.BotApi.Send(msg)
						return
					}
					splitted := strings.Split(u.Message.Text, " ")
					if len(splitted) < 4 {
						msg := tgbotapi.NewMessage(
							u.FromChat().ID,
							"Некорректная команда",
						)
						msg.ParseMode = tgbotapi.ModeMarkdown
						_, _ = telegramBot.BotApi.Send(msg)
						return
					}
					stockName := splitted[1]
					stockVolume, err := strconv.ParseFloat(splitted[2], 64)
					if err != nil {
						stockVolume = 100 * 1000000
					}
					stockSoloVolume, err := strconv.ParseFloat(splitted[3], 64)
					if err != nil {
						stockSoloVolume = 50 * 1000000
					}
					list.Set(&list.StockScrapeInfo{StockTag: stockName, AnomalySizeVolume: stockVolume, AnomalySizeSolo: stockSoloVolume})
					list.SaveInFile()
						
				},
			}),
	)
	telegramBot.AddCommand(
		MakeCommandByFilterDefault(
			"remove",
			"Удалить акцию",
			SimpleActionStruct{
				Telegram: telegramBot,
				SimpleAction: func(telegramBot *TelegramBot, u *telemux.Update) {
					if u.SentFrom().ID != telegramBot.adminId {
						msg := tgbotapi.NewMessage(
							u.FromChat().ID,
							"У вас нет доступа к этой функции бота",
						)
						msg.ParseMode = tgbotapi.ModeMarkdown
						_, _ = telegramBot.BotApi.Send(msg)
						return
					}
					splitted := strings.Split(u.Message.Text, " ")
					if len(splitted) < 2 {
						msg := tgbotapi.NewMessage(
							u.FromChat().ID,
							"Некорректная команда",
						)
						msg.ParseMode = tgbotapi.ModeMarkdown
						_, _ = telegramBot.BotApi.Send(msg)
						return
					}
					stockName := splitted[1]
					list.Remove(&list.StockScrapeInfo{StockTag: stockName})
					list.SaveInFile()
				},
			}),
	)
	telegramBot.initBotMenu()
	telegramBot.dispatchUpdates()
}
