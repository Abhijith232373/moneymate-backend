package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID               uuid.UUID
	StoreID          uuid.UUID
	AvailableBalance float64
	TotalEarnings    float64
	TotalRedeemed    float64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type WalletTransaction struct {
	ID              uuid.UUID
	StoreID         uuid.UUID
	TransactionID   string
	Title           string
	Subtitle        string
	Amount          float64
	TransactionType string
	CreatedAt       time.Time
}

type WalletRepository interface {
	GetWalletByStoreID(ctx context.Context, storeID uuid.UUID) (*Wallet, error)
	UpsertWallet(ctx context.Context, storeID uuid.UUID, availableBalance, totalEarnings, totalRedeemed float64) (*Wallet, error)
	GetWalletTransactions(ctx context.Context, storeID uuid.UUID) ([]*WalletTransaction, error)
	GetWalletTransactionsByType(ctx context.Context, storeID uuid.UUID, txnType string) ([]*WalletTransaction, error)
	CreateWalletTransaction(ctx context.Context, storeID uuid.UUID, transactionID, title, subtitle, txnType string, amount float64) (*WalletTransaction, error)
}
