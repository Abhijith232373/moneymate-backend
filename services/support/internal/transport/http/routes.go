package http

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(router fiber.Router, h *SupportHandler, chatHandler *ChatHandler) {
	support := router.Group("/support")

	support.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"service": "support",
		})
	})



	// Feedback
	support.Post("/feedbacks", h.CreateFeedback)
	support.Get("/feedbacks", h.ListFeedbacks)

	// Complaints
	support.Post("/complaints", h.CreateComplaint)
	support.Get("/complaints", h.ListComplaints)
	support.Get("/complaints/user/:user_id", h.ListComplaintsByUser)

	// Reports
	support.Post("/reports", h.CreateReport)
	support.Get("/reports", h.ListReports)
	support.Get("/reports/user/:reporter_id", h.ListReportsByUser)

	// Chat
	support.Post("/chat/send", chatHandler.SendMessage)
	support.Get("/chat/history/:admin_id", chatHandler.GetChatHistory)
	support.Put("/chat/read/:sender_id", chatHandler.MarkMessagesAsRead)
}

func RegisterAdminRoutes(router fiber.Router, h *SupportHandler, chatHandler *ChatHandler) {
	admin := router.Group("/admin/support")



	// Admin viewing and updating statuses
	admin.Get("/feedbacks", h.ListFeedbacks)
	
	admin.Get("/complaints", h.ListComplaints)

	admin.Get("/reports", h.ListReports)

	// Chat
	admin.Get("/chat/history", chatHandler.GetAdminChatHistory)
	admin.Post("/chat/send", chatHandler.SendMessage)
}
