package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/services/payment/internal/domain"
	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"
)

// RazorpayClient is the contract this usecase depends on — never the concrete
// HTTP implementation. Satisfied by infra/razorpay.Client.
type RazorpayClient interface {
	CreateOrder(amountPaise int64, receipt string) (orderID string, err error)
	VerifyWebhookSignature(payload []byte, signature string) bool
}

type InitiateDepositResponse struct {
	OrderID string
	Amount  int64
	KeyID   string
}

type DepositUsecase interface {
	InitiateDeposit(ctx context.Context, userID string, amountPaise int64) (*InitiateDepositResponse, error)
	ConfirmDeposit(ctx context.Context, orderID, paymentID string) error
	FailDeposit(ctx context.Context, orderID, paymentID string) error
	ListDeposits(ctx context.Context, status *domain.DepositStatus, userID *uuid.UUID, limit, offset int32) ([]*domain.Deposit, int64, error)
}

const (
	minDepositPaise = 100    // ₹1 — Razorpay's own minimum
	maxDepositPaise = 500000 * 100 // ₹5,00,000 — sanity ceiling, tune to your compliance limits
)

type depositUsecase struct {
	deposits                    domain.DepositRepository
	accounts                    domain.AccountRepository
	razorpay                    RazorpayClient
	razorpayKeyID               string
	externalSettlementAccountID uuid.UUID
}

func NewDepositUsecase(
	deposits domain.DepositRepository,
	accounts domain.AccountRepository,
	razorpay RazorpayClient,
	razorpayKeyID string,
	externalSettlementAccountID uuid.UUID,
) DepositUsecase {
	return &depositUsecase{
		deposits: deposits, accounts: accounts, razorpay: razorpay,
		razorpayKeyID: razorpayKeyID, externalSettlementAccountID: externalSettlementAccountID,
	}
}

func (u *depositUsecase) InitiateDeposit(ctx context.Context, userID string, amountPaise int64) (*InitiateDepositResponse, error) {
	if amountPaise < minDepositPaise || amountPaise > maxDepositPaise {
		return nil, apperrors.ErrInvalidInput
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperrors.ErrInvalidInput
	}

	wallet, err := u.accounts.GetWalletByUserID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("get wallet: %w", err)
	}

	depositID := uuid.New()

	orderID, err := u.razorpay.CreateOrder(amountPaise, depositID.String())
	if err != nil {
		return nil, fmt.Errorf("create razorpay order: %w", err)
	}

	deposit := &domain.Deposit{
		ID: depositID, UserID: uid, AccountID: wallet.ID,
		RazorpayOrderID: orderID, Amount: amountPaise,
	}
	if err := u.deposits.Create(ctx, deposit); err != nil {
		return nil, fmt.Errorf("persist deposit: %w", err)
	}

	return &InitiateDepositResponse{OrderID: orderID, Amount: amountPaise, KeyID: u.razorpayKeyID}, nil
}

// ConfirmDeposit is called only from the webhook handler, after signature
// verification. It never trusts a client-supplied "I paid" claim directly.
func (u *depositUsecase) ConfirmDeposit(ctx context.Context, orderID, paymentID string) error {
	_, credited, err := u.deposits.ConfirmPayment(ctx, orderID, paymentID, u.externalSettlementAccountID)
	if err != nil {
		return fmt.Errorf("confirm deposit: %w", err)
	}
	if !credited {
		// Duplicate webhook — not an error, just nothing new to do.
		return nil
	}
	return nil
}

func (u *depositUsecase) FailDeposit(ctx context.Context, orderID, paymentID string) error {
	if err := u.deposits.MarkFailed(ctx, orderID, paymentID); err != nil {
		return fmt.Errorf("mark deposit failed: %w", err)
	}
	return nil
}

func (u *depositUsecase) ListDeposits(ctx context.Context, status *domain.DepositStatus, userID *uuid.UUID, limit, offset int32) ([]*domain.Deposit, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20 // sane default + hard ceiling — prevents an unbounded query from a bad/malicious client
	}
	return u.deposits.List(ctx, status, userID, limit, offset)
}