package telegram

import (
	"log"

	"github.com/AlexChe360/procash/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendThanks(cfg config.Config, chatID int64) {
	bot, err := tgbotapi.NewBotAPI(cfg.TGBotToken)
	if err != nil {
		log.Println("Telegram init error:", err)
		return
	}

	msg := tgbotapi.NewMessage(chatID, "Спасибо за оплату!\nВаш платёж успешно получен.")
	bot.Send(msg)
}
