package telegram

import (
	"errors"

	"BotStocksScrapper/internal/entity"
	"BotStocksScrapper/internal/mappers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramSender struct {
	tgClient *tgbotapi.BotAPI
	chatID   int64
}

func NewSender(tgClient *tgbotapi.BotAPI, chatID int64) TelegramSender {
	return TelegramSender{tgClient: tgClient, chatID: chatID}
}

func (s TelegramSender) Send(stock entity.StockInfo) error {
	msg := mappers.StockInfoToStringMsg(stock)

	if s.tgClient == nil {
		return errors.New("Клиент телеграма не инициализирован")
	}
	chatMessage := tgbotapi.NewMessage(s.chatID, msg)
	_, err := s.tgClient.Send(chatMessage)

	return err
}
