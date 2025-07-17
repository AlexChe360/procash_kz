package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func SendWhatsAppMessage(token, phoneNumberID, to, message string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v23.0/%s/messages", phoneNumberID)

	payload := map[string]any{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "text",
		"text": map[string]string{
			"body": message,
		},
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("❌ Ошибка отправки запроса:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		responseBody, _ := io.ReadAll(resp.Body)
		log.Printf("❌ WhatsApp API error: %s\n", responseBody)
		return fmt.Errorf("WhatsApp API error: %s", responseBody)
	}

	return nil
}
