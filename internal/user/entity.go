package user

import "github.com/google/uuid"

type User struct {
	ID         uuid.UUID `json:"id" db:"id"`
	TelegramID string    `json:"telegram_id" db:"telegram_id"`
	Username   string    `json:"username" db:"username"`
	Password   string    `json:"-" db:"password"`
}

type RegisterPayload struct {
	TelegramID string `json:"telegram_id" validate:"required"`
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"required,min=6,max=20"`
}

type LoginPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type UpdateUser struct {
	ID       uuid.UUID `json:"-"`
	Username string    `json:"username" validate:"required"`
}
