package http

import (
	"crypto/subtle"

	"github.com/gofiber/fiber/v3"
	response "github.com/moneymate-2026/moneymate-backend/shared/pkg/responses"
)

func RequireUserID(c fiber.Ctx) error {
	userID := c.Get("X-User-Id")
	if userID == "" {
		return response.Unauthorized(c, "missing authentication context")
	}
	c.Locals("userID", userID)
	return c.Next()
}

func RequireInternalSecret(secret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		if secret == "" {
			return response.InternalServerError(c)
		}
		provided := c.Get("X-Internal-Secret")
		if subtle.ConstantTimeCompare([]byte(provided), []byte(secret)) != 1 {
			return response.Unauthorized(c, "unauthorized internal call")
		}
		return c.Next()
	}
}