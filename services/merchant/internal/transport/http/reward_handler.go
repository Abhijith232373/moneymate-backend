package http

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/usecases"
)

// RewardHandler exposes REST/HTTP endpoints for the Rewards Center dashboard, history, and bank payouts.
type RewardHandler struct {
	// usecase orchestrates all reward validations, balances, and withdrawal transactions.
	usecase usecases.RewardUseCase
}

// NewRewardHandler initializes and returns a new RewardHandler instance with injected dependencies.
func NewRewardHandler(uc usecases.RewardUseCase) *RewardHandler {
	return &RewardHandler{usecase: uc}
}

// GetRewardSummary handles GET /merchant/:store_id/rewards/summary and returns the current balance and scan counters.
func (h *RewardHandler) GetRewardSummary(c fiber.Ctx) error {
	storeIDStr := c.Params("store_id")
	if storeIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "store_id parameter is required"})
	}
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid store_id UUID format"})
	}

	summary, err := h.usecase.GetRewardSummary(c.Context(), storeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	formattedBalance := fmt.Sprintf("$%.2f", summary.AvailableBalance)
	formattedGrowth := fmt.Sprintf("+%.0f%% this week", summary.WeeklyGrowthPercentage)
	if summary.WeeklyGrowthPercentage < 0 {
		formattedGrowth = fmt.Sprintf("%.0f%% this week", summary.WeeklyGrowthPercentage)
	}

	return c.JSON(RewardSummaryResponse{
		StoreID:                summary.StoreID.String(),
		AvailableBalance:       summary.AvailableBalance,
		TotalScans:             summary.TotalScans,
		PremiumPoints:          summary.PremiumPoints,
		WeeklyGrowthPercentage: summary.WeeklyGrowthPercentage,
		FormattedBalance:       formattedBalance,
		FormattedGrowth:        formattedGrowth,
	})
}

// GetRewardHistory handles GET /merchant/:store_id/rewards/history with query parameters for filtering and search.
func (h *RewardHandler) GetRewardHistory(c fiber.Ctx) error {
	storeIDStr := c.Params("store_id")
	if storeIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "store_id parameter is required"})
	}
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid store_id UUID format"})
	}

	filter := c.Query("filter", "all")
	searchQuery := c.Query("search", "")
	limitStr := c.Query("limit", "50")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	transactions, err := h.usecase.GetRewardHistory(c.Context(), storeID, filter, searchQuery, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var response []RewardTransactionResponse
	for _, tx := range transactions {
		formattedAmt := fmt.Sprintf("+$%.2f", tx.Amount)
		if tx.Amount < 0 {
			formattedAmt = fmt.Sprintf("-$%.2f", -tx.Amount)
		}

		response = append(response, RewardTransactionResponse{
			ID:              tx.ID.String(),
			StoreID:         tx.StoreID.String(),
			CampaignName:    tx.CampaignName,
			DisplayID:       tx.DisplayID,
			Status:          tx.Status,
			Amount:          tx.Amount,
			FormattedAmount: formattedAmt,
			TransactionType: tx.TransactionType,
			CreatedAt:       tx.CreatedAt.Format(time.RFC3339),
			FormattedDate:   formatFriendlyDate(tx.CreatedAt),
		})
	}

	if response == nil {
		response = []RewardTransactionResponse{}
	}

	return c.JSON(response)
}

// RedeemRewards handles POST /merchant/:store_id/rewards/redeem to initiate a bank transfer payout.
func (h *RewardHandler) RedeemRewards(c fiber.Ctx) error {
	storeIDStr := c.Params("store_id")
	if storeIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "store_id parameter is required"})
	}
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid store_id UUID format"})
	}

	var req RedeemBalanceRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON request payload"})
	}

	redemption, updatedSummary, err := h.usecase.RedeemRewards(c.Context(), storeID, req.Amount, req.ConfirmBankTransferAuthorization)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	refID := ""
	if redemption.ReferenceID != nil {
		refID = *redemption.ReferenceID
	}

	remBal := 0.00
	if updatedSummary != nil {
		remBal = updatedSummary.AvailableBalance
	}

	return c.Status(fiber.StatusOK).JSON(RedeemBalanceResponse{
		RedemptionID:              redemption.ID.String(),
		StoreID:                   redemption.StoreID.String(),
		AmountRedeemed:            redemption.Amount,
		FormattedAmountRedeemed:   fmt.Sprintf("$%.2f", redemption.Amount),
		RemainingBalance:          remBal,
		FormattedRemainingBalance: fmt.Sprintf("$%.2f", remBal),
		Status:                    redemption.Status,
		ReferenceID:               refID,
		Message:                   "Bank transfer authorization confirmed. Your redemption is processing.",
	})
}

// formatFriendlyDate converts a UTC timestamp into UI presentation strings like "Today, 14:32 PM" or "Yesterday, 18:45 PM".
func formatFriendlyDate(t time.Time) string {
	now := time.Now()
	y1, m1, d1 := now.Date()
	y2, m2, d2 := t.Date()

	timeStr := t.Format("15:04 PM")
	if y1 == y2 && m1 == m2 && d1 == d2 {
		return fmt.Sprintf("Today, %s", timeStr)
	}
	yesterday := now.AddDate(0, 0, -1)
	y3, m3, d3 := yesterday.Date()
	if y3 == y2 && m3 == m2 && d3 == d2 {
		return fmt.Sprintf("Yesterday, %s", timeStr)
	}
	return fmt.Sprintf("%s, %s", t.Format("Jan 02"), timeStr)
}
