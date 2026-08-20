package http

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/moneymate-2026/moneymate-backend/gateway/internal/middlewares"
	"github.com/moneymate-2026/moneymate-backend/gateway/internal/proxy"
	"github.com/redis/go-redis/v9"
)

func registerAdminRoutes(api fiber.Router, authMiddleware fiber.Handler, authAddr, merchantAddr string, registry *proxy.ServiceRegistry, rdb *redis.Client) {
	admin := api.Group("/admin")
	
	admin.Use(authMiddleware)
	admin.Use(middlewares.RequireRole("admin"))
	// admin.Post("/login", proxy.AuthProxy(authAddr, "/admin/login"))

	admin.Get("/dashboard/stats", proxy.MerchantProxy(merchantAddr, "/admin/dashboard/stats"))
	admin.Post("/refresh", proxy.AuthProxy(authAddr, "/auth/refresh"))

	admin.Get("/config", func(c fiber.Ctx) error {
		ctx := context.Background()
		modules := []string{
			// Master modules
			"auth_routes", "pin_routes", "admin_routes", "merchant_routes", 
			"payment_routes", "secure_routes", "support_routes", "downstream_routes",
			
			// Auth sub-routes
			"auth_user_login", "auth_user_register", "auth_admin_login", "auth_merchant_login", "auth_merchant_register",
			
			// PIN sub-routes
			"pin_management", "pin_verification",
			
			// Admin sub-routes
			"admin_merchants_kyc", "admin_master_campaigns", "admin_platform_config", "admin_system_audit",
			
			// Merchant sub-routes
			"merch_dashboard_analytics", "merch_campaigns_management", "merch_subscriptions", "merch_wallet_payouts",
			
			// Payment sub-routes
			"pay_p2p_transfers", "pay_fiat_deposits", "pay_fiat_withdrawals", "pay_wallet_balances",
			
			// Secure sub-routes
			"secure_service_comm", "secure_identity_sync",
			
			// Support sub-routes
			"support_user_complaints", "support_fraud_reports", "support_live_chat",
			
			// Downstream sub-routes
			"ds_payment_aggregator", "ds_campaign_service", "ds_notification_engine",
		}
		states := make(map[string]bool)
		for _, m := range modules {
			val, _ := rdb.Get(ctx, "config:module:"+m).Result()
			states[m] = (val != "false")
		}
		return c.JSON(fiber.Map{"success": true, "data": states})
	})

	admin.Put("/config", func(c fiber.Ctx) error {
		var req map[string]bool
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"success": false, "error": "invalid payload"})
		}
		ctx := context.Background()
		for k, v := range req {
			if v {
				rdb.Set(ctx, "config:module:"+k, "true", 0)
			} else {
				rdb.Set(ctx, "config:module:"+k, "false", 0)
			}
		}
		return c.JSON(fiber.Map{"success": true, "message": "config updated"})
	})

	admin.Get("/audit", proxy.HTTPProxy(registry, "support", "/admin/support/audit-logs"))
	admin.Post("/audit", proxy.HTTPProxy(registry, "support", "/admin/support/audit-logs"))

	registerAdminUserRoutes(admin, authAddr)
	registerAdminStaffRoutes(admin, authAddr)
	registerAdminRoleRoutes(admin, authAddr)
	registerAdminPermissionRoutes(admin, authAddr)
	registerAdminMerchantRoutes(admin, merchantAddr)
	registerAdminCampaignRoutes(admin, merchantAddr)
	registerAdminKYCRoutes(admin, merchantAddr)
	registerAdminRewardsRoutes(admin, merchantAddr)
	registerAdminSubscriptionRoutes(admin, merchantAddr)
	registerAdminSupportRoutes(admin, registry)
}

func registerAdminStaffRoutes(admin fiber.Router, authAddr string) {
	adminStaff := admin.Group("/staff")
	adminStaff.Post("", proxy.AuthProxy(authAddr, "/admin/staff"))
	adminStaff.Post("/", proxy.AuthProxy(authAddr, "/admin/staff"))
	adminStaff.Get("", proxy.AuthProxy(authAddr, "/admin/staff"))
	adminStaff.Get("/", proxy.AuthProxy(authAddr, "/admin/staff"))
	adminStaff.Get("/:id", proxy.AuthProxy(authAddr, "/admin/staff/:id"))
	adminStaff.Put("/:id", proxy.AuthProxy(authAddr, "/admin/staff/:id"))
	adminStaff.Patch("/:id/status", proxy.AuthProxy(authAddr, "/admin/staff/:id/status"))
	adminStaff.Delete("/:id", proxy.AuthProxy(authAddr, "/admin/staff/:id"))
}

func registerAdminUserRoutes(admin fiber.Router, authAddr string) {
	adminUser := admin.Group("/users")
	adminUser.Post("", proxy.AuthProxy(authAddr, "/admin/users"))
	adminUser.Post("/", proxy.AuthProxy(authAddr, "/admin/users"))
	adminUser.Get("", proxy.AuthProxy(authAddr, "/admin/users"))
	adminUser.Get("/", proxy.AuthProxy(authAddr, "/admin/users"))
	adminUser.Get("/:id", proxy.AuthProxy(authAddr, "/admin/users/:id"))
	adminUser.Put("/:id", proxy.AuthProxy(authAddr, "/admin/users/:id"))
	adminUser.Patch("/:id/status", proxy.AuthProxy(authAddr, "/admin/users/:id/status"))
	adminUser.Delete("/:id", proxy.AuthProxy(authAddr, "/admin/users/:id"))
}

func registerAdminRoleRoutes(admin fiber.Router, authAddr string) {
	adminRole := admin.Group("/roles")
	adminRole.Post("", proxy.AuthProxy(authAddr, "/admin/roles"))
	adminRole.Post("/", proxy.AuthProxy(authAddr, "/admin/roles"))
	adminRole.Get("", proxy.AuthProxy(authAddr, "/admin/roles"))
	adminRole.Get("/", proxy.AuthProxy(authAddr, "/admin/roles"))
	adminRole.Get("/:id", proxy.AuthProxy(authAddr, "/admin/roles/:id"))
	adminRole.Put("/:id", proxy.AuthProxy(authAddr, "/admin/roles/:id"))
	adminRole.Delete("/:id", proxy.AuthProxy(authAddr, "/admin/roles/:id"))
	adminRole.Post("/assign", proxy.AuthProxy(authAddr, "/admin/roles/assign"))
	adminRole.Delete("/users/:userId/roles/:roleId", proxy.AuthProxy(authAddr, "/admin/roles/users/:userId/roles/:roleId"))
	adminRole.Get("/users/:userId", proxy.AuthProxy(authAddr, "/admin/roles/users/:userId"))
}

func registerAdminPermissionRoutes(admin fiber.Router, authAddr string) {
	adminPermission := admin.Group("/permissions")
	adminPermission.Post("", proxy.AuthProxy(authAddr, "/admin/permissions"))
	adminPermission.Post("/", proxy.AuthProxy(authAddr, "/admin/permissions"))
	adminPermission.Get("", proxy.AuthProxy(authAddr, "/admin/permissions"))
	adminPermission.Get("/", proxy.AuthProxy(authAddr, "/admin/permissions"))
	adminPermission.Get("/:id", proxy.AuthProxy(authAddr, "/admin/permissions/:id"))
	adminPermission.Delete("/:id", proxy.AuthProxy(authAddr, "/admin/permissions/:id"))
	adminPermission.Post("/assign", proxy.AuthProxy(authAddr, "/admin/permissions/assign"))
	adminPermission.Delete("/roles/:roleId/permissions/:permissionId", proxy.AuthProxy(authAddr, "/admin/permissions/roles/:roleId/permissions/:permissionId"))
	adminPermission.Get("/roles/:roleId", proxy.AuthProxy(authAddr, "/admin/permissions/roles/:roleId"))
}

func registerAdminMerchantRoutes(admin fiber.Router, merchantAddr string) {
	adminMerchant := admin.Group("/merchants")
	adminMerchant.Get("", proxy.MerchantProxy(merchantAddr, "/admin/merchants"))
	adminMerchant.Get("/", proxy.MerchantProxy(merchantAddr, "/admin/merchants"))
	adminMerchant.Get("/:id", proxy.MerchantProxy(merchantAddr, "/admin/merchants/:id"))
	adminMerchant.Put("/:id/status", proxy.MerchantProxy(merchantAddr, "/admin/merchants/:id/status"))
	adminMerchant.Delete("/:id", proxy.MerchantProxy(merchantAddr, "/admin/merchants/:id"))
	adminMerchant.Get("/:store_id/campaigns", proxy.MerchantProxy(merchantAddr, "/admin/merchants/:store_id/campaigns"))
	adminMerchant.Post("/:store_id/campaigns", proxy.MerchantProxy(merchantAddr, "/admin/merchants/:store_id/campaigns"))
	adminMerchant.Get("/:store_id/kyc", proxy.MerchantProxy(merchantAddr, "/admin/merchants/:store_id/kyc"))
	adminMerchant.Put("/:store_id/kyc/verify", proxy.MerchantProxy(merchantAddr, "/admin/merchants/:store_id/kyc/verify"))
	adminMerchant.Put("/:store_id/subscription", proxy.MerchantProxy(merchantAddr, "/admin/merchants/:store_id/subscription"))
}

func registerAdminCampaignRoutes(admin fiber.Router, merchantAddr string) {
	adminCampaign := admin.Group("/campaigns")
	adminCampaign.Get("", proxy.MerchantProxy(merchantAddr, "/admin/campaigns"))
	adminCampaign.Get("/", proxy.MerchantProxy(merchantAddr, "/admin/campaigns"))
	adminCampaign.Put("/:id/status", proxy.MerchantProxy(merchantAddr, "/admin/campaigns/:id/status"))
	adminCampaign.Delete("/:id", proxy.MerchantProxy(merchantAddr, "/admin/campaigns/:id"))
}

func registerAdminKYCRoutes(admin fiber.Router, merchantAddr string) {
	adminKYC := admin.Group("/kyc")
	adminKYC.Get("", proxy.MerchantProxy(merchantAddr, "/admin/kyc"))
	adminKYC.Get("/", proxy.MerchantProxy(merchantAddr, "/admin/kyc"))
	adminKYC.Get("/:store_id", proxy.MerchantProxy(merchantAddr, "/admin/kyc/:store_id"))
	adminKYC.Put("/:store_id/verify", proxy.MerchantProxy(merchantAddr, "/admin/kyc/:store_id/verify"))
}

func registerAdminRewardsRoutes(admin fiber.Router, merchantAddr string) {
	adminRewards := admin.Group("/rewards")
	adminRewards.Get("/history", proxy.MerchantProxy(merchantAddr, "/admin/rewards/history"))
	adminRewards.Get("/summary", proxy.MerchantProxy(merchantAddr, "/admin/rewards/summary"))
}

func registerAdminSubscriptionRoutes(admin fiber.Router, merchantAddr string) {
	adminSubscriptions := admin.Group("/subscriptions")
	adminSubscriptions.Get("", proxy.MerchantProxy(merchantAddr, "/admin/subscriptions"))
	adminSubscriptions.Get("/", proxy.MerchantProxy(merchantAddr, "/admin/subscriptions"))
}

func registerAdminSupportRoutes(admin fiber.Router, registry *proxy.ServiceRegistry) {
	adminSupport := admin.Group("/support")
	adminSupport.Get("/feedbacks", proxy.HTTPProxy(registry, "support", "/admin/support/feedbacks"))
	adminSupport.Get("/complaints", proxy.HTTPProxy(registry, "support", "/admin/support/complaints"))
	adminSupport.Get("/reports", proxy.HTTPProxy(registry, "support", "/admin/support/reports"))
	
	// Chat
	adminSupport.Get("/chat/history/:user_id", proxy.HTTPProxy(registry, "support", "/admin/support/chat/history/:user_id"))
	adminSupport.Get("/chat/inbox", proxy.HTTPProxy(registry, "support", "/admin/support/chat/inbox"))
	adminSupport.Post("/chat/send", proxy.HTTPProxy(registry, "support", "/admin/support/chat/send"))
}