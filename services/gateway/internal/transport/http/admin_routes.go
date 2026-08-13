package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/moneymate-2026/moneymate-backend/gateway/internal/middlewares"
	"github.com/moneymate-2026/moneymate-backend/gateway/internal/proxy"
)

func registerAdminRoutes(api fiber.Router, authMiddleware fiber.Handler, authAddr, merchantAddr string) {
	admin := api.Group("/admin")
	admin.Use(authMiddleware)
	admin.Use(middlewares.RequireRole("admin"))

	admin.Get("/dashboard", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "admin dashboard data",
			"data": fiber.Map{
				"total_users": 0, "total_merchants": 0,
				"pending_reviews": 0, "total_transactions": 0,
			},
		})
	})
	admin.Post("/refresh", proxy.AuthProxy(authAddr, "/auth/refresh"))

	registerAdminUserRoutes(admin, authAddr)
	registerAdminRoleRoutes(admin, authAddr)
	registerAdminPermissionRoutes(admin, authAddr)
	registerAdminMerchantRoutes(admin, merchantAddr)
	registerAdminCampaignRoutes(admin, merchantAddr)
	registerAdminKYCRoutes(admin, merchantAddr)
	registerAdminRewardsRoutes(admin, merchantAddr)
	registerAdminSubscriptionRoutes(admin, merchantAddr)
}

func registerAdminUserRoutes(admin fiber.Router, authAddr string) {
	adminUser := admin.Group("/users")
	adminUser.Post("/", proxy.AuthProxy(authAddr, "/admin/users"))
	adminUser.Get("/", proxy.AuthProxy(authAddr, "/admin/users"))
	adminUser.Get("/:id", proxy.AuthProxy(authAddr, "/admin/users/:id"))
	adminUser.Put("/:id", proxy.AuthProxy(authAddr, "/admin/users/:id"))
	adminUser.Patch("/:id/status", proxy.AuthProxy(authAddr, "/admin/users/:id/status"))
	adminUser.Delete("/:id", proxy.AuthProxy(authAddr, "/admin/users/:id"))
}

func registerAdminRoleRoutes(admin fiber.Router, authAddr string) {
	adminRole := admin.Group("/roles")
	adminRole.Post("/", proxy.AuthProxy(authAddr, "/admin/roles"))
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
	adminPermission.Post("/", proxy.AuthProxy(authAddr, "/admin/permissions"))
	adminPermission.Get("/", proxy.AuthProxy(authAddr, "/admin/permissions"))
	adminPermission.Get("/:id", proxy.AuthProxy(authAddr, "/admin/permissions/:id"))
	adminPermission.Delete("/:id", proxy.AuthProxy(authAddr, "/admin/permissions/:id"))
	adminPermission.Post("/assign", proxy.AuthProxy(authAddr, "/admin/permissions/assign"))
	adminPermission.Delete("/roles/:roleId/permissions/:permissionId", proxy.AuthProxy(authAddr, "/admin/permissions/roles/:roleId/permissions/:permissionId"))
	adminPermission.Get("/roles/:roleId", proxy.AuthProxy(authAddr, "/admin/permissions/roles/:roleId"))
}

func registerAdminMerchantRoutes(admin fiber.Router, merchantAddr string) {
	adminMerchant := admin.Group("/merchants")
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
	adminCampaign.Get("/", proxy.MerchantProxy(merchantAddr, "/admin/campaigns"))
	adminCampaign.Put("/:id/status", proxy.MerchantProxy(merchantAddr, "/admin/campaigns/:id/status"))
	adminCampaign.Delete("/:id", proxy.MerchantProxy(merchantAddr, "/admin/campaigns/:id"))
}

func registerAdminKYCRoutes(admin fiber.Router, merchantAddr string) {
	adminKYC := admin.Group("/kyc")
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
	adminSubscriptions.Get("/", proxy.MerchantProxy(merchantAddr, "/admin/subscriptions"))
}