package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moneymate-2026/moneymate-backend/services/payment/internal/domain"
	"github.com/moneymate-2026/moneymate-backend/services/payment/sqlc/generated"
	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"
)

type LedgerRepo struct {
	pool *pgxpool.Pool
}

func NewLedgerRepo(pool *pgxpool.Pool) *LedgerRepo {
	return &LedgerRepo{pool: pool}
}

// ExecuteTransfer moves money atomically. It:
//  1. locks both accounts (in a stable, sorted order — no deadlocks)
//  2. checks the sender has enough balance
//  3. writes the transaction + two journal entries (debit + credit)
//  4. updates both balances
//  5. commits — or rolls everything back on any failure
//
// If the (idempotency_key, from_account_id) pair was already used, the INSERT
// hits the unique constraint and the whole tx rolls back cleanly — no partial
// writes. The caller (transfer_usecase) is responsible for turning that into
// a replay of the original result rather than a bare error. See A10.
func (r *LedgerRepo) ExecuteTransfer(ctx context.Context, t *domain.Transaction) (*domain.LedgerResult, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin ledger tx: %w", err)
	}
	defer tx.Rollback(ctx) // no-op if already committed

	q := generated.New(tx)

	// Lock both accounts in sorted-ID order so two concurrent transfers
	// A→B and B→A never deadlock.
	lock1, lock2 := t.FromAccountID, t.ToAccountID
	if lock2.String() < lock1.String() {
		lock1, lock2 = lock2, lock1
	}
	first, err := q.GetAccountByIDForUpdate(ctx, lock1)
	if err != nil {
		return nil, mapDBErr(err)
	}
	second, err := q.GetAccountByIDForUpdate(ctx, lock2)
	if err != nil {
		return nil, mapDBErr(err)
	}

	var from generated.GetAccountByIDForUpdateRow
	if first.ID == t.FromAccountID {
		from = first
	} else {
		from = second
	}

	if from.Balance < t.Amount {
		return nil, apperrors.ErrInsufficientFunds
	}

	now := time.Now().UTC()
	txID := uuid.New()

	inserted, err := q.InsertTransaction(ctx, generated.InsertTransactionParams{
		ID:             txID,
		FromAccountID:  t.FromAccountID,
		ToAccountID:    t.ToAccountID,
		Amount:         t.Amount,
		Column5:        generated.PaymentTxStatus(domain.TxStatusCompleted),
		IdempotencyKey: t.IdempotencyKey,
		Description:    &t.Description,
		CompletedAt:    pgtype.Timestamptz{Time: now, Valid: true},
	})
	if err != nil {
		return nil, mapDBErr(err) // unique-violation ⇒ ErrIdempotencyKeyUsed, handled by usecase
	}

	// Double entry: one DEBIT on the sender, one equal CREDIT on the receiver.
	if err := q.InsertJournalEntry(ctx, generated.InsertJournalEntryParams{
		ID:            uuid.New(),
		TransactionID: txID,
		AccountID:     t.FromAccountID,
		Amount:        t.Amount,
		Column5:       generated.PaymentTxDirectionDebit,
	}); err != nil {
		return nil, mapDBErr(err)
	}
	if err := q.InsertJournalEntry(ctx, generated.InsertJournalEntryParams{
		ID:            uuid.New(),
		TransactionID: txID,
		AccountID:     t.ToAccountID,
		Amount:        t.Amount,
		Column5:       generated.PaymentTxDirectionCredit,
	}); err != nil {
		return nil, mapDBErr(err)
	}

	// Update balances (negative = debit, positive = credit).
	if err := q.AddBalance(ctx, generated.AddBalanceParams{ID: t.FromAccountID, Balance: -t.Amount}); err != nil {
		return nil, mapDBErr(err)
	}
	if err := q.AddBalance(ctx, generated.AddBalanceParams{ID: t.ToAccountID, Balance: t.Amount}); err != nil {
		return nil, mapDBErr(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit ledger tx: %w", err)
	}

	// Read fresh balances for the response (outside the tx, so no locks held).
	result := &domain.LedgerResult{Transaction: rowToTransaction(generated.GetTransactionByIDRow(inserted))}
	if fromAcc, err := r.getAccount(ctx, t.FromAccountID); err == nil {
		result.FromBalance = fromAcc.Balance
	}
	if toAcc, err := r.getAccount(ctx, t.ToAccountID); err == nil {
		result.ToBalance = toAcc.Balance
	}
	return result, nil
}

func (r *LedgerRepo) getAccount(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return NewAccountRepo(r.pool).GetByID(ctx, id)
}