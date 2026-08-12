package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/services/payment/internal/domain"
	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"
)

const (
	minWithdrawalPaise = 100
	maxWithdrawalPaise = 200000 * 100 
)

type WithdrawalUsecase interface {
	RequestWithdrawal(ctx context.Context, userID string, amountPaise int64, idempotencyKey string) (*domain.LedgerResult, error)
	ListWithdrawals(ctx context.Context, userID *uuid.UUID, limit, offset int32) ([]*domain.Transaction, int64, error)
}

type withdrawalUsecase struct {
	accounts                    domain.AccountRepository
	transactions                domain.TransactionRepository
	ledger                      domain.LedgerRepository
	externalSettlementAccountID uuid.UUID
}

func NewWithdrawalUsecase(
	accounts domain.AccountRepository,
	transactions domain.TransactionRepository,
	ledger domain.LedgerRepository,
	externalSettlementAccountID uuid.UUID,
) WithdrawalUsecase {
	return &withdrawalUsecase{
		accounts: accounts, transactions: transactions, ledger: ledger,
		externalSettlementAccountID: externalSettlementAccountID,
	}
}

// RequestWithdrawal debits the user's wallet immediately via the real ledger
// (so balance checks/locking/idempotency all apply exactly like a transfer),
// crediting the external_settlement account. The transaction lands as
// 'completed' in the ledger sense (money left the wallet), but actual bank
// payout is a manual/external step until RazorpayX Payouts is wired in —
// tracked outside this usecase for now.
func (u *withdrawalUsecase) RequestWithdrawal(ctx context.Context, userID string, amountPaise int64, idempotencyKey string) (*domain.LedgerResult, error) {
	if amountPaise < minWithdrawalPaise || amountPaise > maxWithdrawalPaise {
		return nil, apperrors.ErrInvalidInput
	}
	if idempotencyKey == "" {
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

	existing, err := u.transactions.GetByIdempotencyKey(ctx, idempotencyKey, wallet.ID)
	if err == nil && existing != nil {
		return u.replay(ctx, existing)
	}

	result, err := u.ledger.ExecuteTransfer(ctx, &domain.Transaction{
		FromAccountID:  wallet.ID,
		ToAccountID:    u.externalSettlementAccountID,
		Amount:         amountPaise,
		Status:         domain.TxStatusPending,
		IdempotencyKey: idempotencyKey,
		Description:    "wallet withdrawal",
	})
	if err != nil {
		return nil, fmt.Errorf("execute withdrawal: %w", err)
	}
	return result, nil
}

func (u *withdrawalUsecase) replay(ctx context.Context, t *domain.Transaction) (*domain.LedgerResult, error) {
	entries, err := u.transactions.GetEntriesByTransactionID(ctx, t.ID)
	if err != nil {
		return nil, err
	}
	result := &domain.LedgerResult{Transaction: t}
	for _, e := range entries {
		switch e.Direction {
		case "debit":
			result.DebitEntry = e
		case "credit":
			result.CreditEntry = e
		}
	}
	if from, err := u.accounts.GetByID(ctx, t.FromAccountID); err == nil {
		result.FromBalance = from.Balance
	}
	if to, err := u.accounts.GetByID(ctx, t.ToAccountID); err == nil {
		result.ToBalance = to.Balance
	}
	return result, nil
}

func (u *withdrawalUsecase) ListWithdrawals(ctx context.Context, userID *uuid.UUID, limit, offset int32) ([]*domain.Transaction, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return u.transactions.ListWithdrawals(ctx, u.externalSettlementAccountID, userID, limit, offset)
}