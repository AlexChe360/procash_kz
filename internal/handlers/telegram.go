package handlers

import (
	"log"

	"github.com/AlexChe360/procash/internal/config"
	"github.com/AlexChe360/procash/internal/services/telegram"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TelegramWebhook(cfg config.Config, db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Println("Go to TelegramWebhook")
		return telegram.HadleWebhook(cfg, db, c)
	}
}
