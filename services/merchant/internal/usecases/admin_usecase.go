package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/domain"
)

// AdminUseCase encapsulates business logic and governance for platform administrators managing merchant resources.
type AdminUseCase struct {
	adminRepo domain.AdminRepository
}

func NewAdminUseCase(ar domain.AdminRepository) *AdminUseCase {
	return &AdminUseCase{
		adminRepo: ar,
	}
}

// Stores
func (uc *AdminUseCase) GetAllStores(ctx context.Context, limit, offset int) ([]*domain.Store, error) {
	return uc.adminRepo.GetAllStores(ctx, limit, offset)
}

func (uc *AdminUseCase) GetStoreByID(ctx context.Context, storeIDStr string) (*domain.Store, error) {
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return nil, err
	}
	return uc.adminRepo.GetStoreByID(ctx, storeID)
}

func (uc *AdminUseCase) UpdateStoreStatus(ctx context.Context, storeIDStr string, status string) error {
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return err
	}
	return uc.adminRepo.UpdateStoreStatus(ctx, storeID, status)
}

func (uc *AdminUseCase) DeleteStore(ctx context.Context, storeIDStr string) error {
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return err
	}
	return uc.adminRepo.DeleteStore(ctx, storeID)
}

// Campaigns
func (uc *AdminUseCase) GetAllCampaigns(ctx context.Context, limit, offset int) ([]*domain.Campaign, error) {
	return uc.adminRepo.GetAllCampaigns(ctx, limit, offset)
}

func (uc *AdminUseCase) CreateCampaign(ctx context.Context, campaign *domain.Campaign) (*domain.Campaign, error) {
	return uc.adminRepo.CreateCampaign(ctx, campaign)
}

func (uc *AdminUseCase) GetCampaignsByStoreID(ctx context.Context, storeIDStr string) ([]*domain.Campaign, error) {
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return nil, err
	}
	return uc.adminRepo.GetCampaignsByStoreID(ctx, storeID)
}

func (uc *AdminUseCase) UpdateCampaignStatus(ctx context.Context, idStr string, isActive bool) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	return uc.adminRepo.UpdateCampaignStatus(ctx, id, isActive)
}

func (uc *AdminUseCase) DeleteCampaign(ctx context.Context, idStr string) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}
	return uc.adminRepo.DeleteCampaign(ctx, id)
}

// KYC Verification
func (uc *AdminUseCase) GetAllKYCDocuments(ctx context.Context, limit, offset int) ([]*domain.KYCStatusDetail, error) {
	return uc.adminRepo.GetAllKYCDocuments(ctx, limit, offset)
}

func (uc *AdminUseCase) GetKYCByStoreID(ctx context.Context, storeIDStr string) (*domain.KYCStatusDetail, error) {
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return nil, err
	}
	return uc.adminRepo.GetKYCByStoreID(ctx, storeID)
}

func (uc *AdminUseCase) VerifyKYCDocument(ctx context.Context, storeIDStr string, isVerified bool, status string) (*domain.KYCStatusDetail, error) {
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return nil, err
	}
	return uc.adminRepo.VerifyKYCDocument(ctx, storeID, isVerified, status)
}

// Rewards
func (uc *AdminUseCase) GetAllRewardTransactions(ctx context.Context, limit, offset int) ([]*domain.RewardTransaction, error) {
	return uc.adminRepo.GetAllRewardTransactions(ctx, limit, offset)
}

func (uc *AdminUseCase) GetPlatformRewardSummary(ctx context.Context) (map[string]interface{}, error) {
	bal, scans, err := uc.adminRepo.GetPlatformRewardSummary(ctx)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"total_available_balance": bal,
		"total_platform_scans":    scans,
	}, nil
}

// Subscriptions
func (uc *AdminUseCase) GetAllSubscriptions(ctx context.Context, limit, offset int) ([]*domain.MerchantSubscription, error) {
	return uc.adminRepo.GetAllSubscriptions(ctx, limit, offset)
}

func (uc *AdminUseCase) UpdateStoreSubscriptionPlan(ctx context.Context, storeIDStr string, planCode string) (*domain.MerchantSubscription, error) {
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return nil, err
	}
	return uc.adminRepo.UpdateStoreSubscriptionPlan(ctx, storeID, planCode)
}

// Dashboard
func (uc *AdminUseCase) GetAdminDashboardStats(ctx context.Context) (*domain.AdminDashboardStats, error) {
	return uc.adminRepo.GetAdminDashboardStats(ctx)
}

