package middlewares

import (
	"log"
	"strings"
	"github.com/gofiber/fiber/v3"
	"github.com/moneymate-2026/moneymate-backend/gateway/internal/proxy"
	jwtutil "github.com/moneymate-2026/moneymate-backend/shared/pkg/jwt"
)

func RequireAuth(jwtSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format, expected: Bearer <token>",
			})
		}

		token := parts[1]

		claims, err := jwtutil.ParseAccessToken(token, jwtSecret)
		if err != nil {
			log.Print(err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		role := "user"
		if len(claims.Roles) > 0 {
			role = claims.Roles[0]
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", role)
		return c.Next()
	}
}

func RequireTransactionAuth(authClient proxy.AuthClient) fiber.Handler {
	return func(c fiber.Ctx) error {
		transactionToken := c.Get("X-Transaction-Token")
		if transactionToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing x-transaction-token header",
			})
		}

		transactionID := c.Get("X-Transaction-ID")
		if transactionID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "missing x-transaction-id header",
			})
		}

		claims, err := authClient.VerifyTransactionToken(c.Context(), transactionToken, transactionID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired transaction token",
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("transaction_id", claims.TransactionID)

		return c.Next()
	}
}

