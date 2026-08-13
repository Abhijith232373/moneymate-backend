// package http

// import (
// 	"fmt"

// 	"github.com/gofiber/fiber/v3"
// 	"github.com/moneymate-2026/moneymate-backend/gateway/internal/middlewares"
// 	"github.com/moneymate-2026/moneymate-backend/gateway/internal/proxy"
// 	ws "github.com/moneymate-2026/moneymate-backend/gateway/internal/websocket"
// )

// func RegisterRoutes(
// 	app *fiber.App,
// 	authMiddleware fiber.Handler,
// 	authClient proxy.AuthClient,
// 	registry *proxy.ServiceRegistry,
// 	hub *ws.Hub,
// 	authAddr string,
// 	merchantAddr string,
// ) {
// 	api := app.Group("/api/v1")

// 	api.Get("/health", func(c fiber.Ctx) error {
// 		return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 			"status":  "ok",
// 			"service": "gateway",
// 		})
// 	})
	

// 	// ── User Auth ──────────────────────────────────────────────────
// 	userAuth := api.Group("/auth")
// 	userAuth.Get("/health",  proxy.AuthProxy(authAddr, "/auth/health"))
// 	userAuth.Post("/register", proxy.AuthProxy(authAddr, "/auth/user/register"))
// 	userAuth.Post("/login", proxy.AuthProxy(authAddr, "/auth/login"))
// 	userAuth.Post("/admin/login", proxy.AuthProxy(authAddr, "/auth/admin/login"))
// 	userAuth.Post("/logout", authMiddleware, proxy.AuthProxy(authAddr, "/auth/logout"))
// 	userAuth.Post("/otp/send", proxy.AuthProxy(authAddr, "/auth/otp/send"))
// 	userAuth.Post("/otp/verify", proxy.AuthProxy(authAddr, "/auth/otp/verify"))
// 	userAuth.Post("/refresh", proxy.AuthProxy(authAddr, "/auth/refresh"))

// 	// ── User PIN ──────────────────────────────────────────────────
// 	userPin:=api.Group("/pin")
// 		userPin.Post("/",  proxy.AuthProxy(authAddr, "/user/pin"))
// 		userPin.Put("/",  proxy.AuthProxy(authAddr, "/user/pin"))
// 		userPin.Post("/verify",  proxy.AuthProxy(authAddr, "/user/pin"))

// 	// ── Merchant Auth ──────────────────────────────────────────────
// 	merchantAuth := api.Group("/merchant/auth")
// 	merchantAuth.Post("/register", proxy.AuthProxy(authAddr, "/auth/merchant/register"))
// 	merchantAuth.Post("/login", proxy.AuthProxy(authAddr, "/auth/login"))
// 	merchantAuth.Post("/logout", authMiddleware, proxy.AuthProxy(authAddr, "/auth/logout"))
// 	merchantAuth.Post("/refresh", proxy.AuthProxy(authAddr, "/auth/refresh"))




// 	// ── Admin ─────────────────────────────────────────────────
// 	// api.Post("/login", proxy.AuthProxy(authAddr, "/auth/login"))
// 	admin := api.Group("/admin")
// 	admin.Use(authMiddleware)
// 	admin.Use(middlewares.RequireRole("admin"))

// 	admin.Get("/dashboard", func(c fiber.Ctx) error {
// 		return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 			"success": true,
// 			"message": "admin dashboard data",
// 			"data": fiber.Map{
// 				"total_users":        0,
// 				"total_merchants":    0,
// 				"pending_reviews":    0,
// 				"total_transactions": 0,
// 			},
// 		})
// 	})
// 	admin.Post("/refresh", proxy.AuthProxy(authAddr, "/auth/refresh"))

// 		// admin user management
// 		adminUser:=admin.Group("/users")
// 			adminUser.Post("/",proxy.AuthProxy(authAddr,"/admin/users"))
// 			adminUser.Get("/",proxy.AuthProxy(authAddr,"/admin/users"))
// 			adminUser.Get("/:id",proxy.AuthProxy(authAddr,"/admin/users/:id"))
// 			adminUser.Put("/:id",proxy.AuthProxy(authAddr,"/admin/users/:id"))
// 			adminUser.Patch("/:id/status",proxy.AuthProxy(authAddr,"/admin/users/:id/status"))
// 			adminUser.Delete("/:id",proxy.AuthProxy(authAddr,"/admin/users/:id"))
			
// 		//admin role management
// 		adminRole:=admin.Group("/roles")
// 			adminRole.Post("/",proxy.AuthProxy(authAddr,"/admin/roles"))//create role
// 			adminRole.Get("/",proxy.AuthProxy(authAddr,"/admin/roles"))//list roles
// 			adminRole.Get("/:id",proxy.AuthProxy(authAddr,"/admin/roles/:id"))//get role
// 			adminRole.Put("/:id",proxy.AuthProxy(authAddr,"/admin/roles/:id"))//edit role
// 			adminRole.Delete("/:id",proxy.AuthProxy(authAddr,"/admin/roles/:id"))//delete role
// 			adminRole.Post("/assign",proxy.AuthProxy(authAddr,"/admin/roles/assign"))//assign role
// 			adminRole.Delete("/users/:userId/roles/:roleId",proxy.AuthProxy(authAddr,"/admin/roles/users/:userId/roles/:roleId"))//remove role
// 			adminRole.Get("/users/:userId",proxy.AuthProxy(authAddr,"/admin/roles/users/:userId"))//get user role
		
			
// 		// admin permission management
// 		adminPermission := admin.Group("/permissions")
// 			adminPermission.Post("/", proxy.AuthProxy(authAddr, "/admin/permissions"))                                          
// 			adminPermission.Get("/", proxy.AuthProxy(authAddr, "/admin/permissions"))                                           
// 			adminPermission.Get("/:id", proxy.AuthProxy(authAddr, "/admin/permissions/:id"))                                   
// 			adminPermission.Delete("/:id", proxy.AuthProxy(authAddr, "/admin/permissions/:id"))                                 
// 			adminPermission.Post("/assign", proxy.AuthProxy(authAddr, "/admin/permissions/assign"))                              
// 			adminPermission.Delete("/roles/:roleId/permissions/:permissionId", proxy.AuthProxy(authAddr, "/admin/permissions/roles/:roleId/permissions/:permissionId"))
// 			adminPermission.Get("/roles/:roleId", proxy.AuthProxy(authAddr, "/admin/permissions/roles/:roleId"))                  
	
			
// 		// admin merchant management
// 		adminMerchant := admin.Group("/merchants")
// 			adminMerchant.Get("/:id", proxy.MerchantProxy(merchantAddr, "/admin/merchants/:id"))
// 			adminMerchant.Put("/:id/status", proxy.MerchantProxy(merchantAddr, "/admin/merchants/:id/status"))
// 			adminMerchant.Delete("/:id", proxy.MerchantProxy(merchantAddr, "/admin/merchants/:id"))

// 			adminMerchant.Get("/:store_id/campaigns", proxy.MerchantProxy(merchantAddr, "/admin/merchants/:store_id/campaigns"))
// 			adminMerchant.Post("/:store_id/campaigns", proxy.MerchantProxy(merchantAddr, "/admin/merchants/:store_id/campaigns"))
// 			adminMerchant.Get("/:store_id/kyc", proxy.MerchantProxy(merchantAddr, "/admin/merchants/:store_id/kyc"))
// 			adminMerchant.Put("/:store_id/kyc/verify", proxy.MerchantProxy(merchantAddr, "/admin/merchants/:store_id/kyc/verify"))
// 			adminMerchant.Put("/:store_id/subscription", proxy.MerchantProxy(merchantAddr, "/admin/merchants/:store_id/subscription"))

		
// 		// admin campaigns management
// 		adminCampaign := admin.Group("/campaigns")
// 			adminCampaign.Get("/", proxy.MerchantProxy(merchantAddr, "/admin/campaigns"))
// 			adminCampaign.Put("/:id/status", proxy.MerchantProxy(merchantAddr, "/admin/campaigns/:id/status"))
// 			adminCampaign.Delete("/:id", proxy.MerchantProxy(merchantAddr, "/admin/campaigns/:id"))

// 		// admin kyc management
// 		adminKYC := admin.Group("/kyc")
// 			adminKYC.Get("/", proxy.MerchantProxy(merchantAddr, "/admin/kyc"))
// 			adminKYC.Get("/:store_id", proxy.MerchantProxy(merchantAddr, "/admin/kyc/:store_id"))
// 			adminKYC.Put("/:store_id/verify", proxy.MerchantProxy(merchantAddr, "/admin/kyc/:store_id/verify"))

// 		// admin rewards management
// 		adminRewards := admin.Group("/rewards")
// 			adminRewards.Get("/history", proxy.MerchantProxy(merchantAddr, "/admin/rewards/history"))
// 			adminRewards.Get("/summary", proxy.MerchantProxy(merchantAddr, "/admin/rewards/summary"))

// 		// admin subscriptions management
// 		adminSubscriptions := admin.Group("/subscriptions")
// 			adminSubscriptions.Get("/", proxy.MerchantProxy(merchantAddr, "/admin/subscriptions"))

// 	// ── Secure (authenticated user) ────────────────────────────────
// 	secure := api.Group("/secure")
// 	secure.Use(authMiddleware)
// 	secure.Get("/profile", func(c fiber.Ctx) error {
// 		userID := c.Locals("user_id").(string)
// 		return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 			"success": true,
// 			"data": fiber.Map{
// 				"user_id": userID,
// 			},
// 		})
// 	})

// 	secure.Post("/pin", proxy.AuthProxy(authAddr, "/user/pin"))
// 	secure.Put("/pin", proxy.AuthProxy(authAddr, "/user/pin"))
// 	secure.Post("/pin/verify", proxy.AuthProxy(authAddr, "/user/pin/verify"))

// 	// ── Merchant (Unauthenticated) ─────────────────────────────────
// 	merchantUnauth := api.Group("/merchant")
// 	merchantUnauth.Post("/register", proxy.MerchantProxy(merchantAddr, "/merchant/register"))
// 	merchantUnauth.Post("/login", proxy.MerchantProxy(merchantAddr, "/merchant/login"))
// 	merchantUnauth.Get("/public/campaigns", proxy.MerchantProxy(merchantAddr, "/merchant/public/campaigns"))

// 	// ── Merchant (authenticated + role=merchant/admin) ───────────────────
// 	merchant := api.Group("/merchant")
// 	merchant.Get("/health", proxy.MerchantProxy(merchantAddr, "/merchant/health"))
// 	// The merchant service has its own auth middleware for these routes
// 	// so we bypass the gateway's auth middleware here.
// 	merchant.Get("/dashboard", proxy.MerchantProxy(merchantAddr, "/merchant/dashboard"))
// 	merchant.Get("/:store_id/dashboard", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/dashboard"))

// 	merchant.Get("/status/:owner_id", proxy.MerchantProxy(merchantAddr, "/merchant/status/:owner_id"))
// 	merchant.Get("/pending", proxy.MerchantProxy(merchantAddr, "/merchant/pending"))

// 	merchant.Get("/profile", proxy.MerchantProxy(merchantAddr, "/merchant/profile"))
// 	merchant.Get("/:store_id/profile", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/profile"))
// 	merchant.Put("/profile", proxy.MerchantProxy(merchantAddr, "/merchant/profile"))
// 	merchant.Put("/:store_id/profile", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/profile"))
// 	merchant.Post("/profile", proxy.MerchantProxy(merchantAddr, "/merchant/profile"))
// 	merchant.Post("/:store_id/profile", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/profile"))

// 	merchant.Post("/campaigns", proxy.MerchantProxy(merchantAddr, "/merchant/campaigns"))
// 	merchant.Post("/:store_id/campaigns", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/campaigns"))
// 	merchant.Get("/campaigns", proxy.MerchantProxy(merchantAddr, "/merchant/campaigns"))
// 	merchant.Get("/:store_id/campaigns", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/campaigns"))
// 	merchant.Put("/campaigns/:id/status", proxy.MerchantProxy(merchantAddr, "/merchant/campaigns/:id/status"))
// 	merchant.Put("/:store_id/campaigns/:id/status", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/campaigns/:id/status"))

// 	merchant.Get("/rewards/summary", proxy.MerchantProxy(merchantAddr, "/merchant/rewards/summary"))
// 	merchant.Get("/:store_id/rewards/summary", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/rewards/summary"))
// 	merchant.Get("/rewards/history", proxy.MerchantProxy(merchantAddr, "/merchant/rewards/history"))
// 	merchant.Get("/:store_id/rewards/history", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/rewards/history"))
// 	merchant.Post("/rewards/redeem", proxy.MerchantProxy(merchantAddr, "/merchant/rewards/redeem"))
// 	merchant.Post("/:store_id/rewards/redeem", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/rewards/redeem"))

// 	merchant.Get("/subscriptions/plans", proxy.MerchantProxy(merchantAddr, "/merchant/subscriptions/plans"))
// 	merchant.Get("/:store_id/subscriptions/plans", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/subscriptions/plans"))
// 	merchant.Get("/subscriptions/current", proxy.MerchantProxy(merchantAddr, "/merchant/subscriptions/current"))
// 	merchant.Get("/:store_id/subscriptions/current", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/subscriptions/current"))
// 	merchant.Post("/subscriptions/change", proxy.MerchantProxy(merchantAddr, "/merchant/subscriptions/change"))
// 	merchant.Post("/:store_id/subscriptions/change", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/subscriptions/change"))
// 	merchant.Post("/subscriptions/upgrade/initiate", proxy.MerchantProxy(merchantAddr, "/merchant/subscriptions/upgrade/initiate"))
// 	merchant.Post("/:store_id/subscriptions/upgrade/initiate", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/subscriptions/upgrade/initiate"))
// 	merchant.Post("/subscriptions/upgrade/verify", proxy.MerchantProxy(merchantAddr, "/merchant/subscriptions/upgrade/verify"))
// 	merchant.Post("/:store_id/subscriptions/upgrade/verify", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/subscriptions/upgrade/verify"))

// 	merchant.Get("/wallet", proxy.MerchantProxy(merchantAddr, "/merchant/wallet"))
// 	merchant.Get("/:store_id/wallet", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/wallet"))

// 	merchant.Get("/earnings", proxy.MerchantProxy(merchantAddr, "/merchant/earnings"))
// 	merchant.Get("/:store_id/earnings", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/earnings"))
// 	merchant.Post("/earnings/payouts", proxy.MerchantProxy(merchantAddr, "/merchant/earnings/payouts"))
// 	merchant.Post("/:store_id/earnings/payouts", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/earnings/payouts"))

// 	merchant.Get("/kyc", proxy.MerchantProxy(merchantAddr, "/merchant/kyc"))
// 	merchant.Get("/:store_id/kyc", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/kyc"))
// 	merchant.Get("/kyc/status", proxy.MerchantProxy(merchantAddr, "/merchant/kyc/status"))
// 	merchant.Get("/:store_id/kyc/status", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/kyc/status"))
// 	merchant.Put("/kyc", proxy.MerchantProxy(merchantAddr, "/merchant/kyc"))
// 	merchant.Put("/:store_id/kyc", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/kyc"))
// 	merchant.Post("/kyc/update", proxy.MerchantProxy(merchantAddr, "/merchant/kyc/update"))
// 	merchant.Post("/:store_id/kyc/update", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/kyc/update"))


// 	// ── Downstream service proxies ─────────────────────────────────
// 	downstreamServices := []string{"payment", "campaign", "debt", "pod", "scheduler", "referral", "rewards", "routing", "notification"}
// 	for _, svc := range downstreamServices {
// 		svcName := svc
// 		api.All(fmt.Sprintf("/%s/*", svcName), authMiddleware, proxy.ProxyToService(registry, svcName))
// 	}

// 	ws.RegisterWebSocketRoutes(app, hub, authClient)
// }


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
	registerAdminRoutes(api, cfg.AuthMiddleware, cfg.AuthAddr, cfg.MerchantAddr)
	registerMerchantRoutes(api, cfg.MerchantAddr)
	registerPaymentRoutes(api, cfg.AuthMiddleware, cfg.Registry)
	registerSecureRoutes(api, cfg.AuthMiddleware, cfg.AuthAddr)
	registerDownstreamRoutes(api, cfg.AuthMiddleware, cfg.Registry)

	ws.RegisterWebSocketRoutes(cfg.App, cfg.Hub, cfg.AuthClient)
}