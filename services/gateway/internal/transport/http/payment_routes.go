package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/moneymate-2026/moneymate-backend/gateway/internal/proxy"
)

func registerPaymentRoutes(api fiber.Router, authMiddleware fiber.Handler, registry *proxy.ServiceRegistry) {
	payment := api.Group("/payment")
	payment.Use(authMiddleware)

	payment.Get("/wallets/me", proxy.HTTPProxy(registry, "payment", "/payment/wallets/me"))
	payment.Post("/deposits", proxy.HTTPProxy(registry, "payment", "/payment/deposits"))
	payment.Post("/deposits/confirm", proxy.HTTPProxy(registry, "payment", "/payment/deposits/confirm"))
	payment.Get("/deposits", proxy.HTTPProxy(registry, "payment", "/payment/deposits"))
	payment.Post("/withdrawals", proxy.HTTPProxy(registry, "payment", "/payment/withdrawals"))
	payment.Get("/withdrawals", proxy.HTTPProxy(registry, "payment", "/payment/withdrawals"))
}