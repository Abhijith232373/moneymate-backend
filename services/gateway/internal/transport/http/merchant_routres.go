package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/moneymate-2026/moneymate-backend/gateway/internal/proxy"
)

func registerMerchantRoutes(api fiber.Router, merchantAddr string) {
	merchantUnauth := api.Group("/merchant")
	merchantUnauth.Post("/register", proxy.MerchantProxy(merchantAddr, "/merchant/register"))
	merchantUnauth.Post("/login", proxy.MerchantProxy(merchantAddr, "/merchant/login"))
	merchantUnauth.Get("/public/campaigns", proxy.MerchantProxy(merchantAddr, "/merchant/public/campaigns"))
	merchantUnauth.Get("/public/subscriptions/plans", proxy.MerchantProxy(merchantAddr, "/merchant/public/subscriptions/plans"))

	merchant := api.Group("/merchant")
	merchant.Get("/health", proxy.MerchantProxy(merchantAddr, "/merchant/health"))
	merchant.Get("/dashboard", proxy.MerchantProxy(merchantAddr, "/merchant/dashboard"))
	merchant.Get("/:store_id/dashboard", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/dashboard"))
	merchant.Get("/status/:owner_id", proxy.MerchantProxy(merchantAddr, "/merchant/status/:owner_id"))
	merchant.Get("/pending", proxy.MerchantProxy(merchantAddr, "/merchant/pending"))

	merchant.Get("/profile", proxy.MerchantProxy(merchantAddr, "/merchant/profile"))
	merchant.Get("/:store_id/profile", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/profile"))
	merchant.Put("/profile", proxy.MerchantProxy(merchantAddr, "/merchant/profile"))
	merchant.Put("/:store_id/profile", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/profile"))
	merchant.Post("/profile", proxy.MerchantProxy(merchantAddr, "/merchant/profile"))
	merchant.Post("/:store_id/profile", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/profile"))

	merchant.Post("/campaigns", proxy.MerchantProxy(merchantAddr, "/merchant/campaigns"))
	merchant.Post("/:store_id/campaigns", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/campaigns"))
	merchant.Get("/campaigns", proxy.MerchantProxy(merchantAddr, "/merchant/campaigns"))
	merchant.Get("/:store_id/campaigns", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/campaigns"))
	merchant.Put("/campaigns/:id/status", proxy.MerchantProxy(merchantAddr, "/merchant/campaigns/:id/status"))
	merchant.Put("/:store_id/campaigns/:id/status", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/campaigns/:id/status"))

	merchant.Get("/rewards/summary", proxy.MerchantProxy(merchantAddr, "/merchant/rewards/summary"))
	merchant.Get("/:store_id/rewards/summary", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/rewards/summary"))
	merchant.Get("/rewards/history", proxy.MerchantProxy(merchantAddr, "/merchant/rewards/history"))
	merchant.Get("/:store_id/rewards/history", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/rewards/history"))
	merchant.Post("/rewards/redeem", proxy.MerchantProxy(merchantAddr, "/merchant/rewards/redeem"))
	merchant.Post("/:store_id/rewards/redeem", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/rewards/redeem"))

	merchant.Get("/subscriptions/plans", proxy.MerchantProxy(merchantAddr, "/merchant/subscriptions/plans"))
	merchant.Get("/:store_id/subscriptions/plans", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/subscriptions/plans"))
	merchant.Get("/subscriptions/current", proxy.MerchantProxy(merchantAddr, "/merchant/subscriptions/current"))
	merchant.Get("/:store_id/subscriptions/current", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/subscriptions/current"))
	merchant.Post("/subscriptions/change", proxy.MerchantProxy(merchantAddr, "/merchant/subscriptions/change"))
	merchant.Post("/:store_id/subscriptions/change", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/subscriptions/change"))
	merchant.Post("/subscriptions/upgrade/initiate", proxy.MerchantProxy(merchantAddr, "/merchant/subscriptions/upgrade/initiate"))
	merchant.Post("/:store_id/subscriptions/upgrade/initiate", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/subscriptions/upgrade/initiate"))
	merchant.Post("/subscriptions/upgrade/verify", proxy.MerchantProxy(merchantAddr, "/merchant/subscriptions/upgrade/verify"))
	merchant.Post("/:store_id/subscriptions/upgrade/verify", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/subscriptions/upgrade/verify"))

	merchant.Get("/wallet", proxy.MerchantProxy(merchantAddr, "/merchant/wallet"))
	merchant.Get("/:store_id/wallet", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/wallet"))
	merchant.Get("/earnings", proxy.MerchantProxy(merchantAddr, "/merchant/earnings"))
	merchant.Get("/:store_id/earnings", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/earnings"))
	merchant.Post("/earnings/payouts", proxy.MerchantProxy(merchantAddr, "/merchant/earnings/payouts"))
	merchant.Post("/:store_id/earnings/payouts", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/earnings/payouts"))

	merchant.Get("/kyc", proxy.MerchantProxy(merchantAddr, "/merchant/kyc"))
	merchant.Get("/:store_id/kyc", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/kyc"))
	merchant.Get("/kyc/status", proxy.MerchantProxy(merchantAddr, "/merchant/kyc/status"))
	merchant.Get("/:store_id/kyc/status", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/kyc/status"))
	merchant.Put("/kyc", proxy.MerchantProxy(merchantAddr, "/merchant/kyc"))
	merchant.Put("/:store_id/kyc", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/kyc"))
	merchant.Post("/kyc/update", proxy.MerchantProxy(merchantAddr, "/merchant/kyc/update"))
	merchant.Post("/:store_id/kyc/update", proxy.MerchantProxy(merchantAddr, "/merchant/:store_id/kyc/update"))
}