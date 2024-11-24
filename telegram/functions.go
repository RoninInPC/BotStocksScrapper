package telegram

import (
	"github.com/and3rson/telemux/v2"
)

// можно сделать любую структуру с любыми аргументами, которые будут подставляться в Action
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
