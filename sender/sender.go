package sender

import (
	"errors"

	"BotStocksScrapper/entity"
	"BotStocksScrapper/mappers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender struct {
	tgClient *tgbotapi.BotAPI
	chatID   int64
}

func NewSender(tgClient *tgbotapi.BotAPI, chatID int64) *Sender {
	return &Sender{tgClient: tgClient, chatID: chatID}
}

func (s *Sender) SendMsg(stock entity.StockInfo) error {
	msg := mappers.StockInfoToStringMsg(stock)

	if s.tgClient == nil {
		return errors.New("Клиент телеграма не инициализирован")
	}
	chatMessage := tgbotapi.NewMessage(s.chatID, msg)
	_, err := s.tgClient.Send(chatMessage)

	return err
}
