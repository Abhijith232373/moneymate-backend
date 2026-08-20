package middlewares

import (
	"github.com/gofiber/fiber/v3"
)

func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok || role == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "access denied: no role found",
			})
		}

		for _, allowed := range allowedRoles {
			if role == allowed {
				return c.Next()
			}
		}

		// Allow any role that isn't 'user' or 'merchant' to access admin areas
		// if the allowed role requirement contains 'admin'
		for _, allowed := range allowedRoles {
			if allowed == "admin" && role != "user" && role != "merchant" {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "access denied: insufficient permissions",
		})
	}
}