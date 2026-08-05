package http

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/usecases"
)

type EarningsHandler struct {
	usecase usecases.EarningsUseCase
}

func NewEarningsHandler(uc usecases.EarningsUseCase) *EarningsHandler {
	return &EarningsHandler{usecase: uc}
}

func (h *EarningsHandler) GetEarnings(c fiber.Ctx) error {
	storeIDStr := resolveMerchantID(c)
	if storeIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "store_id is required"})
	}
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid store_id"})
	}

	stats, requestedMap, err := h.usecase.GetEarningsData(c.Context(), storeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch earnings data"})
	}

	return c.JSON(EarningsResponse{
		TotalScans:          stats.TotalScans,
		TotalEarned:         stats.TotalEarned,
		FormattedTotal:      fmt.Sprintf("₹%.2f", stats.TotalEarned),
		RequestedMilestones: requestedMap,
	})
}

func (h *EarningsHandler) RequestPayout(c fiber.Ctx) error {
	storeIDStr := resolveMerchantID(c)
	if storeIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "store_id is required"})
	}
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid store_id"})
	}

	var req RequestPayoutRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.MilestoneScans <= 0 || req.RewardAmount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid milestone parameters"})
	}

	_, err = h.usecase.RequestPayout(c.Context(), storeID, req.MilestoneScans, req.RewardAmount)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Payout requested successfully"})
}
