package http

import (
	"slices"
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/domain"
	sharedjwt "github.com/moneymate-2026/moneymate-backend/shared/pkg/jwt"
	response "github.com/moneymate-2026/moneymate-backend/shared/pkg/responses"
)

const (
	localsRecipientID   = "recipientID"
	localsRecipientType = "recipientType"
)

// RequireRecipient verifies the standard access token and records who the
// caller is. recipient_type is "merchant" when the token carries a merchant
// role, otherwise "user". Every /notification route is authenticated — there
// is no anonymous access.
func RequireRecipient(cfg sharedjwt.Config) fiber.Handler {
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
		c.Locals(localsRecipientID, claims.UserID)
		if slices.Contains(claims.Roles, "merchant") {
			c.Locals(localsRecipientType, string(domain.RecipientMerchant))
		} else {
			c.Locals(localsRecipientType, string(domain.RecipientUser))
		}
		return c.Next()
	}
}

func recipientIDFromLocals(c fiber.Ctx) string {
	id, _ := c.Locals(localsRecipientID).(string)
	return id
}

func recipientTypeFromLocals(c fiber.Ctx) domain.RecipientType {
	t, _ := c.Locals(localsRecipientType).(string)
	return domain.RecipientType(t)
}
