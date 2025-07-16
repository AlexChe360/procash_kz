package handlers

import (
	"bytes"
	"html/template"
	"log"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func PrivacyHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {

		tmplPath := filepath.Join("static", "privacy", "index.html")
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			log.Println("template parse error:", err)
			return c.Status(500).SendString("Template error")
		}

		data := map[string]string{}

		var outputBuffer bytes.Buffer
		if err := tmpl.Execute(&outputBuffer, data); err != nil {
			log.Println("template exec error:", err)
			return c.Status(500).SendString("Template exec error")
		}

		return c.Type("html").SendStream(&outputBuffer)

	}
}
