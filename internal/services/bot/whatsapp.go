package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/AlexChe360/procash/internal/config"
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

	url := fmt.Sprintf("https://graph.facebook.com/v23.0/%s/messages", w.PhoneID)

	body := map[string]any{
		"message_product": "whatsapp",
		"to":              to,
		"type":            "typing_on",
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+w.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.APIClient.Do(req)
	if err != nil {
		log.Println("❌ Ошибка отправки запроса:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("❌ WhatsApp API error (%d): %s", resp.StatusCode, bodyBytes)
	}

	time.Sleep(duration)
}

func (w *WhatsappClient) SendMessage(to string, text string) error {

	url := fmt.Sprintf("https://graph.facebook.com/v23.0/%s/messages", w.PhoneID)

	body := map[string]any{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "text",
		"text": map[string]string{
			"body": text,
		},
	}
	jsonBody, _ := json.Marshal(body)

	log.Printf("➡️ Отправка WhatsApp на %s: %s", to, text)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+w.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.APIClient.Do(req)
	if err != nil {
		log.Println("❌ Ошибка отправки запроса:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("❌ WhatsApp API error (%d): %s", resp.StatusCode, bodyBytes)
		return fmt.Errorf("whatsapp api error: %s", bodyBytes)
	}

	log.Println("✅ Успешно отправлено WhatsApp сообщение")

	return err
}
