package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/moneymate-2026/moneymate-backend/gateway/internal/proxy"
)

func registerSecureRoutes(api fiber.Router, authMiddleware fiber.Handler, authAddr string) {
	secure := api.Group("/secure")
	secure.Use(authMiddleware)
	secure.Get("/profile", func(c fiber.Ctx) error {
		userID := c.Locals("user_id").(string)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    fiber.Map{"user_id": userID},
		})
	})
	secure.Post("/pin", proxy.AuthProxy(authAddr, "/user/pin"))
	secure.Put("/pin", proxy.AuthProxy(authAddr, "/user/pin"))
	secure.Post("/pin/verify", proxy.AuthProxy(authAddr, "/user/pin/verify"))
}

func registerDownstreamRoutes(api fiber.Router, authMiddleware fiber.Handler, registry *proxy.ServiceRegistry) {
	downstreamServices := []string{"payment", "campaign", "debt", "pod", "scheduler", "referral", "rewards", "routing", "notification"}
	for _, svc := range downstreamServices {
		svcName := svc
		api.All("/"+svcName+"/*", authMiddleware, proxy.ProxyToService(registry, svcName))
	}
}