package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/moneymate-2026/moneymate-backend/gateway/internal/proxy"
)

func registerSupportRoutes(api fiber.Router, authMiddleware fiber.Handler, registry *proxy.ServiceRegistry) {
	supportUnauth := api.Group("/support")
	supportUnauth.Get("/health", proxy.HTTPProxy(registry, "support", "/support/health"))

	support := api.Group("/support")
	support.Use(authMiddleware)

	// Feedback
	support.Post("/feedbacks", proxy.HTTPProxy(registry, "support", "/support/feedbacks"))

	// Complaints
	support.Post("/complaints", proxy.HTTPProxy(registry, "support", "/support/complaints"))
	support.Get("/complaints/me", proxy.HTTPProxy(registry, "support", "/support/complaints/me"))

	// Reports
	support.Post("/reports", proxy.HTTPProxy(registry, "support", "/support/reports"))
	support.Get("/reports/me", proxy.HTTPProxy(registry, "support", "/support/reports/me"))

	// Chat
	support.Post("/chat/send", proxy.HTTPProxy(registry, "support", "/support/chat/send"))
	support.Put("/chat/read/:sender_id", proxy.HTTPProxy(registry, "support", "/support/chat/read/:sender_id"))
}
