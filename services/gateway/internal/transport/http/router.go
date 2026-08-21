

package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"github.com/moneymate-2026/moneymate-backend/gateway/internal/middlewares"
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
	Redis          *redis.Client
}

func RegisterRoutes(cfg RouteConfig) {
	api := cfg.App.Group("/api/v1")

	api.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok", "service": "gateway"})
	})

	registerAuthRoutes(api.Group("/", middlewares.MaintenanceMode(cfg.Redis, "auth_routes")), cfg.AuthAddr, cfg.AuthMiddleware)
	registerPinRoutes(api.Group("/", middlewares.MaintenanceMode(cfg.Redis, "pin_routes")), cfg.AuthAddr, cfg.AuthMiddleware)
	registerProfilePictureRoutes(api.Group("/", middlewares.MaintenanceMode(cfg.Redis, "profile_picture_routes")), cfg.AuthAddr, cfg.AuthMiddleware)
	registerAdminRoutes(api.Group("/", middlewares.MaintenanceMode(cfg.Redis, "admin_routes")), cfg.AuthMiddleware, cfg.AuthAddr, cfg.MerchantAddr, cfg.Registry, cfg.Redis)
	registerMerchantRoutes(api.Group("/", middlewares.MaintenanceMode(cfg.Redis, "merchant_routes")), cfg.MerchantAddr)
	registerPaymentRoutes(api.Group("/", middlewares.MaintenanceMode(cfg.Redis, "payment_routes")), cfg.AuthMiddleware, cfg.Registry)
	registerSecureRoutes(api.Group("/", middlewares.MaintenanceMode(cfg.Redis, "secure_routes")), cfg.AuthMiddleware, cfg.AuthAddr)
	registerSupportRoutes(api.Group("/", middlewares.MaintenanceMode(cfg.Redis, "support_routes")), cfg.AuthMiddleware, cfg.Registry)
	registerDownstreamRoutes(api.Group("/", middlewares.MaintenanceMode(cfg.Redis, "downstream_routes")), cfg.AuthMiddleware, cfg.Registry)
  registerRewardRoutes(api.Group("/", middlewares.MaintenanceMode(cfg.Redis, "reward_routes")), cfg.AuthMiddleware, cfg.Registry)
  
	// Backward compatibility for frontend calls without /api/v1 prefix
	registerAuthRoutes(cfg.App.Group("/", middlewares.MaintenanceMode(cfg.Redis, "auth_routes")), cfg.AuthAddr, cfg.AuthMiddleware)
	registerPinRoutes(cfg.App.Group("/", middlewares.MaintenanceMode(cfg.Redis, "pin_routes")), cfg.AuthAddr, cfg.AuthMiddleware)
	registerAdminRoutes(cfg.App.Group("/", middlewares.MaintenanceMode(cfg.Redis, "admin_routes")), cfg.AuthMiddleware, cfg.AuthAddr, cfg.MerchantAddr, cfg.Registry, cfg.Redis)
	registerMerchantRoutes(cfg.App.Group("/", middlewares.MaintenanceMode(cfg.Redis, "merchant_routes")), cfg.MerchantAddr)
	registerPaymentRoutes(cfg.App.Group("/", middlewares.MaintenanceMode(cfg.Redis, "payment_routes")), cfg.AuthMiddleware, cfg.Registry)
	registerSecureRoutes(cfg.App.Group("/", middlewares.MaintenanceMode(cfg.Redis, "secure_routes")), cfg.AuthMiddleware, cfg.AuthAddr)
	registerSupportRoutes(cfg.App.Group("/", middlewares.MaintenanceMode(cfg.Redis, "support_routes")), cfg.AuthMiddleware, cfg.Registry)
	registerDownstreamRoutes(cfg.App.Group("/", middlewares.MaintenanceMode(cfg.Redis, "downstream_routes")), cfg.AuthMiddleware, cfg.Registry)
  registerRewardRoutes(cfg.App.Group("/", middlewares.MaintenanceMode(cfg.Redis, "reward_routes")), cfg.AuthMiddleware, cfg.Registry)
  
	ws.RegisterWebSocketRoutes(cfg.App, cfg.Hub, cfg.AuthClient)
}



