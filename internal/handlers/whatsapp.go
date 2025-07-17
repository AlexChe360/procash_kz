package handlers

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/AlexChe360/procash/internal/config"
	"github.com/AlexChe360/procash/internal/models"
	"github.com/AlexChe360/procash/internal/services/freedom"
	"github.com/AlexChe360/procash/internal/services/rkeeper"
	"github.com/AlexChe360/procash/internal/services/whatsapp"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func WhatsappWebhook(cfg config.Config, db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body map[string]any
		if err := c.BodyParser(&body); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		log.Println("📩 WhatsApp Webhook received")

		entries, ok := body["entry"].([]any)
		if !ok || len(entries) == 0 {
			log.Println("⚠️ Нет entry")
			return c.SendStatus(fiber.StatusOK)
		}
		entry, ok := entries[0].(map[string]any)
		if !ok {
			log.Println("⚠️ Неверный формат entry")
			return c.SendStatus(fiber.StatusOK)
		}

		changesList, ok := entry["changes"].([]any)
		if !ok || len(changesList) == 0 {
			log.Println("⚠️ Нет changes")
			return c.SendStatus(fiber.StatusOK)
		}
		change, ok := changesList[0].(map[string]any)
		if !ok {
			log.Println("⚠️ Неверный формат changes[0]")
			return c.SendStatus(fiber.StatusOK)
		}

		value, ok := change["value"].(map[string]any)
		if !ok {
			log.Println("⚠️ Нет value")
			return c.SendStatus(fiber.StatusOK)
		}
		messages, ok := value["messages"].([]any)
		if !ok || len(messages) == 0 {
			log.Println("⚠️ Нет входящих сообщений")
			return c.SendStatus(fiber.StatusOK)
		}

		messageData, ok := messages[0].(map[string]any)
		if !ok {
			log.Println("⚠️ Неверный формат message")
			return c.SendStatus(fiber.StatusOK)
		}

		textMap, ok := messageData["text"].(map[string]any)
		if !ok {
			log.Println("⚠️ Нет поля text")
			return c.SendStatus(fiber.StatusOK)
		}
		text, ok := textMap["body"].(string)
		if !ok {
			log.Println("⚠️ Нет поля body")
			return c.SendStatus(fiber.StatusOK)
		}

		log.Println("📩 Текст сообщения:", text)

		from, ok := messageData["from"].(string)
		if !ok {
			log.Println("⚠️ Нет поля from")
			return c.SendStatus(fiber.StatusOK)
		}

		meta := ""
		if strings.Contains(text, "meta=") {
			meta = text[strings.Index(text, "meta=")+5:]
		}
		parts := strings.Split(meta, "-")
		if len(parts) != 2 {
			log.Println("⚠️ Неверный формат meta:", meta)
			return c.SendStatus(fiber.StatusOK)
		}

		restaurantID, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Println("⚠️ Невалидный restaurantID:", parts[0])
			return c.SendStatus(fiber.StatusOK)
		}

		tableNumber := parts[1]

		tableCode, err := rkeeper.GetTableCode(cfg, restaurantID, tableNumber)
		if err != nil {
			log.Println("❌ tableCode:", err)
			return c.SendStatus(fiber.StatusOK)
		}

		orderGUID, waiterID, err := rkeeper.GetOrderInfo(cfg, restaurantID, tableCode)
		if err != nil {
			log.Println("❌ orderInfo:", err)
			return c.SendStatus(fiber.StatusOK)
		}

		items, totalSum, err := rkeeper.GetOrderDetails(cfg, restaurantID, orderGUID)
		if err != nil {
			log.Println("❌ orderDetails:", err)
			return c.SendStatus(fiber.StatusOK)
		}

		waiterName, err := rkeeper.GetWaiterName(cfg, restaurantID, waiterID)
		if err != nil {
			log.Println("❌ waiterName:", err)
			waiterName = "Неизвестно"
		}

		payment, err := freedom.GenerateURL(cfg, totalSum, "Оплата счёта")
		if err != nil {
			log.Println("❌ FreedomPay:", err)
			return c.SendStatus(fiber.StatusOK)
		}

		err = whatsapp.SendButtons(
			cfg,
			from,
			tableNumber,
			waiterName,
			items,
			totalSum,
			payment["redirect_url"],
		)
		if err != nil {
			log.Println("❌ Отправка WhatsApp:", err)
		} else {
			log.Println("✅ Отправлено WhatsApp:", from)
		}

		_ = db.Create(&models.WhatsappLog{
			Phone:        from,
			RestaurantID: restaurantID,
			TableNumber:  tableNumber,
			WaiterName:   waiterName,
			OrderGUID:    orderGUID,
			Amount:       totalSum,
			PayURL:       payment["redirect_url"],
			CreatedAt:    Now(),
		})

		return c.SendStatus(fiber.StatusOK)
	}
}

func Now() time.Time {
	return time.Now().UTC()
}
