package middlewares

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

func getSubRouteKey(moduleKey, path string) string {
	switch moduleKey {
	case "auth_routes":
		if strings.Contains(path, "/merchant/auth/register") {
			return "auth_merchant_register"
		}
		if strings.Contains(path, "/merchant/auth/login") {
			return "auth_merchant_login"
		}
		if strings.Contains(path, "/auth/admin/login") {
			return "auth_admin_login"
		}
		if strings.Contains(path, "/auth/user/register") || strings.Contains(path, "/auth/register") {
			return "auth_user_register"
		}
		if strings.Contains(path, "/auth/login") {
			return "auth_user_login"
		}
	case "pin_routes":
		if strings.Contains(path, "/pin/verify") {
			return "pin_verification"
		}
		if strings.Contains(path, "/pin") {
			return "pin_management"
		}
	case "admin_routes":
		if strings.Contains(path, "/admin/merchants") || strings.Contains(path, "/admin/kyc") {
			return "admin_merchants_kyc"
		}
		if strings.Contains(path, "/admin/campaigns") {
			return "admin_master_campaigns"
		}
		if strings.Contains(path, "/admin/config") {
			return "admin_platform_config"
		}
		if strings.Contains(path, "/admin/audit") {
			return "admin_system_audit"
		}
	case "merchant_routes":
		if strings.Contains(path, "/merchant/campaigns") {
			return "merch_campaigns_management"
		}
		if strings.Contains(path, "/merchant/subscriptions") {
			return "merch_subscriptions"
		}
		if strings.Contains(path, "/merchant/wallet") || strings.Contains(path, "/merchant/earnings") || strings.Contains(path, "/merchant/rewards") {
			return "merch_wallet_payouts"
		}
		if strings.Contains(path, "/merchant/dashboard") {
			return "merch_dashboard_analytics"
		}
	case "payment_routes":
		if strings.Contains(path, "/transfers") {
			return "pay_p2p_transfers"
		}
		if strings.Contains(path, "/deposits") {
			return "pay_fiat_deposits"
		}
		if strings.Contains(path, "/withdrawals") {
			return "pay_fiat_withdrawals"
		}
		if strings.Contains(path, "/wallets") || strings.Contains(path, "/transactions") {
			return "pay_wallet_balances"
		}
	case "secure_routes":
		if strings.Contains(path, "/sync") {
			return "secure_identity_sync"
		}
		return "secure_service_comm"
	case "support_routes":
		if strings.Contains(path, "/complaints") {
			return "support_user_complaints"
		}
		if strings.Contains(path, "/reports") {
			return "support_fraud_reports"
		}
		if strings.Contains(path, "/chat") {
			return "support_live_chat"
		}
	case "downstream_routes":
		if strings.Contains(path, "/payment") {
			return "ds_payment_aggregator"
		}
		if strings.Contains(path, "/campaign") {
			return "ds_campaign_service"
		}
		if strings.Contains(path, "/notification") {
			return "ds_notification_engine"
		}
	}
	return ""
}

func MaintenanceMode(rdb *redis.Client, moduleKey string) fiber.Handler {
	return func(c fiber.Ctx) error {
		// CRITICAL EXEMPTION: Never block the configuration endpoints. 
		// If an admin toggles off the 'admin_routes' master switch, they would lock 
		// themselves out of the system and never be able to turn it back on.
		if strings.Contains(c.Path(), "/admin/config") || strings.Contains(c.Path(), "/auth/admin/login") {
			return c.Next()
		}

		ctx := context.Background()
		
		// 1. Check Master Module toggle first
		val, err := rdb.Get(ctx, "config:module:"+moduleKey).Result()
		if err == nil && val == "false" {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"success": false,
				"error":   "Maintenance Mode",
				"message": "This module is temporarily unavailable due to system maintenance.",
			})
		}

		// 2. Check Specific Sub-route toggle
		subKey := getSubRouteKey(moduleKey, c.Path())
		if subKey != "" {
			subVal, subErr := rdb.Get(ctx, "config:module:"+subKey).Result()
			if subErr == nil && subVal == "false" {
				return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
					"success": false,
					"error":   "Maintenance Mode",
					"message": "This specific service is temporarily unavailable due to maintenance.",
				})
			}
		}

		return c.Next()
	}
}
