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
	support.Get("/feedbacks", proxy.HTTPProxy(registry, "support", "/support/feedbacks"))

	// Complaints
	support.Post("/complaints", proxy.HTTPProxy(registry, "support", "/support/complaints"))
	support.Get("/complaints", proxy.HTTPProxy(registry, "support", "/support/complaints"))
	support.Get("/complaints/user/:user_id", proxy.HTTPProxy(registry, "support", "/support/complaints/user/:user_id"))

	// Reports
	support.Post("/reports", proxy.HTTPProxy(registry, "support", "/support/reports"))
	support.Get("/reports", proxy.HTTPProxy(registry, "support", "/support/reports"))
	support.Get("/reports/user/:reporter_id", proxy.HTTPProxy(registry, "support", "/support/reports/user/:reporter_id"))

	// Chat
	support.Post("/chat/send", proxy.HTTPProxy(registry, "support", "/support/chat/send"))
	support.Get("/chat/history/:admin_id", proxy.HTTPProxy(registry, "support", "/support/chat/history/:admin_id"))
	support.Put("/chat/read/:sender_id", proxy.HTTPProxy(registry, "support", "/support/chat/read/:sender_id"))
}
