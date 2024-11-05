package users

import "github.com/google/uuid"

type User struct {
	ID         uuid.UUID `json:"id"`
	TelegramID uuid.UUID `json:"telegram_id" validate:"required"`
	Username   string    `json:"username" validate:"required,min=6,max=20"`
	Password   string    `json:"password" validate:"required,min=6,max=20"`
}


