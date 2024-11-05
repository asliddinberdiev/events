package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ParseInt(value string, defaultValue int) int {
	if parsedValue, err := strconv.Atoi(value); err == nil {
		return parsedValue
	}
	return defaultValue
}

func ParseQueryByID(c *fiber.Ctx) (uuid.UUID, error) {
	idString := c.Params("id", "")
	id, err := uuid.Parse(idString)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
