package http

import (
	"github.com/asliddinberdiev/events/internal/http/middleware"
	"github.com/asliddinberdiev/events/internal/user"
	"github.com/gofiber/fiber/v2"
)

func NewRouter() *fiber.App {
	app := fiber.New()
	app.Use(middleware.Cors())


	userRepo := user.NewPostgresUserRepository()
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	basepath := "/api/v1"
	public := app.Group(basepath)
	private := app.Group(basepath, middleware.AuthJWT)
	
	userHandler.RegisterRoutesPublic(public)
	userHandler.RegisterRoutesPrivate(private)

	return app
}
