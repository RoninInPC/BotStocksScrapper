package telegram

import (
	"github.com/and3rson/telemux/v2"
)

type Action interface {
	Action(u *telemux.Update)
}

type SimpleAction func(telegramBot *TelegramBot, u *telemux.Update)

type SimpleActionStruct struct {
	SimpleAction SimpleAction
	Telegram     *TelegramBot
}

func (s SimpleActionStruct) Action(u *telemux.Update) {
	s.SimpleAction(s.Telegram, u)
}
