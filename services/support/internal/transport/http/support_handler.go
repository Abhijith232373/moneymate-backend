package http

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/abijith/moneymate-backend/services/support/internal/usecase"
	"github.com/moneymate-2026/moneymate-backend/shared/pkg/responses"
)

type SupportHandler struct {
	uc usecase.SupportUseCase
}

func NewSupportHandler(uc usecase.SupportUseCase) *SupportHandler {
	return &SupportHandler{uc: uc}
}

type createFeedbackReq struct {
	UserID      string `json:"user_id"`
	UserType    string `json:"user_type"`
	Rating      int    `json:"rating"`
	Description string `json:"description"`
}

func (h *SupportHandler) CreateFeedback(c fiber.Ctx) error {
	var req createFeedbackReq
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request")
	}

	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		return response.BadRequest(c, nil, "invalid user_id")
	}

	fb, err := h.uc.CreateFeedback(c.Context(), uid, req.UserType, req.Rating, req.Description)
	if err != nil {
		return response.InternalServerError(c)
	}

	return response.Created(c, "Feedback created successfully", fb)
}

func (h *SupportHandler) ListFeedbacks(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	fbs, err := h.uc.ListFeedbacks(c.Context(), int32(limit), int32(offset))
	if err != nil {
		return response.InternalServerError(c)
	}

	return response.OK(c, "Feedbacks fetched successfully", fbs)
}

type createComplaintReq struct {
	UserID      string `json:"user_id"`
	UserType    string `json:"user_type"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (h *SupportHandler) CreateComplaint(c fiber.Ctx) error {
	var req createComplaintReq
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request")
	}

	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		return response.BadRequest(c, nil, "invalid user_id")
	}

	comp, err := h.uc.CreateComplaint(c.Context(), uid, req.UserType, req.Title, req.Description)
	if err != nil {
		return response.InternalServerError(c)
	}

	return response.Created(c, "Complaint created successfully", comp)
}

func (h *SupportHandler) ListComplaints(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	comps, err := h.uc.ListComplaints(c.Context(), int32(limit), int32(offset))
	if err != nil {
		return response.InternalServerError(c)
	}

	return response.OK(c, "Complaints fetched successfully", comps)
}

func (h *SupportHandler) ListComplaintsByUser(c fiber.Ctx) error {
	userID := c.Params("user_id")
	userType := c.Query("user_type")

	uid, err := uuid.Parse(userID)
	if err != nil {
		return response.BadRequest(c, nil, "invalid user_id")
	}

	comps, err := h.uc.ListComplaintsByUser(c.Context(), uid, userType)
	if err != nil {
		return response.InternalServerError(c)
	}

	return response.OK(c, "Complaints fetched successfully", comps)
}



type createReportReq struct {
	ReporterID   string `json:"reporter_id"`
	ReporterType string `json:"reporter_type"`
	ReportedVPA  string `json:"reported_vpa"`
	Title        string `json:"title"`
	Description  string `json:"description"`
}

func (h *SupportHandler) CreateReport(c fiber.Ctx) error {
	var req createReportReq
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request")
	}

	uid, err := uuid.Parse(req.ReporterID)
	if err != nil {
		return response.BadRequest(c, nil, "invalid reporter_id")
	}

	rep, err := h.uc.CreateReport(c.Context(), uid, req.ReporterType, req.ReportedVPA, req.Title, req.Description)
	if err != nil {
		return response.InternalServerError(c)
	}

	return response.Created(c, "Report created successfully", rep)
}

func (h *SupportHandler) ListReports(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	reps, err := h.uc.ListReports(c.Context(), int32(limit), int32(offset))
	if err != nil {
		return response.InternalServerError(c)
	}

	return response.OK(c, "Reports fetched successfully", reps)
}

func (h *SupportHandler) ListReportsByUser(c fiber.Ctx) error {
	reporterID := c.Params("reporter_id")
	reporterType := c.Query("reporter_type")

	uid, err := uuid.Parse(reporterID)
	if err != nil {
		return response.BadRequest(c, nil, "invalid reporter_id")
	}

	reps, err := h.uc.ListReportsByUser(c.Context(), uid, reporterType)
	if err != nil {
		return response.InternalServerError(c)
	}

	return response.OK(c, "Reports fetched successfully", reps)
}

