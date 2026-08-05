package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Campaign represents a merchant's offer or loyalty campaign.
type Campaign struct {
	ID             uuid.UUID
	StoreID        uuid.UUID
	Name           string
	OfferType      string
	RewardValue    float64
	MinBillAmount  float64
	TargetAudience string
	BannerURL      *string
	StartDate      time.Time
	EndDate        time.Time
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// CampaignRepository defines the strict data access contract for campaigns.
type CampaignRepository interface {
	CreateCampaign(ctx context.Context, campaign *Campaign) (*Campaign, error)
	GetCampaignsByStoreID(ctx context.Context, storeID uuid.UUID) ([]*Campaign, error)

	GetCampaignByID(ctx context.Context, id uuid.UUID) (*Campaign, error)
	UpdateCampaignStatus(ctx context.Context, id, storeID uuid.UUID, isActive bool) error
}
