package http

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	sharedjwt "github.com/moneymate-2026/moneymate-backend/shared/pkg/jwt"
	response "github.com/moneymate-2026/moneymate-backend/shared/pkg/responses"
)

const localsUserID = "userID"

func RequireInternalSecret(secret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		provided := c.Get("X-Internal-Secret")
		if provided == "" || provided != secret {
			return response.Unauthorized(c, "internal access required")
		}
		return c.Next()
	}
}

func RequireUserID(cfg sharedjwt.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Unauthorized(c, "authentication required")
		}

		claims, err := sharedjwt.ParseAccessToken(parts[1], cfg.AccessSecret)
		if err != nil {
			return response.Unauthorized(c, "authentication required")
		}

		c.Locals(localsUserID, claims.UserID)
		return c.Next()
	}
}

func userIDFromLocals(c fiber.Ctx) string {
	id, _ := c.Locals(localsUserID).(string)
	return id
}

func resolveEnv(key string) string {
	return os.Getenv(key)
}
