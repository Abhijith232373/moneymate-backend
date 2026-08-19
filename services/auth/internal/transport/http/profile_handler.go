package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	usecase "github.com/moneymate-2026/moneymate-backend/auth/internal/usecases"
	response "github.com/moneymate-2026/moneymate-backend/shared/pkg/responses"
)

type ProfilePictureHandler struct {
	usecase usecase.ProfilePictureUsecase
}

func NewProfilePictureHandler(usecase usecase.ProfilePictureUsecase) *ProfilePictureHandler {
	return &ProfilePictureHandler{usecase: usecase}
}

func (h *ProfilePictureHandler) Presign(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "authentication required")
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return response.Unauthorized(c, "invalid user session")
	}

	var req struct {
		ContentType string `json:"content_type" validate:"required"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return response.BadRequest(c, formatValidationErrors(err), "validation failed")
	}

	result, err := h.usecase.PresignUpload(c.Context(), uid, req.ContentType)
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "upload url generated", result)
}

func (h *ProfilePictureHandler) Set(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return response.Unauthorized(c, "authentication required")
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return response.Unauthorized(c, "invalid user session")
	}

	var req struct {
		URL string `json:"url" validate:"required,url"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return response.BadRequest(c, formatValidationErrors(err), "validation failed")
	}

	user, err := h.usecase.SetProfilePicture(c.Context(), uid, req.URL)
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "profile picture updated", user)
}