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

		entry := body["entry"].([]any)[0].(map[string]any)
		changes := entry["changes"].([]any)[0].(map[string]any)
		value := changes["value"].(map[string]any)
		messages := value["messages"].([]any)
		if len(messages) == 0 {
			log.Println("⚠️ Нет входящих сообщений")
			return c.SendStatus(fiber.StatusOK)
		}

		message := messages[0].(map[string]any)
		text := message["text"].(map[string]any)["body"].(string)
		from := message["from"].(string)

		meta := ""
		if strings.Contains(text, "meta=") {
			meta = text[strings.Index(text, "meta=")+5:]
		}
		parts := strings.Split(meta, "-")
		if len(parts) != 2 {
			log.Println("⚠️ Неверный формат meta")
			return c.SendStatus(fiber.StatusOK)
		}

		restaurantID, err1 := strconv.Atoi(parts[0])
		tableNumber := parts[1]
		if err1 != nil {
			return c.SendStatus(fiber.StatusOK)
		}

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
