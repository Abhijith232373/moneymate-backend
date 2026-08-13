	// payment/internal/transport/http/deposit_handler.go
	package http

	import (
		"github.com/gofiber/fiber/v3"
		"github.com/google/uuid"

		"github.com/moneymate-2026/moneymate-backend/services/payment/internal/domain"
		"github.com/moneymate-2026/moneymate-backend/services/payment/internal/usecases"
		"github.com/moneymate-2026/moneymate-backend/shared/pkg/money"
		response "github.com/moneymate-2026/moneymate-backend/shared/pkg/responses"
	)

	type DepositHandler struct {
		deposits usecases.DepositUsecase
		razorpay usecases.RazorpayClient
	}

	func NewDepositHandler(deposits usecases.DepositUsecase, razorpay usecases.RazorpayClient) *DepositHandler {
		return &DepositHandler{deposits: deposits, razorpay: razorpay}
	}

	func (h *DepositHandler) Initiate(c fiber.Ctx) error {
		userID := userIDFromLocals(c)
		if userID == "" {
			return response.Unauthorized(c, "authentication required")
		}

		var req initiateDepositRequest
		if err := c.Bind().Body(&req); err != nil {
			return response.BadRequest(c, nil, "invalid request body")
		}
		if err := validate.Struct(req); err != nil {
			return response.BadRequest(c, nil, "validation failed")
		}

		amountPaise := int64(req.AmountRupees * 100)

		resp, err := h.deposits.InitiateDeposit(c.Context(), userID, amountPaise)
		if err != nil {
			return handleError(c, err)
		}
		return response.OK(c, "deposit order created", fiber.Map{
			"order_id": resp.OrderID,
			"amount":   money.FormatPaise(resp.Amount),
			"key_id":   resp.KeyID,
		})
	}

	// Confirm is called by the client after Razorpay Checkout succeeds,
	// with the three values Razorpay's checkout callback returns.
	func (h *DepositHandler) Confirm(c fiber.Ctx) error {
		userID := userIDFromLocals(c)
		if userID == "" {
			return response.Unauthorized(c, "authentication required")
		}

		var req confirmDepositRequest
		if err := c.Bind().Body(&req); err != nil {
			return response.BadRequest(c, nil, "invalid request body")
		}
		if err := validate.Struct(req); err != nil {
			return response.BadRequest(c, nil, "validation failed")
		}

		if err := h.razorpay.VerifySignature(req.OrderID, req.PaymentID, req.Signature); err != nil {
			return response.Unauthorized(c, "payment signature verification failed")
		}

		if err := h.deposits.ConfirmDeposit(c.Context(), req.OrderID, req.PaymentID); err != nil {
			return handleError(c, err)
		}
		return response.OK(c, "deposit confirmed", nil)
	}

	func (h *DepositHandler) List(c fiber.Ctx) error {
		userID := userIDFromLocals(c)
		if userID == "" {
			return response.Unauthorized(c, "authentication required")
		}
		uid, err := uuidFromString(userID)
		if err != nil {
			return response.BadRequest(c, nil, "invalid user id")
		}

		deposits, total, err := h.deposits.ListDeposits(c.Context(), nil, &uid, 20, 0)
		if err != nil {
			return handleError(c, err)
		}
		return response.OK(c, "deposits listed", fiber.Map{
			"deposits": toDepositResponses(deposits),
			"total":    total,
		})
	}

	func toDepositResponses(deposits []*domain.Deposit) []fiber.Map {
		out := make([]fiber.Map, len(deposits))
		for i, d := range deposits {
			out[i] = fiber.Map{
				"id":         d.ID.String(),
				"amount":     money.FormatPaise(d.Amount),
				"status":     string(d.Status),
				"created_at": d.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
		}
		return out
	}


	func uuidFromString(s string) (uuid.UUID, error) {
		return uuid.Parse(s)
	}