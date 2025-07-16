package handlers

import (
	"log"
	"time"

	"github.com/AlexChe360/procash/internal/config"
	"github.com/AlexChe360/procash/internal/models"
	"github.com/AlexChe360/procash/internal/services/telegram"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func FreedomCallback(cfg config.Config, db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		form := struct {
			Status     string `form:"pg_result"`
			PaymentID  string `form:"pg_payment_id"`
			OrderID    string `form:"pg_order_id"`
			UserID     string `form:"pg_user_id"`
			Amount     int    `form:"amount"`
			TelegramID int64  `form:"tg_id"`
		}{}

		if err := c.BodyParser(&form); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if form.Status != "1" {
			log.Println("Оплата не успешна")
			return c.SendString("not success")
		}

		payment := models.Payment{
			OrderID:   form.OrderID,
			Amount:    int64(form.Amount),
			PaidAt:    time.Now(),
			CreatedAt: time.Now(),
		}
		if err := db.Create(&payment).Error; err != nil {
			log.Println("Ошибка сохраннеия оплпаты:", err)
		}

		if form.TelegramID != 0 {
			telegram.SendThanks(cfg, form.TelegramID)
		}

		return c.SendString("ok")
	}
}
