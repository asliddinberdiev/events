package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/asliddinberdiev/events/conf"
	"github.com/asliddinberdiev/events/internal/common"
	"github.com/asliddinberdiev/events/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const UserKey = "user_id"

func CreateJWT(secret string, userID uuid.UUID) (string, error) {
	expiration := time.Now().Add(time.Second * time.Duration(conf.Envs.App.JWTExpirationInSeconds)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		UserKey:     userID,
		"expiredAt": expiration,
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AuthJWT(ctx *fiber.Ctx) error {
	tokenString := getTokenFromRequest(ctx)

	token, err := validateToken(tokenString)
	if err != nil || !token.Valid {
		return permissionDenied(ctx)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return permissionDenied(ctx)
	}

	userID, userExists := claims[UserKey].(string)
	if !userExists {
		return permissionDenied(ctx)
	}

	ctx.Locals(UserKey, userID)

	return ctx.Next()
}

func getTokenFromRequest(c *fiber.Ctx) string {
	tokenAuth := c.Get("Authorization")

	if tokenAuth != "" {
		parts := strings.Split(tokenAuth, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}
	return ""
}

func validateToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(conf.Envs.App.JWTSecret), nil
	})
}

func permissionDenied(c *fiber.Ctx) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "permission denied"})
}

func MakeRequestSearch(c *fiber.Ctx) *common.SearchRequest {
	var query common.SearchRequest

	query.Limit = uint16(utils.ParseInt(c.Query("limit"), 10))
	query.Page = uint16(utils.ParseInt(c.Query("page"), 1))
	query.Search = c.Query("search", "")

	return &query
}

func MakeRequest(c *fiber.Ctx) *common.Request {
	var request common.Request

	if userID, ok := c.Context().Value(UserKey).(string); ok {
		request.UserID = userID
	}

	return &request
}

