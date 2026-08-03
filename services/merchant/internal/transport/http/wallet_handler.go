package http

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/usecases"
)

type WalletHandler struct {
	usecase usecases.WalletUseCase
}

func NewWalletHandler(uc usecases.WalletUseCase) *WalletHandler {
	return &WalletHandler{usecase: uc}
}

func (h *WalletHandler) GetWallet(c fiber.Ctx) error {
	storeIDStr := resolveMerchantID(c)
	if storeIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "store_id is required"})
	}
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid store_id"})
	}

	filterType := c.Query("filter", "all")

	walletData, err := h.usecase.GetWalletData(c.Context(), storeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch wallet data"})
	}

	txns, err := h.usecase.GetWalletTransactions(c.Context(), storeID, filterType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch wallet transactions"})
	}

	var txnsResponse []WalletTransactionResponse
	for _, txn := range txns {
		sign := "+"
		if txn.Amount < 0 {
			sign = "-"
		}
		formattedAmount := fmt.Sprintf("%s₹%.2f", sign, abs(txn.Amount))
		
		// Map UI specific text formats
		dateStr := txn.CreatedAt.Format("Jan 02, 2006")
		timeStr := txn.CreatedAt.Format("03:04 PM")
		// To match UI exact format like "Today, 10:30 AM" or "Yesterday, 04:15 PM" we can let frontend handle it or send basic strings
		// We'll send standard formats and frontend can format "Today/Yesterday".

		txnsResponse = append(txnsResponse, WalletTransactionResponse{
			ID:              txn.ID.String(),
			TransactionID:   txn.TransactionID,
			Title:           txn.Title,
			Subtitle:        txn.Subtitle,
			Amount:          txn.Amount,
			FormattedAmount: formattedAmount,
			Date:            dateStr,
			Time:            timeStr,
		})
	}

	if txnsResponse == nil {
		txnsResponse = []WalletTransactionResponse{}
	}

	overview := WalletOverviewResponse{
		AvailableBalance:       walletData.AvailableBalance,
		FormattedBalance:       fmt.Sprintf("₹%.2f", walletData.AvailableBalance),
		TotalEarnings:          walletData.TotalEarnings,
		FormattedTotalEarnings: fmt.Sprintf("₹%.2f", walletData.TotalEarnings),
		TotalRedeemed:          walletData.TotalRedeemed,
		FormattedTotalRedeemed: fmt.Sprintf("₹%.2f", walletData.TotalRedeemed),
	}

	response := WalletResponse{
		Overview:     overview,
		Transactions: txnsResponse,
	}

	return c.JSON(response)
}

func abs(val float64) float64 {
	if val < 0 {
		return -val
	}
	return val
}
