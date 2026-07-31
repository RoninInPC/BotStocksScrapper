package botservice

import (
	"BotStocksScrapper/telegram"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBotService(t *testing.T) {
	bot, err := telegram.InitBot("7097380480:AAEu5s1B4Yp38UF-CrB3jVhn5qvKhiIF-2A")
	assert.Equal(t, nil, err)
	service := BotService{TelegramBot: bot}
	go service.Work()
	time.Sleep(time.Minute)
}
