package database

import (
	"log"

	"github.com/AlexChe360/procash/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Open(cfg config.Config) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("database open:", err)
	}

	if err := AutoMigrate(db); err != nil {
		log.Fatal("auto-migrate", err)
	}

	return db
}
