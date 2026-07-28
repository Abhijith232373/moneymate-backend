package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/usecases"
)

type DashboardHandler struct {
	usecase *usecases.DashboardUseCase
}

func NewDashboardHandler(uc *usecases.DashboardUseCase) *DashboardHandler {
	return &DashboardHandler{usecase: uc}
}

// GetDashboard serves the aggregated statistics, transactions, and promotional campaigns for the merchant dashboard.
func (h *DashboardHandler) GetDashboard(c fiber.Ctx) error {
	id := resolveMerchantID(c)

	out, err := h.usecase.GetDashboard(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}

	var stats []StatCardResponse
	for _, s := range out.Stats {
		stats = append(stats, StatCardResponse{
			Title:          s.Title,
			Value:          s.Value,
			Icon:           s.Icon,
			IconColorClass: s.IconColorClass,
			BorderClass:    s.BorderClass,
			Trend: TrendResponse{
				Text: s.TrendText,
				Type: s.TrendType,
			},
		})
	}

	var txs []DashboardTransactionResponse
	for _, tx := range out.Transactions {
		txs = append(txs, DashboardTransactionResponse{
			Time:     tx.Time,
			Customer: tx.Customer,
			Initial:  tx.Initial,
			Color:    tx.Color,
			Amount:   tx.Amount,
			Reward:   tx.Reward,
			Status:   tx.Status,
		})
	}

	var camps []DashboardCampaignResponse
	for _, camp := range out.Campaigns {
		camps = append(camps, DashboardCampaignResponse{
			ID:     camp.ID,
			Name:   camp.Name,
			Type:   camp.Type,
			Status: camp.Status,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": DashboardResponse{
			Stats:        stats,
			Transactions: txs,
			Campaigns:    camps,
			Balance:      out.Balance,
			MerchantID:   out.MerchantID,
			BusinessName: out.BusinessName,
		},
	})
}
