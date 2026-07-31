package telegram

import (
	"BotStocksScrapper/list"
	"BotStocksScrapper/mappers"
	"github.com/and3rson/telemux/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	Commands Commands
	BotApi   *tgbotapi.BotAPI
}

func InitBot(token string) (*TelegramBot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TelegramBot{Commands: make(Commands, 0), BotApi: api}, nil
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
	telegramBot.initBotMenu()
	telegramBot.dispatchUpdates()
}
