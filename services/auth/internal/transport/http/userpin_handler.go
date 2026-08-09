package http

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/auth/internal/infra/tokenissuer"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/usecases"
	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"
)

type UserPinHandler struct {
	usecase usecase.UserPinUsecase
	issuer  *tokenissuer.Issuer
}

func NewUserPinHandler(u usecase.UserPinUsecase, issuer *tokenissuer.Issuer) *UserPinHandler {
	return &UserPinHandler{usecase: u, issuer: issuer}
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

	token, expiresAt, err := h.issuer.IssueTransactionToken(userID)
	if err != nil {
		appErr := apperrors.NewAppError(fiber.StatusInternalServerError, "INTERNAL_ERROR", "Failed to issue transaction token", err)
		return c.Status(appErr.StatusCode).JSON(appErr)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":           "PIN verified successfully",
		"transaction_token": token,
		"expires_in":        int(time.Until(expiresAt).Seconds()),
	})
}

