package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/auth/internal/usecases"
	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"
)

type UserPinHandler struct {
	usecase usecase.UserPinUsecase
}

func NewUserPinHandler(u usecase.UserPinUsecase) *UserPinHandler {
	return &UserPinHandler{usecase: u}
}

func (h *UserPinHandler) SetPIN(c fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		appErr := apperrors.NewAppError(fiber.StatusUnauthorized, "UNAUTHORIZED", "Invalid user ID", err)
		return c.Status(appErr.StatusCode).JSON(appErr)
	}

	var req setPINRequest
	if err := c.Bind().Body(&req); err != nil {
		appErr := apperrors.NewAppError(fiber.StatusBadRequest, "INVALID_INPUT", "Invalid request body", err)
		return c.Status(appErr.StatusCode).JSON(appErr)
	}

	ucReq := usecase.SetPINRequest{PIN: req.PIN}
	if err := h.usecase.SetPIN(c.Context(), userID, ucReq); err != nil {
		appErr := apperrors.ParseError(err)
		return c.Status(appErr.StatusCode).JSON(appErr)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "PIN set successfully",
	})
}

func (h *UserPinHandler) UpdatePIN(c fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		appErr := apperrors.NewAppError(fiber.StatusUnauthorized, "UNAUTHORIZED", "Invalid user ID", err)
		return c.Status(appErr.StatusCode).JSON(appErr)
	}

	var req updatePINRequest
	if err := c.Bind().Body(&req); err != nil {
		appErr := apperrors.NewAppError(fiber.StatusBadRequest, "INVALID_INPUT", "Invalid request body", err)
		return c.Status(appErr.StatusCode).JSON(appErr)
	}

	ucReq := usecase.UpdatePINRequest{
		OldPIN: req.OldPIN,
		NewPIN: req.NewPIN,
	}
	if err := h.usecase.UpdatePIN(c.Context(), userID, ucReq); err != nil {
		appErr := apperrors.ParseError(err)
		return c.Status(appErr.StatusCode).JSON(appErr)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "PIN updated successfully",
	})
}

func (h *UserPinHandler) VerifyPIN(c fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		appErr := apperrors.NewAppError(fiber.StatusUnauthorized, "UNAUTHORIZED", "Invalid user ID", err)
		return c.Status(appErr.StatusCode).JSON(appErr)
	}

	var req verifyPINRequest
	if err := c.Bind().Body(&req); err != nil {
		appErr := apperrors.NewAppError(fiber.StatusBadRequest, "INVALID_INPUT", "Invalid request body", err)
		return c.Status(appErr.StatusCode).JSON(appErr)
	}

	ucReq := usecase.VerifyPINRequest{PIN: req.PIN}
	if err := h.usecase.VerifyPIN(c.Context(), userID, ucReq); err != nil {
		appErr := apperrors.ParseError(err)
		return c.Status(appErr.StatusCode).JSON(appErr)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "PIN verified successfully",
	})
}
