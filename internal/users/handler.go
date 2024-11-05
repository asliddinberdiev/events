package users

import "github.com/gofiber/fiber/v2"

type UserHandler struct {
	service *UserService
}

func NewUserHandler(serivce *UserService) *UserHandler {
	return &UserHandler{service: serivce}
}

func (h *UserHandler) RegisterRoutes(app *fiber.App, basePath string) {
	auth := app.Group(basePath + "/auth")
	{
		auth.Post("/register")
		auth.Post("/login")
	}
}
