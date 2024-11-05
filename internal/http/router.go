package http

import (
	"github.com/asliddinberdiev/events/internal/http/middleware"
	"github.com/asliddinberdiev/events/internal/users"
	"github.com/gofiber/fiber/v2"
)

func NewRouter() *fiber.App {
	app := fiber.New()
	app.Use(middleware.Cors())

	basepath := "/api/v1"

	userRepo := users.NewPostgresUserRepository()
	userService := users.NewUserService(userRepo)
	userHandler := users.NewUserHandler(userService)
	userHandler.RegisterRoutes(app, basepath)

	return app
}
