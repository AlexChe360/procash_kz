package main

import (
	"log"

	"github.com/AlexChe360/procash/internal/config"
	"github.com/AlexChe360/procash/internal/database"
	"github.com/AlexChe360/procash/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env not found, using system env vars")
	}

	cfg := config.Load()
	db := database.Open(cfg)

	app := fiber.New()

	handlers.RegisterRouter(app, cfg, db)

	log.Printf("Server running on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
