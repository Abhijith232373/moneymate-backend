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
	// api.Post("/login", proxy.AuthProxy(authAddr, "/auth/login"))
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
	admin.Get("/merchants", proxy.HTTPProxy(registry, "merchant", "/admin/merchants"))
	admin.Post("/refresh", proxy.AuthProxy(authAddr, "/auth/refresh"))
		// admin user management
		adminUser:=admin.Group("/users")
			adminUser.Post("/",proxy.AuthProxy(authAddr,"/admin/users"))
			adminUser.Get("/",proxy.AuthProxy(authAddr,"/admin/users"))
			adminUser.Get("/:id",proxy.AuthProxy(authAddr,"/admin/users/:id"))
			adminUser.Put("/:id",proxy.AuthProxy(authAddr,"/admin/users/:id"))
			adminUser.Patch("/:id/status",proxy.AuthProxy(authAddr,"/admin/users/:id/status"))
			adminUser.Delete("/:id",proxy.AuthProxy(authAddr,"/admin/users/:id"))
		//admin role management
		adminRole:=admin.Group("/roles")
			adminRole.Post("/",proxy.AuthProxy(authAddr,"/admin/roles"))//create role
			adminRole.Get("/",proxy.AuthProxy(authAddr,"/admin/roles"))//list roles
			adminRole.Get("/:id",proxy.AuthProxy(authAddr,"/admin/roles/:id"))//get role
			adminRole.Put("/:id",proxy.AuthProxy(authAddr,"/admin/roles/:id"))//edit role
			adminRole.Delete("/:id",proxy.AuthProxy(authAddr,"/admin/roles/:id"))//delete role
			adminRole.Post("/assign",proxy.AuthProxy(authAddr,"/admin/roles/assign"))//assign role
			adminRole.Delete("/users/:userId/roles/:roleId",proxy.AuthProxy(authAddr,"/admin/roles/users/:userId/roles/:roleId"))//remove role
			adminRole.Get("/users/:userId",proxy.AuthProxy(authAddr,"/admin/roles/users/:userId"))//get user role

		// admin merchant management
		adminMerchant := admin.Group("/merchants")
		adminMerchant.Get("/:id", proxy.HTTPProxy(registry, "merchant", "/admin/merchants/:id"))
		adminMerchant.Put("/:id/status", proxy.HTTPProxy(registry, "merchant", "/admin/merchants/:id/status"))
		adminMerchant.Delete("/:id", proxy.HTTPProxy(registry, "merchant", "/admin/merchants/:id"))

		adminMerchant.Get("/:store_id/campaigns", proxy.HTTPProxy(registry, "merchant", "/admin/merchants/:store_id/campaigns"))
		adminMerchant.Get("/:store_id/kyc", proxy.HTTPProxy(registry, "merchant", "/admin/merchants/:store_id/kyc"))
		adminMerchant.Put("/:store_id/kyc/verify", proxy.HTTPProxy(registry, "merchant", "/admin/merchants/:store_id/kyc/verify"))
		adminMerchant.Put("/:store_id/subscription", proxy.HTTPProxy(registry, "merchant", "/admin/merchants/:store_id/subscription"))


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

	secure.Post("/pin", proxy.AuthProxy(authAddr, "/user/pin"))
	secure.Put("/pin", proxy.AuthProxy(authAddr, "/user/pin"))
	secure.Post("/pin/verify", proxy.AuthProxy(authAddr, "/user/pin/verify"))

	// ── Merchant (authenticated + role=merchant) ───────────────────
	merchant := api.Group("/merchant")
	merchant.Use(authMiddleware)
	merchant.Use(middlewares.RequireRole("merchant"))
	merchant.Get("/dashboard", proxy.HTTPProxy(registry, "merchant", "/merchant/dashboard"))
	merchant.Get("/:store_id/dashboard", proxy.HTTPProxy(registry, "merchant", "/merchant/:store_id/dashboard"))

	merchant.Post("/register", proxy.HTTPProxy(registry, "merchant", "/merchant/register"))
	merchant.Get("/status/:owner_id", proxy.HTTPProxy(registry, "merchant", "/merchant/status/:owner_id"))
	merchant.Get("/pending", proxy.HTTPProxy(registry, "merchant", "/merchant/pending"))

	merchant.Get("/profile", proxy.HTTPProxy(registry, "merchant", "/merchant/profile"))
	merchant.Get("/:store_id/profile", proxy.HTTPProxy(registry, "merchant", "/merchant/:store_id/profile"))
	merchant.Put("/profile", proxy.HTTPProxy(registry, "merchant", "/merchant/profile"))
	merchant.Put("/:store_id/profile", proxy.HTTPProxy(registry, "merchant", "/merchant/:store_id/profile"))
	merchant.Post("/profile", proxy.HTTPProxy(registry, "merchant", "/merchant/profile"))
	merchant.Post("/:store_id/profile", proxy.HTTPProxy(registry, "merchant", "/merchant/:store_id/profile"))

	merchant.Post("/campaigns", proxy.HTTPProxy(registry, "merchant", "/merchant/campaigns"))
	merchant.Post("/:store_id/campaigns", proxy.HTTPProxy(registry, "merchant", "/merchant/:store_id/campaigns"))
	merchant.Get("/campaigns", proxy.HTTPProxy(registry, "merchant", "/merchant/campaigns"))
	merchant.Get("/:store_id/campaigns", proxy.HTTPProxy(registry, "merchant", "/merchant/:store_id/campaigns"))

	merchant.Get("/rewards/summary", proxy.HTTPProxy(registry, "merchant", "/merchant/rewards/summary"))
	merchant.Get("/:store_id/rewards/summary", proxy.HTTPProxy(registry, "merchant", "/merchant/:store_id/rewards/summary"))
	merchant.Get("/rewards/history", proxy.HTTPProxy(registry, "merchant", "/merchant/rewards/history"))
	merchant.Get("/:store_id/rewards/history", proxy.HTTPProxy(registry, "merchant", "/merchant/:store_id/rewards/history"))
	merchant.Post("/rewards/redeem", proxy.HTTPProxy(registry, "merchant", "/merchant/rewards/redeem"))
	merchant.Post("/:store_id/rewards/redeem", proxy.HTTPProxy(registry, "merchant", "/merchant/:store_id/rewards/redeem"))

	merchant.Get("/subscriptions/plans", proxy.HTTPProxy(registry, "merchant", "/merchant/subscriptions/plans"))
	merchant.Get("/:store_id/subscriptions/plans", proxy.HTTPProxy(registry, "merchant", "/merchant/:store_id/subscriptions/plans"))
	merchant.Get("/subscriptions/current", proxy.HTTPProxy(registry, "merchant", "/merchant/subscriptions/current"))
	merchant.Get("/:store_id/subscriptions/current", proxy.HTTPProxy(registry, "merchant", "/merchant/:store_id/subscriptions/current"))
	merchant.Post("/subscriptions/change", proxy.HTTPProxy(registry, "merchant", "/merchant/subscriptions/change"))
	merchant.Post("/:store_id/subscriptions/change", proxy.HTTPProxy(registry, "merchant", "/merchant/:store_id/subscriptions/change"))
	merchant.Post("/subscriptions/upgrade", proxy.HTTPProxy(registry, "merchant", "/merchant/subscriptions/upgrade"))
	merchant.Post("/:store_id/subscriptions/upgrade", proxy.HTTPProxy(registry, "merchant", "/merchant/:store_id/subscriptions/upgrade"))

	merchant.Get("/kyc", proxy.HTTPProxy(registry, "merchant", "/merchant/kyc"))
	merchant.Get("/:store_id/kyc", proxy.HTTPProxy(registry, "merchant", "/merchant/:store_id/kyc"))
	merchant.Get("/kyc/status", proxy.HTTPProxy(registry, "merchant", "/merchant/kyc/status"))
	merchant.Get("/:store_id/kyc/status", proxy.HTTPProxy(registry, "merchant", "/merchant/:store_id/kyc/status"))
	merchant.Put("/kyc", proxy.HTTPProxy(registry, "merchant", "/merchant/kyc"))
	merchant.Put("/:store_id/kyc", proxy.HTTPProxy(registry, "merchant", "/merchant/:store_id/kyc"))
	merchant.Post("/kyc/update", proxy.HTTPProxy(registry, "merchant", "/merchant/kyc/update"))
	merchant.Post("/:store_id/kyc/update", proxy.HTTPProxy(registry, "merchant", "/merchant/:store_id/kyc/update"))


	// ── Downstream service proxies ─────────────────────────────────
	downstreamServices := []string{"payment", "campaign", "debt", "pod", "scheduler", "referral", "rewards", "routing", "notification"}
	for _, svc := range downstreamServices {
		svcName := svc
		api.All(fmt.Sprintf("/%s/*", svcName), authMiddleware, proxy.ProxyToService(registry, svcName))
	}

	ws.RegisterWebSocketRoutes(app, hub, authClient)
}