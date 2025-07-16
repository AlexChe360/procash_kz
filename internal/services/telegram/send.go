package telegram

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/AlexChe360/procash/internal/config"
	"github.com/AlexChe360/procash/internal/models"
	"github.com/AlexChe360/procash/internal/services/freedom"
	"github.com/AlexChe360/procash/internal/services/rkeeper"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

func SendOrderInfo(cfg config.Config, db *gorm.DB, chatId int64, restaurantIdStr, tableNumberStr string) {
	bot, err := tgbotapi.NewBotAPI(cfg.TGBotToken)
	if err != nil {
		return
	}

	restaurantID, _ := strconv.Atoi(restaurantIdStr)
	tableNumber, _ := strconv.Atoi(tableNumberStr)

	tableCode, err := rkeeper.GetTableCode(cfg, restaurantID, tableNumberStr)
	guid, waiterID, _ := rkeeper.GetOrderInfo(cfg, restaurantID, tableCode)
	_, amount, _ := rkeeper.GetOrderDetails(cfg, restaurantID, guid)
	waiterName, _ := rkeeper.GetWaiterName(cfg, restaurantID, waiterID)

	description := fmt.Sprintf("Оплата счета: #%s", guid)

	payment, err := freedom.GenerateURL(cfg, amount, description)
	if err != nil {
		log.Println("Ошибка при создании оплаты:", err)
		msg := tgbotapi.NewMessage(chatId, "Произошла ошибка при генерации оплаты. Попробуйте позже")
		bot.Send(msg)
		return
	}

	_ = db.Create(&models.TelegramLog{
		TelegramID:   chatId,
		RestaurantID: restaurantID,
		TableNumber:  tableNumberStr,
		WaiterName:   waiterName,
		OrderGUID:    guid,
		Amount:       amount,
		PayURL:       payment["redirect_url"],
		CreatedAt:    time.Now().UTC(),
	})

	text := fmt.Sprintf(
		"*Счёт к оплате*\n\nСтол: `%s`\nСумма: *%d* ₸\nОфициант: `%s`\n\n[Оплатить счёт](%s)",
		tableNumber,
		amount,
		waiterName,
		payment["redirect_url"],
	)

	message := tgbotapi.NewMessage(chatId, text)
	message.ParseMode = "Markdown"
	message.DisableWebPagePreview = true

	_, err = bot.Send(message)
	if err != nil {
		log.Println("Ошибка отправки Telegram-сообщения:", err)
	}
}
