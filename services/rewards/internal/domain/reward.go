package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RecipientType string

const (
	RecipientTypeUser     RecipientType = "user"
	RecipientTypeMerchant RecipientType = "merchant"
)

type PayoutStatus string

const (
	PayoutStatusPending   PayoutStatus = "pending"
	PayoutStatusCompleted PayoutStatus = "completed"
	PayoutStatusFailed    PayoutStatus = "failed"
)

type RewardRule struct {
	ID                        uuid.UUID
	Name                      string
	MinPercentageBPS          int32
	MaxPercentageBPS          int32
	MinTransactionAmountPaise int64
	MaxPayoutAmountPaise      int64
	Active                    bool
	CreatedBy                 *uuid.UUID
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}

type RewardPayout struct {
	ID                   uuid.UUID
	TransactionID        uuid.UUID
	RecipientID          uuid.UUID
	RecipientAccountID   uuid.UUID
	RecipientType        RecipientType
	RuleID               *uuid.UUID
	OriginalAmountPaise  int64
	RewardPercentageBPS  int32
	RewardAmountPaise    int64
	Status               PayoutStatus
	PaymentTransactionID *uuid.UUID
	FailureReason        *string
	EventPayload         []byte
	CreatedAt            time.Time
	UpdatedAt            time.Time
	CompletedAt          *time.Time
}

type RewardRepository interface {
	CreateRule(ctx context.Context, rule RewardRule) (*RewardRule, error)
	ListRules(ctx context.Context, limit, offset int32) ([]*RewardRule, error)
	GetRule(ctx context.Context, id uuid.UUID) (*RewardRule, error)
	GetActiveRule(ctx context.Context) (*RewardRule, error)
	UpdateRule(ctx context.Context, rule RewardRule) (*RewardRule, error)
	DeactivateRule(ctx context.Context, id uuid.UUID) (*RewardRule, error)

	InsertPayout(ctx context.Context, payout RewardPayout) (*RewardPayout, error)
	GetPayoutByID(ctx context.Context, id uuid.UUID) (*RewardPayout, error)
	ListPayoutsByRecipient(ctx context.Context, recipientID uuid.UUID, status *PayoutStatus, limit, offset int32) ([]*RewardPayout, error)
	ListPayoutsByTransaction(ctx context.Context, transactionID uuid.UUID) ([]*RewardPayout, error)
	MarkCompleted(ctx context.Context, payoutID, paymentTransactionID uuid.UUID) (*RewardPayout, error)
	MarkFailed(ctx context.Context, payoutID uuid.UUID, reason string) (*RewardPayout, error)
	ListFailedPayouts(ctx context.Context, limit int32) ([]*RewardPayout, error)
}
