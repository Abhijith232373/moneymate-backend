package http

import (
	"github.com/gofiber/fiber/v3"

	"github.com/moneymate-2026/moneymate-backend/services/payment/internal/domain"
	"github.com/moneymate-2026/moneymate-backend/services/payment/internal/usecases"
	"github.com/moneymate-2026/moneymate-backend/shared/pkg/money"
	response "github.com/moneymate-2026/moneymate-backend/shared/pkg/responses"
)

type WalletHandler struct {
	wallets usecases.WalletUsecase
}

func NewWalletHandler(wallets usecases.WalletUsecase) *WalletHandler {
	return &WalletHandler{wallets: wallets}
}
func (h *WalletHandler) CreateWallet(c fiber.Ctx) error {
	userID := userIDFromLocals(c)
	if userID == "" {
		return response.Unauthorized(c, "authentication required")
	}
	// self-service creation — handle not supplied by the client; this path
	// is effectively dead once auth-svc always creates the wallet at
	// registration, but kept for now in case a wallet is ever missing and
	// needs manual recovery.
	return response.BadRequest(c, nil, "wallets are created automatically at registration")
}

// CreateWalletInternal is called by auth-svc's Register flow, immediately
// after a new user is created. Protected by RequireInternalSecret — never
// exposed through the gateway.
func (h *WalletHandler) CreateWalletInternal(c fiber.Ctx) error {
	var req createWalletInternalRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return response.BadRequest(c, nil, "validation failed")
	}

	acc, err := h.wallets.CreateWallet(c.Context(), req.UserID, req.Handle)
	if err != nil {
		return handleError(c, err)
	}
	return response.Created(c, "wallet created", toWalletResponse(acc))
}
// GetMyWallet replaces the old GetWalletByUser(:user_id) route — a user can
// only ever look up their own wallet through this endpoint. Looking up an
// arbitrary user's wallet by path param is intentionally removed.
func (h *WalletHandler) GetMyWallet(c fiber.Ctx) error {
	userID := userIDFromLocals(c)
	if userID == "" {
		return response.Unauthorized(c, "authentication required")
	}
	acc, err := h.wallets.GetWallet(c.Context(), userID)
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "wallet found", toWalletResponse(acc))
}

func (h *WalletHandler) GetWalletByID(c fiber.Ctx) error {
	acc, err := h.wallets.GetByID(c.Context(), c.Params("id"))
	if err != nil {
		return handleError(c, err)
	}
	// Only the owner can view their own wallet by ID either.
	userID := userIDFromLocals(c)
	if acc.UserID == nil || acc.UserID.String() != userID {
		return response.Forbidden(c, nil, "you do not have access to this wallet")
	}
	return response.OK(c, "wallet found", toWalletResponse(acc))
}

func toWalletResponse(a *domain.Account) walletResponse {
	var userID string
	if a.UserID != nil {
		userID = a.UserID.String()
	}
	return walletResponse{
		ID:       a.ID.String(),
		UserID:   userID,
		Type:     string(a.Type),
		Currency: a.Currency,
		Balance:  money.FormatPaise(a.Balance),
	}
}
