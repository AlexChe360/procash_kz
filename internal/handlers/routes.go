package handlers

import (
	"github.com/AlexChe360/procash/internal/config"
	"github.com/AlexChe360/procash/internal/services/bot"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRouter(app *fiber.App, cfg config.Config, db *gorm.DB, bot bot.BotClient) {
	// QR лендинг
	app.Get("/", QRHandler(cfg, db))
	app.Get("/privacy", PrivacyHandler())
	app.Static("/assets", "./static/assets")

	// Webhook
	app.Post("/freedom", FreedomCallback(cfg, db))
	app.Post("/telegram", TelegramWebhook(cfg, db, bot))
	app.Post("/whatsapp", WhatsappWebhook(cfg, db))

	// Пинг для проверки
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})
}
