package http

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/gofiber/fiber/v3"

	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/usecases"
	response "github.com/moneymate-2026/moneymate-backend/shared/pkg/responses"
)

type InboxHandler struct {
	inbox *usecases.InboxUsecase
}

func NewInboxHandler(inbox *usecases.InboxUsecase) *InboxHandler {
	return &InboxHandler{inbox: inbox}
}

// List returns the caller's saved notifications, newest first.
func (h *InboxHandler) List(c fiber.Ctx) error {
	recipientID, err := uuid.Parse(recipientIDFromLocals(c))
	if err != nil {
		return response.Unauthorized(c, "authentication required")
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	msgs, err := h.inbox.List(c.Context(), recipientTypeFromLocals(c), recipientID, page, pageSize)
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "inbox fetched", msgs)
}

// Get returns one message (ownership enforced by the repo query).
func (h *InboxHandler) Get(c fiber.Ctx) error {
	recipientID, err := uuid.Parse(recipientIDFromLocals(c))
	if err != nil {
		return response.Unauthorized(c, "authentication required")
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid notification id")
	}
	msg, err := h.inbox.Get(c.Context(), id, recipientTypeFromLocals(c), recipientID)
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "notification fetched", msg)
}

// MarkRead flips the unread flag. Returns ok even when nothing matched —
// marking an unknown id is not an error for the client.
func (h *InboxHandler) MarkRead(c fiber.Ctx) error {
	recipientID, err := uuid.Parse(recipientIDFromLocals(c))
	if err != nil {
		return response.Unauthorized(c, "authentication required")
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid notification id")
	}
	if err := h.inbox.MarkRead(c.Context(), id, recipientTypeFromLocals(c), recipientID); err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "marked as read", nil)
}
