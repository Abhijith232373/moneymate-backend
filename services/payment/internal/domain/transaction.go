package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TxStatus string

const (
	TxStatusPending   TxStatus = "pending"
	TxStatusCompleted TxStatus = "completed"
	TxStatusFailed    TxStatus = "failed"
	TxStatusReversed  TxStatus = "reversed"
)

type Transaction struct {
	ID             uuid.UUID
	FromAccountID  uuid.UUID
	ToAccountID    uuid.UUID
	Amount         int64 // paise
	Status         TxStatus
	IdempotencyKey string
	Description    string
	CreatedAt      time.Time
	CompletedAt    *time.Time
}

type JournalEntry struct {
	ID            uuid.UUID
	TransactionID uuid.UUID
	AccountID     uuid.UUID
	Amount        int64 // paise
	Direction     string // "debit" | "credit"
	CreatedAt     time.Time
}

type TransactionRepository interface {
	Create(ctx context.Context, t *Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*Transaction, error)
	// GetByIdempotencyKey now scoped to the sending account — see A4.
	GetByIdempotencyKey(ctx context.Context, key string, fromAccountID uuid.UUID) (*Transaction, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status TxStatus) error
	ListByAccount(ctx context.Context, accountID uuid.UUID) ([]*Transaction, error)
	GetEntriesByTransactionID(ctx context.Context, txID uuid.UUID) ([]*JournalEntry, error)
}