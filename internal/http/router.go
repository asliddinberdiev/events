package http

import (
	"github.com/asliddinberdiev/events/internal/http/middleware"
	"github.com/gofiber/fiber/v2"
)

func NewRouter() *fiber.App {
	app := fiber.New()
	app.Use(middleware.Cors())

	return app
}
