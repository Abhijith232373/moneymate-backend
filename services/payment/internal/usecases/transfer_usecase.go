package usecases

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/services/payment/internal/domain"
	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"
)

type TransferInput struct {
	AuthenticatedUserID string // from the verified access token — NOT client-supplied
	FromAccountID        string
	ToAccountID           string
	AmountPaise           int64 // already converted by the handler
	IdempotencyKey        string
	Description           string
}

type TransferUsecase interface {
	Transfer(ctx context.Context, in TransferInput) (*domain.LedgerResult, error)
	GetByID(ctx context.Context, id string) (*domain.Transaction, error)
}

type transferUsecase struct {
	accounts     domain.AccountRepository
	transactions domain.TransactionRepository
	ledger       domain.LedgerRepository
}

func NewTransferUsecase(
	accounts domain.AccountRepository,
	transactions domain.TransactionRepository,
	ledger domain.LedgerRepository,
) TransferUsecase {
	return &transferUsecase{accounts: accounts, transactions: transactions, ledger: ledger}
}

func (u *transferUsecase) Transfer(ctx context.Context, in TransferInput) (*domain.LedgerResult, error) {
	// ── Validation ────────────────────────────────────────────────
	if in.AmountPaise <= 0 {
		return nil, apperrors.ErrInvalidInput
	}
	key := strings.TrimSpace(in.IdempotencyKey)
	if key == "" {
		return nil, apperrors.ErrInvalidInput
	}
	authUserID, err := uuid.Parse(in.AuthenticatedUserID)
	if err != nil {
		return nil, apperrors.ErrUnauthorized
	}
	fromID, err := uuid.Parse(in.FromAccountID)
	if err != nil {
		return nil, apperrors.ErrInvalidInput
	}
	toID, err := uuid.Parse(in.ToAccountID)
	if err != nil {
		return nil, apperrors.ErrInvalidInput
	}
	if fromID == toID {
		return nil, apperrors.ErrInvalidInput
	}

	// ── Ownership check: the caller must own the source wallet ───
	// This is the domain-layer enforcement of Decision 3. It runs regardless
	// of which transport (HTTP today, gRPC/internal later) calls this usecase.
	fromAcc, err := u.accounts.GetByID(ctx, fromID)
	if err != nil {
		return nil, err
	}
	if fromAcc.Type != domain.AccountTypeWallet {
		return nil, apperrors.ErrInvalidInput
	}
	if fromAcc.UserID == nil || *fromAcc.UserID != authUserID {
		return nil, apperrors.ErrForbidden
	}

	toAcc, err := u.accounts.GetByID(ctx, toID)
	if err != nil {
		return nil, err
	}
	if toAcc.Type != domain.AccountTypeWallet {
		return nil, apperrors.ErrInvalidInput
	}

	// ── Idempotency: replay, never double-charge ─────────────────
	existing, err := u.transactions.GetByIdempotencyKey(ctx, key, fromID)
	if err == nil && existing != nil {
		return u.replay(ctx, existing)
	}
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, err
	}

	// ── Atomic double-entry write ────────────────────────────────
	result, err := u.ledger.ExecuteTransfer(ctx, &domain.Transaction{
		FromAccountID:  fromID,
		ToAccountID:    toID,
		Amount:         in.AmountPaise,
		Status:         domain.TxStatusPending,
		IdempotencyKey: key,
		Description:    in.Description,
	})
	if err != nil {
		if errors.Is(err, apperrors.ErrIdempotencyKeyUsed) {
			// Lost a race against a concurrent identical request. The other
			// request's insert won; fetch its result and hand back the same
			// success response instead of an error to a legitimately-retrying
			// client.
			winner, getErr := u.transactions.GetByIdempotencyKey(ctx, key, fromID)
			if getErr != nil {
				return nil, getErr
			}
			return u.replay(ctx, winner)
		}
		return nil, err
	}
	return result, nil
}

// replay rebuilds a LedgerResult for an already-processed idempotency key so
// the caller gets the original outcome instead of a new transfer.
func (u *transferUsecase) replay(ctx context.Context, t *domain.Transaction) (*domain.LedgerResult, error) {
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

func (u *transferUsecase) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	txID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperrors.ErrInvalidInput
	}
	return u.transactions.GetByID(ctx, txID)
}