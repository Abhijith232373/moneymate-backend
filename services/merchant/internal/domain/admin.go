package domain

import (
	"context"
	"github.com/google/uuid"
)

// AdminRepository defines the strict data access contract for platform administrators
// to perform CRUD and governance operations across all merchants, campaigns, KYC records, rewards, and subscriptions.
type AdminRepository interface {
	// Stores
	GetAllStores(ctx context.Context, limit, offset int) ([]*Store, error)
	GetStoreByID(ctx context.Context, storeID uuid.UUID) (*Store, error)
	UpdateStoreStatus(ctx context.Context, storeID uuid.UUID, status string) error
	DeleteStore(ctx context.Context, storeID uuid.UUID) error

	// Campaigns
	GetAllCampaigns(ctx context.Context, limit, offset int) ([]*Campaign, error)
	CreateCampaign(ctx context.Context, campaign *Campaign) (*Campaign, error)
	GetCampaignsByStoreID(ctx context.Context, storeID uuid.UUID) ([]*Campaign, error)
	UpdateCampaignStatus(ctx context.Context, id uuid.UUID, isActive bool) error
	DeleteCampaign(ctx context.Context, id uuid.UUID) error

	// KYC Verification
	GetAllKYCDocuments(ctx context.Context, limit, offset int) ([]*KYCStatusDetail, error)
	GetKYCByStoreID(ctx context.Context, storeID uuid.UUID) (*KYCStatusDetail, error)
	VerifyKYCDocument(ctx context.Context, storeID uuid.UUID, isVerified bool, status string) (*KYCStatusDetail, error)

	// Rewards
	GetAllRewardTransactions(ctx context.Context, limit, offset int) ([]*RewardTransaction, error)
	GetPlatformRewardSummary(ctx context.Context) (float64, int64, error)

	// Subscriptions
	GetAllSubscriptions(ctx context.Context, limit, offset int) ([]*MerchantSubscription, error)
	UpdateStoreSubscriptionPlan(ctx context.Context, storeID uuid.UUID, planCode string) (*MerchantSubscription, error)
}
