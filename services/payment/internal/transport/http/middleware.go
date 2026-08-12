package http

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	authclient "github.com/moneymate-2026/moneymate-backend/services/payment/internal/adapter/authClient"
	sharedjwt "github.com/moneymate-2026/moneymate-backend/shared/pkg/jwt"
	response "github.com/moneymate-2026/moneymate-backend/shared/pkg/responses"
)

const (
	localsUserID = "userID"
)

// RequireUserID verifies the standard access token (same one auth-svc issues
// on login) and puts the caller's user_id in fiber locals. Applied to every
// /payment route — there is no anonymous access to wallet data.
//
// Mirrors the RequireUserID middleware already used in auth-svc / merchant-svc;
// if the gateway already verifies and forwards a trusted header (e.g.
// X-User-Id) once S1's gateway routing lands, swap the body of this function
// to read that header instead of re-parsing the JWT — the rest of the service
// doesn't change.
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

func RequireTransactionToken(authClient *authclient.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Get("X-Transaction-Token")
		if token == "" {
			return response.Unauthorized(c, "transaction token required")
		}
		verifiedUserID, err := authClient.VerifyTransactionToken(c.Context(), token, "")
		if err != nil {
			return response.Unauthorized(c, "invalid, expired, or already-used transaction token")
		}

		authUserID, _ := c.Locals(localsUserID).(string)
		if verifiedUserID != authUserID {
			return response.Forbidden(c, nil, "transaction token does not match authenticated user")
		}
		return c.Next()
	}
}

func userIDFromLocals(c fiber.Ctx) string {
	id, _ := c.Locals(localsUserID).(string)
	return id
}

func RequireInternalSecret(secret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		provided := c.Get("X-Internal-Secret")
		if provided == "" || provided != secret {
			return response.Unauthorized(c, "internal access required")
		}
		return c.Next()
	}
}