package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AlexChe360/procash/internal/config"
)

type ButtonTemplatePayload struct {
	MessagingProduct string   `json:"messaging+product"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Template         Template `json:"template"`
}

type Template struct {
	Name       string      `json:"name"`
	Language   Language    `json:"language"`
	Components []Component `json:"components"`
}

type Language struct {
	Code string `json:"code"`
}

type Component struct {
	Type       string      `json:"type"`
	SubType    string      `json:"sub_type,omitempty"`
	Index      int         `json:"index,omitempty"`
	Parameters []Parameter `json:"parameters"`
}

type Parameter struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func SendButtons(cfg config.Config, phone, tableNumber, waiterName string, items []map[string]any, total int, payURL string) error {
	itemsText := ""
	for _, item := range items {
		itemsText += fmt.Sprintf("%s — %v ₸; ", item["name"], item["amount"])
	}

	payload := ButtonTemplatePayload{
		MessagingProduct: "whatsapp",
		To:               phone,
		Type:             "template",
		Template: Template{
			Name:     "procash_ru",
			Language: Language{Code: "ru"},
			Components: []Component{
				{
					Type: "body",
					Parameters: []Parameter{
						{Type: "text", Text: tableNumber},
						{Type: "text", Text: waiterName},
						{Type: "text", Text: itemsText},
						{Type: "text", Text: fmt.Sprintf("%d ₸", total)},
					},
				},
				{Type: "button", SubType: "url", Index: 0, Parameters: []Parameter{{Type: "text", Text: payURL}}},
				{Type: "button", SubType: "url", Index: 1, Parameters: []Parameter{{Type: "text", Text: payURL}}},
				{Type: "button", SubType: "url", Index: 2, Parameters: []Parameter{{Type: "text", Text: payURL}}},
			},
		},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", fmt.Sprintf("https://graph.facebook.com/v23.0/%s/messages", cfg.WhatsappPhoneID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.WhatsapApiToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
