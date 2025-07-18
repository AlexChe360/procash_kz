package bot

import (
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramClient struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramClient(api *tgbotapi.BotAPI) *TelegramClient {
	return &TelegramClient{bot: api}
}

func (t *TelegramClient) SendTyping(chatID string, duration time.Duration) {
	id, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		log.Println("Telegram typing: invalid chatID:", err)
		return
	}

	msg := tgbotapi.NewChatAction(id, tgbotapi.ChatTyping)
	t.bot.Send(msg)
	time.Sleep(duration)
}

func (t *TelegramClient) SendMessage(chatID string, text string) error {
	id, _ := strconv.ParseInt(chatID, 10, 64)
	msg := tgbotapi.NewMessage(id, text)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true
	_, err := t.bot.Send(msg)
	return err
}

func (c *TelegramClient) Send(msg tgbotapi.Chattable) (tgbotapi.Message, error) {
	return c.bot.Send(msg)
}
