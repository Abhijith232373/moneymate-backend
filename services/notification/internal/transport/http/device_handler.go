package http

import (
	"github.com/google/uuid"
	"github.com/gofiber/fiber/v3"

	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/usecases"
	response "github.com/moneymate-2026/moneymate-backend/shared/pkg/responses"
)

type DeviceHandler struct {
	devices *usecases.DeviceUsecase
}

func NewDeviceHandler(devices *usecases.DeviceUsecase) *DeviceHandler {
	return &DeviceHandler{devices: devices}
}

type registerDeviceRequest struct {
	DeviceID   string `json:"device_id"`
	Token      string `json:"token"`
	Platform   string `json:"platform"`
	AppVersion string `json:"app_version"`
}

// Register stores (or refreshes) the caller's push token on this device.
func (h *DeviceHandler) Register(c fiber.Ctx) error {
	var req registerDeviceRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request body")
	}

	recipientID, err := uuid.Parse(recipientIDFromLocals(c))
	if err != nil {
		return response.Unauthorized(c, "authentication required")
	}

	dev, err := h.devices.Register(c.Context(), usecases.RegisterDeviceInput{
		RecipientType: recipientTypeFromLocals(c),
		RecipientID:   recipientID,
		DeviceID:      req.DeviceID,
		Token:         req.Token,
		Platform:      req.Platform,
		AppVersion:    req.AppVersion,
	})
	if err != nil {
		return handleError(c, err)
	}
	return response.Created(c, "device registered", dev)
}

// Revoke removes one device — used on logout.
func (h *DeviceHandler) Revoke(c fiber.Ctx) error {
	recipientID, err := uuid.Parse(recipientIDFromLocals(c))
	if err != nil {
		return response.Unauthorized(c, "authentication required")
	}
	if err := h.devices.Revoke(c.Context(), recipientTypeFromLocals(c), recipientID, c.Params("device_id")); err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "device revoked", nil)
}
