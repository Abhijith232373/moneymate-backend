package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type EarningsStats struct {
	StoreID     uuid.UUID
	TotalScans  int64
	TotalEarned float64
	UpdatedAt   time.Time
}

type EarningsPayoutRequest struct {
	ID             uuid.UUID
	StoreID        uuid.UUID
	MilestoneScans int32
	RewardAmount   float64
	Status         string
	CreatedAt      time.Time
}

type EarningsRepository interface {
	GetEarningsStats(ctx context.Context, storeID uuid.UUID) (*EarningsStats, error)
	UpsertEarningsStats(ctx context.Context, storeID uuid.UUID, totalScans int64, totalEarned float64) (*EarningsStats, error)
	GetRequestedMilestones(ctx context.Context, storeID uuid.UUID) ([]int32, error)
	CreatePayoutRequest(ctx context.Context, storeID uuid.UUID, milestoneScans int32, rewardAmount float64) (*EarningsPayoutRequest, error)
}
