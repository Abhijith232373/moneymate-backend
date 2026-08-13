// payment/internal/transport/http/withdrawal_handler.go
package http

import (
	"github.com/gofiber/fiber/v3"

	"github.com/moneymate-2026/moneymate-backend/services/payment/internal/usecases"
	"github.com/moneymate-2026/moneymate-backend/shared/pkg/money"
	response "github.com/moneymate-2026/moneymate-backend/shared/pkg/responses"
)

type WithdrawalHandler struct {
	withdrawals usecases.WithdrawalUsecase
}

func NewWithdrawalHandler(withdrawals usecases.WithdrawalUsecase) *WithdrawalHandler {
	return &WithdrawalHandler{withdrawals: withdrawals}
}

func (h *WithdrawalHandler) Request(c fiber.Ctx) error {
	userID := userIDFromLocals(c)
	if userID == "" {
		return response.Unauthorized(c, "authentication required")
	}

	var req requestWithdrawalRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return response.BadRequest(c, nil, "validation failed")
	}

	amountPaise := int64(req.AmountRupees * 100)

	result, err := h.withdrawals.RequestWithdrawal(c.Context(), userID, amountPaise, req.IdempotencyKey)
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "withdrawal processed", fiber.Map{
		"transaction_id": result.Transaction.ID.String(),
		"status":         string(result.Transaction.Status),
		"remaining_balance": money.FormatPaise(result.FromBalance),
	})
}

func (h *WithdrawalHandler) List(c fiber.Ctx) error {
	userID := userIDFromLocals(c)
	if userID == "" {
		return response.Unauthorized(c, "authentication required")
	}
	uid, err := uuidFromString(userID)
	if err != nil {
		return response.BadRequest(c, nil, "invalid user id")
	}

	txs, total, err := h.withdrawals.ListWithdrawals(c.Context(), &uid, 20, 0)
	if err != nil {
		return handleError(c, err)
	}
	out := make([]fiber.Map, len(txs))
	for i, t := range txs {
		out[i] = fiber.Map{
			"id":         t.ID.String(),
			"amount":     money.FormatPaise(t.Amount),
			"status":     string(t.Status),
			"created_at": t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	return response.OK(c, "withdrawals listed", fiber.Map{"withdrawals": out, "total": total})
}