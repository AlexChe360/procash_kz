package telegram

import (
	"encoding/json"
	"strings"

	"github.com/AlexChe360/procash/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func HadleWebhook(cfg config.Config, db *gorm.DB, c *fiber.Ctx) error {
	var update tgbotapi.Update
	if err := json.Unmarshal(c.Body(), &update); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if update.Message != nil && update.Message.IsCommand() {
		cmd := update.Message.Command()
		args := update.Message.CommandArguments()

		switch cmd {
		case "start":
			parts := strings.Split(args, "_")
			if len(parts) != 2 {
				return c.SendStatus(fiber.StatusBadRequest)
			}
			restaurantId := parts[0]
			tableNumber := parts[1]
			go SendOrderInfo(cfg, db, update.Message.Chat.ID, restaurantId, tableNumber)
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
