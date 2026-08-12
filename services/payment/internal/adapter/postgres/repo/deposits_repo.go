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
)
type DepositRepo struct {
	pool *pgxpool.Pool
	q    *generated.Queries
}

func NewDepositRepo(pool *pgxpool.Pool) *DepositRepo {
	return &DepositRepo{pool: pool, q: generated.New(pool)}
}

func (r *DepositRepo) Create(ctx context.Context, d *domain.Deposit) error {
	row, err := r.q.CreateDeposit(ctx, generated.CreateDepositParams{
		ID:              d.ID,
		UserID:          d.UserID,
		AccountID:       d.AccountID,
		RazorpayOrderID: d.RazorpayOrderID,
		Amount:          d.Amount,
	})
	if err != nil {
		return mapDBErr(err)
	}
	d.CreatedAt = row.CreatedAt
	return nil
}

func (r *DepositRepo) GetByOrderID(ctx context.Context, orderID string) (*domain.Deposit, error) {
	row, err := r.q.GetDepositByOrderID(ctx, orderID)
	if err != nil {
		return nil, mapDBErr(err)
	}
	return toDomainDeposit(row), nil
}

// ConfirmPayment atomically marks the deposit paid, credits the wallet, and
// records a double-entry transaction (external settlement -> wallet).
// Idempotent: if the deposit isn't in 'created' status (already processed),
// it returns credited=false without touching the balance — safe against
// Razorpay's at-least-once webhook delivery.
func (r *DepositRepo) ConfirmPayment(ctx context.Context, orderID, paymentID string, externalSettlementAccountID uuid.UUID) (*domain.Deposit, bool, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, false, fmt.Errorf("begin confirm tx: %w", err)
	}
	defer tx.Rollback(ctx)

	q := generated.New(tx)

	deposit, err := q.GetDepositByOrderID(ctx, orderID)
	if err != nil {
		return nil, false, mapDBErr(err)
	}

	now := time.Now().UTC()
	affected, err := q.MarkDepositPaidIfCreated(ctx, generated.MarkDepositPaidIfCreatedParams{
		ID:                deposit.ID,
		RazorpayPaymentID: &paymentID,
		CompletedAt:       pgtype.Timestamptz{Time: now, Valid: true},
	})
	if err != nil {
		return nil, false, mapDBErr(err)
	}
	if affected == 0 {
		// Already processed (duplicate webhook) — commit the no-op and return as-is.
		if err := tx.Commit(ctx); err != nil {
			return nil, false, fmt.Errorf("commit confirm tx (no-op): %w", err)
		}
		return toDomainDeposit(deposit), false, nil
	}

	if err := q.AddBalance(ctx, generated.AddBalanceParams{ID: deposit.AccountID, Balance: deposit.Amount}); err != nil {
		return nil, false, mapDBErr(err)
	}

	txID := uuid.New()
	desc := "wallet deposit via Razorpay"
	if _, err := q.InsertTransaction(ctx, generated.InsertTransactionParams{
		ID:             txID,
		FromAccountID:  externalSettlementAccountID,
		ToAccountID:    deposit.AccountID,
		Amount:         deposit.Amount,
		Column5:        generated.PaymentTxStatus(domain.TxStatusCompleted),
		IdempotencyKey: orderID,
		Description:    &desc,
		CompletedAt:    pgtype.Timestamptz{Time: now, Valid: true},
	}); err != nil {
		return nil, false, mapDBErr(err)
	}

	if err := q.InsertJournalEntry(ctx, generated.InsertJournalEntryParams{
		ID: uuid.New(), TransactionID: txID, AccountID: externalSettlementAccountID,
		Amount: deposit.Amount, Column5: generated.PaymentTxDirectionDebit,
	}); err != nil {
		return nil, false, mapDBErr(err)
	}
	if err := q.InsertJournalEntry(ctx, generated.InsertJournalEntryParams{
		ID: uuid.New(), TransactionID: txID, AccountID: deposit.AccountID,
		Amount: deposit.Amount, Column5: generated.PaymentTxDirectionCredit,
	}); err != nil {
		return nil, false, mapDBErr(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, false, fmt.Errorf("commit confirm tx: %w", err)
	}

	deposit.Status = generated.PaymentDepositStatusPaid
	return toDomainDeposit(deposit), true, nil
}

func (r *DepositRepo) MarkFailed(ctx context.Context, orderID, paymentID string) error {
	deposit, err := r.q.GetDepositByOrderID(ctx, orderID)
	if err != nil {
		return mapDBErr(err)
	}
	return mapDBErr(r.q.MarkDepositFailed(ctx, generated.MarkDepositFailedParams{
		ID:                deposit.ID,
		RazorpayPaymentID: &paymentID,
	}))
}

func (r *DepositRepo) List(ctx context.Context, status *domain.DepositStatus, userID *uuid.UUID, limit, offset int32) ([]*domain.Deposit, int64, error) {
	var statusParam *generated.PaymentDepositStatus
	if status != nil {
		s := generated.PaymentDepositStatus(*status)
		statusParam = &s
	}
	var userIDParam pgtype.UUID
	if userID != nil {
		userIDParam = pgtype.UUID{Bytes: *userID, Valid: true}
	}

	rows, err := r.q.ListDeposits(ctx, generated.ListDepositsParams{
		Limit: limit, Offset: offset, Status: statusParam, UserID: userIDParam,
	})
	if err != nil {
		return nil, 0, mapDBErr(err)
	}
	total, err := r.q.CountDeposits(ctx, generated.CountDepositsParams{Status: statusParam, UserID: userIDParam})
	if err != nil {
		return nil, 0, mapDBErr(err)
	}
	deposits := make([]*domain.Deposit, len(rows))
	for i, row := range rows {
		deposits[i] = toDomainDeposit(row)
	}
	return deposits, total, nil
}

func toDomainDeposit(row generated.PaymentDeposit) *domain.Deposit {
	d := &domain.Deposit{
		ID: row.ID, UserID: row.UserID, AccountID: row.AccountID,
		RazorpayOrderID: row.RazorpayOrderID, Amount: row.Amount,
		Status: domain.DepositStatus(row.Status), CreatedAt: row.CreatedAt,
	}
	if row.RazorpayPaymentID != nil {
		d.RazorpayPaymentID = *row.RazorpayPaymentID
	}
	if row.CompletedAt.Valid {
		t := row.CompletedAt.Time
		d.CompletedAt = &t
	}
	return d
}