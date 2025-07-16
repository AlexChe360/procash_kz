package database

import (
	"github.com/AlexChe360/procash/internal/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Payment{},
		&models.WhatsappLog{},
		&models.TelegramLog{},
	)
}
