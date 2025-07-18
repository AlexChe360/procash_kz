package telegram

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/AlexChe360/procash/internal/config"
	"github.com/AlexChe360/procash/internal/models"
	"github.com/AlexChe360/procash/internal/services/bot"
	"github.com/AlexChe360/procash/internal/services/freedom"
	"github.com/AlexChe360/procash/internal/services/rkeeper"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

func SendOrderInfo(cfg config.Config, db *gorm.DB, bot bot.BotClient, chatId int64, restaurantIdStr, tableNumberStr string) {

	bot.SendTyping(strconv.FormatInt(chatId, 10), 2*time.Second)

	restaurantID, _ := strconv.Atoi(restaurantIdStr)
	tableNumber, _ := strconv.Atoi(tableNumberStr)

	bot.SendTyping(strconv.FormatInt(chatId, 10), 1*time.Second)
	tableCode, err := rkeeper.GetTableCode(cfg, restaurantID, tableNumberStr)
	if err != nil {
		log.Println("❌ Ошибка при получении tableCode:", err)
		msg := tgbotapi.NewMessage(chatId, "Ошибка при получении информации о столе. Попробуйте позже.")
		bot.Send(msg)
		return
	}

	guid, waiterID, err := rkeeper.GetOrderInfo(cfg, restaurantID, tableCode)
	if err != nil {
		log.Println("❌ Ошибка при получении заказа:", err)
		msg := tgbotapi.NewMessage(chatId, "Ошибка при получении заказа. Попробуйте позже.")
		bot.Send(msg)
		return
	}

	_, amount, err := rkeeper.GetOrderDetails(cfg, restaurantID, guid)
	if err != nil {
		log.Println("❌ Ошибка при получении деталей заказа:", err)
		msg := tgbotapi.NewMessage(chatId, "Ошибка при получении суммы заказа. Попробуйте позже.")
		bot.Send(msg)
		return
	}

	waiterName, err := rkeeper.GetWaiterName(cfg, restaurantID, waiterID)
	if err != nil {
		log.Println("❌ Ошибка при получении имени официанта:", err)
		msg := tgbotapi.NewMessage(chatId, "Ошибка при получении имени официанта. Попробуйте позже.")
		bot.Send(msg)
		return
	}

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
		strconv.Itoa(tableNumber),
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
