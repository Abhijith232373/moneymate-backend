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

	// Complaints
	support.Post("/complaints", h.CreateComplaint)
	support.Get("/complaints/me", h.ListComplaintsByUser)

	// Reports
	support.Post("/reports", h.CreateReport)
	support.Get("/reports/me", h.ListReportsByUser)

	// Chat
	support.Post("/chat/send", chatHandler.SendMessage)
	support.Put("/chat/read/:sender_id", chatHandler.MarkMessagesAsRead)
}

func RegisterAdminRoutes(router fiber.Router, h *SupportHandler, chatHandler *ChatHandler) {
	admin := router.Group("/admin/support")



	// Admin viewing and updating statuses
	admin.Get("/feedbacks", h.ListFeedbacks)
	
	admin.Get("/complaints", h.ListComplaints)

	admin.Get("/reports", h.ListReports)

	admin.Get("/audit-logs", h.ListAuditLogs)
	admin.Post("/audit-logs", h.CreateAuditLog)

	// Chat
	admin.Get("/chat/history/:user_id", chatHandler.GetChatHistoryForAdmin)
	admin.Get("/chat/inbox", chatHandler.GetAdminChatHistory) // Inbox view for all chats
	admin.Post("/chat/send", chatHandler.SendMessage)
}
