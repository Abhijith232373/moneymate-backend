package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type QRTransaction struct {
	ID                uuid.UUID
	StoreID           uuid.UUID
	CustomerDisplayID string
	BillAmount        float64
	RewardIssued      float64
	CreatedAt         time.Time
}

type QRRepository interface {
	CreateQRTransaction(ctx context.Context, storeID uuid.UUID, customerDisplayID string, billAmount, rewardIssued float64) (*QRTransaction, error)
	GetQRTransactionsByStoreID(ctx context.Context, storeID uuid.UUID, limit, offset int) ([]*QRTransaction, error)
	GetTodayQRScanCount(ctx context.Context, storeID uuid.UUID) (int64, error)
	GetTodayQRScanVolume(ctx context.Context, storeID uuid.UUID) (float64, error)
}
