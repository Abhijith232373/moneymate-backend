package middlewares

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

func MaintenanceMode(rdb *redis.Client, moduleKey string) fiber.Handler {
	return func(c fiber.Ctx) error {
		// By default, assume enabled unless explicitly set to false
		val, err := rdb.Get(context.Background(), "config:module:"+moduleKey).Result()
		if err == nil && val == "false" {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"success": false,
				"error":   "503 Maintenance Mode",
				"message": "This module is temporarily unavailable due to system maintenance.",
			})
		}
		return c.Next()
	}
}
