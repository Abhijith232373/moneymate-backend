package http

import (
	"fmt"

	"github.com/gofiber/fiber/v3"

	"github.com/moneymate-2026/moneymate-backend/gateway/internal/middlewares"
	"github.com/moneymate-2026/moneymate-backend/gateway/internal/proxy"
	ws "github.com/moneymate-2026/moneymate-backend/gateway/internal/websocket"
)

func RegisterRoutes(
	app *fiber.App,
	authMiddleware fiber.Handler,
	authClient proxy.AuthClient,
	registry *proxy.ServiceRegistry,
	hub *ws.Hub,
	authAddr string,
) {
	api := app.Group("/api/v1")

	api.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"service": "gateway",
		})
	})

	// ── User Auth ──────────────────────────────────────────────────
	userAuth := api.Group("/auth")
	userAuth.Post("/register", proxy.AuthProxy(authAddr, "/auth/user/register"))
	userAuth.Post("/login", proxy.AuthProxy(authAddr, "/auth/login"))
	userAuth.Post("/logout", authMiddleware, proxy.AuthProxy(authAddr, "/auth/logout"))
	userAuth.Post("/otp/send", proxy.AuthProxy(authAddr, "/auth/otp/send"))
	userAuth.Post("/otp/verify", proxy.AuthProxy(authAddr, "/auth/otp/verify"))
	userAuth.Post("/refresh", proxy.AuthProxy(authAddr, "/auth/refresh"))

	// ── Merchant Auth ──────────────────────────────────────────────
	merchantAuth := api.Group("/merchant/auth")
	merchantAuth.Post("/register", proxy.AuthProxy(authAddr, "/auth/merchant/register"))
	merchantAuth.Post("/login", proxy.AuthProxy(authAddr, "/auth/login"))
	merchantAuth.Post("/logout", authMiddleware, proxy.AuthProxy(authAddr, "/auth/logout"))
	merchantAuth.Post("/refresh", proxy.AuthProxy(authAddr, "/auth/refresh"))

	// ── Admin ─────────────────────────────────────────────────
	adminRoutes := api.Group("/admin")
	adminRoutes.Post("/login", proxy.AuthProxy(authAddr, "/auth/login"))
	adminRoutes.Post("/refresh", proxy.AuthProxy(authAddr, "/auth/refresh"))
		// admin user management
		adminUser:=adminRoutes.Group("/users")
			adminRoutes.Post("/",proxy.AuthProxy(authAddr,"/admin/users"))
			adminUser.Get("/",proxy.AuthProxy(authAddr,"/admin/users"))
			adminUser.Get("/:id",proxy.AuthProxy(authAddr,"/admin/users/:id"))
			adminUser.Put("/:id",proxy.AuthProxy(authAddr,"/admin/users/:id"))
			adminUser.Patch("/:id/status",proxy.AuthProxy(authAddr,"/admin/users/:id/status"))
			adminUser.Delete("/:id",proxy.AuthProxy(authAddr,"/admin/users/:id"))


	// ── Secure (authenticated user) ────────────────────────────────
	secure := api.Group("/secure")
	secure.Use(authMiddleware)
	secure.Get("/profile", func(c fiber.Ctx) error {
		userID := c.Locals("user_id").(string)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"user_id": userID,
			},
		})
	})

	// ── Merchant (authenticated + role=merchant) ───────────────────
	merchant := api.Group("/merchant")
	merchant.Use(authMiddleware)
	merchant.Use(middlewares.RequireRole("merchant"))
	merchant.Get("/dashboard", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "merchant dashboard placeholder",
		})
	})

	// ── Admin (authenticated + role=admin) ─────────────────────────
	admin := api.Group("/admin")
	admin.Use(authMiddleware)
	admin.Use(middlewares.RequireRole("admin"))

	admin.Get("/dashboard", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "admin dashboard data",
			"data": fiber.Map{
				"total_users":        0,
				"total_merchants":    0,
				"pending_reviews":    0,
				"total_transactions": 0,
			},
		})
	})

	admin.Get("/merchants", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "merchant list placeholder",
			"data":    []interface{}{},
		})
	})

	// ── Downstream service proxies ─────────────────────────────────
	downstreamServices := []string{"payment", "merchant", "campaign", "debt", "pod", "scheduler", "referral", "rewards", "routing", "notification"}
	for _, svc := range downstreamServices {
		svcName := svc
		api.All(fmt.Sprintf("/%s/*", svcName), authMiddleware, proxy.ProxyToService(registry, svcName))
	}

	ws.RegisterWebSocketRoutes(app, hub, authClient)
}