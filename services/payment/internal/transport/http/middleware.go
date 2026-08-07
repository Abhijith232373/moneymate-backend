package http

import (
	"strings"

	"github.com/gofiber/fiber/v3"

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

// RequireTransactionToken verifies the short-lived transaction_token minted
// by POST /auth/pin/verify. It must be present, valid, unexpired, and its
// subject must match the already-authenticated user_id from RequireUserID.
// Apply this ONLY on top of RequireUserID, and ONLY on /payment/transfers.
func RequireTransactionToken(cfg sharedjwt.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Get("X-Transaction-Token")
		if token == "" {
			return response.Unauthorized(c, "transaction token required")
		}
		claims, err := sharedjwt.ParseTransactionToken(token, cfg.AccessSecret)
		if err != nil {
			return response.Unauthorized(c, "invalid or expired transaction token")
		}
		userID, _ := c.Locals(localsUserID).(string)
		if claims.UserID != userID {
			return response.Forbidden(c, nil, "transaction token does not match authenticated user")
		}
		return c.Next()
	}
}

func userIDFromLocals(c fiber.Ctx) string {
	id, _ := c.Locals(localsUserID).(string)
	return id
}