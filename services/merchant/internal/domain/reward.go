package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// RewardSummary represents the financial balance and scan statistics for a merchant store in the Rewards Center.
type RewardSummary struct {
	ID                     uuid.UUID
	StoreID                uuid.UUID
	AvailableBalance       float64
	TotalScans             int64
	PremiumPoints          int64
	WeeklyGrowthPercentage float64
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// RewardTransaction represents an immutable ledger entry of a QR scan earning or payout redemption.
type RewardTransaction struct {
	ID              uuid.UUID
	StoreID         uuid.UUID
	CampaignName    string
	DisplayID       string
	Status          string
	Amount          float64
	TransactionType string
	CreatedAt       time.Time
}

// RedemptionRequest represents an audit trail for a merchant withdrawing funds to their verified bank account.
type RedemptionRequest struct {
	ID                     uuid.UUID
	StoreID                uuid.UUID
	Amount                 float64
	BankTransferAuthorized bool
	Status                 string
	ReferenceID            *string
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// RewardRepository defines the data access contract for managing high-concurrency rewards and financial ledgers.
type RewardRepository interface {
	// GetSummaryByStoreID retrieves the current reward balance and quick stats for a given merchant store.
	GetSummaryByStoreID(ctx context.Context, storeID uuid.UUID) (*RewardSummary, error)
	// InitializeSummary creates a default 0-balance reward record for a new merchant store upon onboarding.
	InitializeSummary(ctx context.Context, storeID uuid.UUID) (*RewardSummary, error)
	// GetTransactionsByStoreID queries paginated transaction history with filtering by time range ("all", "this_month", "this_week") and keyword search.
	GetTransactionsByStoreID(ctx context.Context, storeID uuid.UUID, filter string, searchQuery string, limit, offset int) ([]*RewardTransaction, error)
	// RedeemBalance atomically verifies funds, deducts the requested amount from available balance, and logs a redemption request and ledger transaction.
	RedeemBalance(ctx context.Context, storeID uuid.UUID, amount float64, bankAuthorized bool) (*RedemptionRequest, error)
}
