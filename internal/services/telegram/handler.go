package telegram

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/AlexChe360/procash/internal/config"
	"github.com/AlexChe360/procash/internal/services/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func HadleWebhook(cfg config.Config, db *gorm.DB, bot bot.BotClient, c *fiber.Ctx) error {
	var update tgbotapi.Update
	if err := json.Unmarshal(c.Body(), &update); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if update.Message != nil && update.Message.IsCommand() {
		cmd := update.Message.Command()
		args := update.Message.CommandArguments()

		switch cmd {
		case "start":
			log.Printf("📥 Получен args: %s", args)
			parts := strings.Split(args, "_")
			if len(parts) != 2 {
				log.Printf("⚠️ Invalid /start payload: %s", args)
				return c.SendStatus(fiber.StatusBadRequest)
			}
			restaurantId := parts[0]
			tableNumber := parts[1]

			log.Printf("✅ /start with restaurantId=%s, tableNumber=%s", restaurantId, tableNumber)
			go SendOrderInfo(cfg, db, bot, update.Message.Chat.ID, restaurantId, tableNumber)
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
