package user

import (
	"log"

	"github.com/asliddinberdiev/events/conf"
	"github.com/asliddinberdiev/events/internal/common"
	"github.com/asliddinberdiev/events/internal/http/middleware"
	"github.com/asliddinberdiev/events/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(serivce *UserService) *UserHandler {
	return &UserHandler{service: serivce}
}
func (h *UserHandler) RegisterRoutesPublic(app fiber.Router) {
	auth := app.Group("/auth")
	{
		auth.Post("/register", h.Register)
		auth.Post("/login", h.Login)
	}
}	

func (h *UserHandler) RegisterRoutesPrivate(app fiber.Router) {
	users := app.Group("/users")
	{
		users.Get("", h.getAll)
		users.Get("/:id", h.getByID)
		users.Put("/:id", h.update)
		users.Delete("/:id", h.delete)
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var input RegisterPayload

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid body")
	}

	if err := middleware.Validate.Struct(input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invaled Body fields")
	}

	existingTelegramID, err := h.service.GetByTelegramID(c.Context(), input.TelegramID)
	if err != nil {
		log.Println("err: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error checking for existing telegram_id")
	}

	if existingTelegramID != nil {
		return fiber.NewError(fiber.StatusConflict, "Telegram id already exists")
	}

	existingUser, err := h.service.GetByUsername(c.Context(), input.Username)
	if err != nil {
		log.Println("err: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error checking for existing user")
	}

	if existingUser != nil {
		return fiber.NewError(fiber.StatusConflict, "Username already exists")
	}

	hashPassword, err := middleware.HashPassword(input.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error: creating hashed password")
	}

	newUser, err := h.service.Create(c.Context(), RegisterPayload{
		TelegramID: input.TelegramID,
		Username:   input.Username,
		Password:   hashPassword,
	})
	if err != nil {
		log.Println("panisssssss: ", err)

		return fiber.NewError(fiber.StatusInternalServerError, "Error: creating user")
	}

	token, err := middleware.CreateJWT(conf.Envs.App.JWTSecret, newUser.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error: creating token")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"token": token,
		"user":  newUser,
	})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var input LoginPayload

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid body")
	}

	if err := middleware.Validate.Struct(input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invaled Body fields")
	}

	existingUser, err := h.service.GetByUsername(c.Context(), input.Username)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error checking for existing user")
	}

	if existingUser == nil {
		return fiber.NewError(fiber.StatusConflict, "Wrong username or password")
	}

	if !middleware.ComparePassword(existingUser.Password, input.Password) {
		return fiber.NewError(fiber.StatusConflict, "Wrong username or password")
	}

	token, err := middleware.CreateJWT(conf.Envs.App.JWTSecret, existingUser.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error: creating token")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"token": token,
		"user":  existingUser,
	})
}

func (h *UserHandler) getAll(c *fiber.Ctx) error {
	query := middleware.MakeRequestSearch(c)

	params := common.SearchRequest{Limit: uint16(query.Limit), Page: uint16(query.Page)}

	res, err := h.service.GetAll(c.Context(), params)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed get all user")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total": res.Total,
		"list":  res.List,
	})
}

func (h *UserHandler) getByID(c *fiber.Ctx) error {
	id, err := utils.ParseQueryByID(c)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	res, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed get by id user")
	}

	if res == nil {
		return fiber.NewError(fiber.StatusNotFound, "not found")
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *UserHandler) update(c *fiber.Ctx) error {
	id, err := utils.ParseQueryByID(c)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	var input UpdateUser
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid body")
	}

	if err := middleware.Validate.Struct(input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invaled Body fields")
	}
	input.ID = id

	existingUsername, err := h.service.GetByUsername(c.Context(), input.Username)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error checking for existing user")
	}

	if existingUsername == nil {
		return fiber.NewError(fiber.StatusConflict, "Username already exists")
	}

	existingUser, err := h.service.GetByID(c.Context(), input.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error checking for existing user")
	}

	if existingUser == nil {
		return fiber.NewError(fiber.StatusNotFound, "not found")
	}

	if err := h.service.Update(c.Context(), input); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed update user")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":       input.ID,
		"username": input.Username,
	})
}

func (h *UserHandler) delete(c *fiber.Ctx) error {
	id, err := utils.ParseQueryByID(c)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	existingUser, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error checking for existing user")
	}

	if existingUser == nil {
		return fiber.NewError(fiber.StatusNotFound, "not found")
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed delete user")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": id})
}
