package domain

import (
	"context"
)

// LedgerResult is what a transfer returns: the transaction, its two journal
// entries, and the resulting balances (paise) on both accounts.
type LedgerResult struct {
	Transaction *Transaction
	DebitEntry  *JournalEntry
	CreditEntry *JournalEntry
	FromBalance int64 // paise after the transfer
	ToBalance   int64 // paise after the transfer
}

// LedgerRepository performs the atomic double-entry write. It MUST be a single
// database transaction so a failure never leaves a half-written ledger.
type LedgerRepository interface {
	ExecuteTransfer(ctx context.Context, t *Transaction) (*LedgerResult, error)
}