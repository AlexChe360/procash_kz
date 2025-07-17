package config

import "os"

type Config struct {
	Port                string
	DBUrl               string
	RKeeperToken        string
	RKeeperBaseURL      string
	TGBotToken          string
	TelegramBotUsername string
	WhatsappPhone       string
	WhatsappPhoneID     string
	WhatsapApiToken     string
	MerchantID          string
	PaymentSecretKey    string
	PaymentUserId       string
	PaymentURL          string
	DefaultRestaurantID string
	WhatsappApiVersion  string
}

func Load() Config {
	return Config{
		Port:                getEnv("PORT", "8000"),
		DBUrl:               getEnv("DATABASE_URL", "file:procash.db"),
		RKeeperToken:        getEnv("RKEEPER_API_TOKE", ""),
		RKeeperBaseURL:      getEnv("RKEEPER_BASE_URL", ""),
		TGBotToken:          getEnv("TG_BOT_TOKEN", ""),
		TelegramBotUsername: getEnv("TG_BOT_NAME", ""),
		WhatsappPhone:       getEnv("WA_PHONE", ""),
		WhatsappPhoneID:     getEnv("WA_PHONE_NUMBER_ID", ""),
		WhatsapApiToken:     getEnv("WA_BOT_TOKEN", ""),
		MerchantID:          getEnv("FREEDOM_MERCHANT_ID", ""),
		PaymentSecretKey:    getEnv("FREEDOM_SECRET_KEY", ""),
		PaymentUserId:       getEnv("FREEDOM_USER_ID", ""),
		PaymentURL:          getEnv("FREEDOM_API_URL", ""),
		DefaultRestaurantID: getEnv("DEFAULT_RESTAURANT_ID", ""),
		WhatsappApiVersion:  getEnv("WA_VESRION_API", ""),
	}
}

func getEnv(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
