package http

import (
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(router fiber.Router, h *SupportHandler, authMiddleware fiber.Handler) {
	support := router.Group("/support")

	support.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"service": "support",
		})
	})

	if authMiddleware != nil {
		support.Use(authMiddleware)
	}

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
}

func RegisterAdminRoutes(router fiber.Router, h *SupportHandler, authMiddleware fiber.Handler) {
	admin := router.Group("/admin/support")

	if authMiddleware != nil {
		admin.Use(authMiddleware)
	}

	// Admin viewing and updating statuses
	admin.Get("/feedbacks", h.ListFeedbacks)
	
	admin.Get("/complaints", h.ListComplaints)

	admin.Get("/reports", h.ListReports)
}
