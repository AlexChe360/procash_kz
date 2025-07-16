package handlers

import (
	"bytes"
	"html/template"
	"log"
	"path/filepath"

	"github.com/AlexChe360/procash/internal/config"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func QRHandler(cfg config.Config, db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		restaurantID := c.Query("restaurantId")
		tableNumber := c.Query("tableNumber")

		if restaurantID == "" || tableNumber == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing restaurantId or tableNumber")
		}

		data := map[string]string{
			"TableNumber": tableNumber,
			"TelegramURL": "https://t.me/" + cfg.TelegramBotUsername + "?start=" + restaurantID + "_" + tableNumber,
			"WhatsappURL": "https://wa.me/" + cfg.WhatsappPhone + "?text=meta=" + restaurantID + "-" + tableNumber,
		}

		tmplPath := filepath.Join("static", "order", "index.html")
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			log.Println("template parse error:", err)
			return c.Status(500).SendString("Template error")
		}

		var outputBuffer bytes.Buffer
		if err := tmpl.Execute(&outputBuffer, data); err != nil {
			log.Println("template exec error:", err)
			return c.Status(500).SendString("Template exec error")
		}

		return c.Type("html").SendStream(&outputBuffer)

	}
}
