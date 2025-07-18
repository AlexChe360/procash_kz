package bot

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/AlexChe360/procash/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type WhatsappClient struct {
	Token     string
	PhoneID   string
	APIClient *http.Client
}

func NewWhatsappClient(cfg config.Config) *WhatsappClient {
	return &WhatsappClient{
		Token:     cfg.WhatsapApiToken,
		PhoneID:   cfg.WhatsappPhoneID,
		APIClient: http.DefaultClient,
	}
}

func (w *WhatsappClient) SendTyping(to string, duration time.Duration) {
	body := map[string]any{
		"message_product": "whatsapp",
		"to":              to,
		"type":            "typing_on",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(
		"POST",
		"https://graph.facebook.com/v23.0/"+w.PhoneID+"/messages",
		bytes.NewBuffer(jsonBody))

	req.Header.Set("Authorization", "Bearer "+w.Token)
	req.Header.Set("Content-Type", "application/json")
	w.APIClient.Do(req)
	time.Sleep(duration)
}

func (w *WhatsappClient) SendMessage(to string, text string) error {
	body := map[string]any{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "text",
		"text": map[string]string{
			"body": text,
		},
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST",
		"https://graph.facebook.com/v18.0/"+w.PhoneID+"/messages",
		bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+w.Token)
	req.Header.Set("Content-Type", "application/json")
	_, err := w.APIClient.Do(req)
	return err
}

func (w *WhatsappClient) Send(msg tgbotapi.Chattable) (tgbotapi.Message, error) {
	log.Println("⚠️ WhatsAppClient.Send called with Telegram msg — ignoring")
	return tgbotapi.Message{}, nil
}
