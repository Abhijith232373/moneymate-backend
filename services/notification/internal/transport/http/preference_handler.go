package http

import (
	"github.com/google/uuid"
	"github.com/gofiber/fiber/v3"

	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/domain"
	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/usecases"
	response "github.com/moneymate-2026/moneymate-backend/shared/pkg/responses"
)

type PreferenceHandler struct {
	preferences *usecases.PreferenceUsecase
}

func NewPreferenceHandler(preferences *usecases.PreferenceUsecase) *PreferenceHandler {
	return &PreferenceHandler{preferences: preferences}
}

type upsertPreferenceRequest struct {
	Category domain.Category `json:"category"`
	Enabled  *bool           `json:"enabled"`
}

// List returns every category with its current on/off state.
func (h *PreferenceHandler) List(c fiber.Ctx) error {
	recipientID, err := uuid.Parse(recipientIDFromLocals(c))
	if err != nil {
		return response.Unauthorized(c, "authentication required")
	}
	prefs, err := h.preferences.List(c.Context(), recipientTypeFromLocals(c), recipientID)
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "preferences fetched", prefs)
}

// Upsert toggles one category for the caller.
func (h *PreferenceHandler) Upsert(c fiber.Ctx) error {
	var req upsertPreferenceRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request body")
	}
	if req.Enabled == nil {
		return response.BadRequest(c, nil, "enabled is required")
	}
	recipientID, err := uuid.Parse(recipientIDFromLocals(c))
	if err != nil {
		return response.Unauthorized(c, "authentication required")
	}
	pref, err := h.preferences.Upsert(c.Context(), recipientTypeFromLocals(c), recipientID, req.Category, *req.Enabled)
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "preference updated", pref)
}
