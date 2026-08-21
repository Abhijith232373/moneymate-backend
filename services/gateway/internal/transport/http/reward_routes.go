package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/moneymate-2026/moneymate-backend/gateway/internal/proxy"
)

func registerRewardRoutes(api fiber.Router, authMiddleware fiber.Handler, registry *proxy.ServiceRegistry) {
	rewards := api.Group("/rewards")
	rewards.Use(authMiddleware)
	rewards.Get("/me", proxy.HTTPProxy(registry, "rewards", "/rewards/me"))
	rewards.Get("/", proxy.HTTPProxy(registry, "rewards", "/rewards"))
}
