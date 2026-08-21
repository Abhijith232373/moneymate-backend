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

	// Dashboard
	GetAdminDashboardStats(ctx context.Context) (*AdminDashboardStats, error)
}

type AdminDashboardStats struct {
	TotalRevenue       float64                  `json:"total_revenue"`
	RewardPool         float64                  `json:"reward_pool"`
	SystemWallet       float64                  `json:"system_wallet"`
	DailyTransactions  int64                    `json:"daily_transactions"`
	RecentTransactions []AdminRecentTransaction `json:"recent_transactions"`
}

type AdminRecentTransaction struct {
	ID        string  `json:"id"`
	User      string  `json:"user"`
	Amount    string  `json:"amount"`
	Date      string  `json:"date"`
	Status    string  `json:"status"`
}
