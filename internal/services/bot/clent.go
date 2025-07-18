package bot

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotClient interface {
	SendTyping(chatId string, duration time.Duration)
	SendMessage(chatId, text string) error
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

type WhatsAppClient interface {
	SendTyping(to string, duration time.Duration)
	SendMessage(to, text string) error
}
