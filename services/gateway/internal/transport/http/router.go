

package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/moneymate-2026/moneymate-backend/gateway/internal/proxy"
	ws "github.com/moneymate-2026/moneymate-backend/gateway/internal/websocket"
)

type RouteConfig struct {
	App            *fiber.App
	AuthMiddleware fiber.Handler
	AuthClient     proxy.AuthClient
	Registry       *proxy.ServiceRegistry
	Hub            *ws.Hub
	AuthAddr       string
	MerchantAddr   string
}

func RegisterRoutes(cfg RouteConfig) {
	api := cfg.App.Group("/api/v1")

	api.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok", "service": "gateway"})
	})

	registerAuthRoutes(api, cfg.AuthAddr, cfg.AuthMiddleware)
	registerPinRoutes(api, cfg.AuthAddr, cfg.AuthMiddleware)
	registerAdminRoutes(api, cfg.AuthMiddleware, cfg.AuthAddr, cfg.MerchantAddr, cfg.Registry)
	registerMerchantRoutes(api, cfg.MerchantAddr)
	registerPaymentRoutes(api, cfg.AuthMiddleware, cfg.Registry)
	registerSecureRoutes(api, cfg.AuthMiddleware, cfg.AuthAddr)
	registerSupportRoutes(api, cfg.AuthMiddleware, cfg.Registry)
	registerDownstreamRoutes(api, cfg.AuthMiddleware, cfg.Registry)

	
	ws.RegisterWebSocketRoutes(cfg.App, cfg.Hub, cfg.AuthClient)
}



