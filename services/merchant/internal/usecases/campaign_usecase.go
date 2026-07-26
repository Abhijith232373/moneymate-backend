package usecases

import (
	"context"

	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/domain"
)

type CampaignUseCase interface {
	CreateCampaign(ctx context.Context, campaign *domain.Campaign) (*domain.Campaign, error)
	GetCampaignsByStoreID(ctx context.Context, storeID uuid.UUID) ([]*domain.Campaign, error)
	GetCampaignsByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*domain.Campaign, error)
	UpdateCampaignStatus(ctx context.Context, id, storeID uuid.UUID, isActive bool) error
}

type campaignUseCase struct {
	campaignRepo domain.CampaignRepository
	storeRepo    domain.MerchantRepository
}

func NewCampaignUseCase(cr domain.CampaignRepository, sr domain.MerchantRepository) CampaignUseCase {
	return &campaignUseCase{
		campaignRepo: cr,
		storeRepo:    sr,
	}
}

func (uc *campaignUseCase) CreateCampaign(ctx context.Context, c *domain.Campaign) (*domain.Campaign, error) {
	// Optional: validate store exists or is active
	c.IsActive = true
	return uc.campaignRepo.CreateCampaign(ctx, c)
}

func (uc *campaignUseCase) GetCampaignsByStoreID(ctx context.Context, storeID uuid.UUID) ([]*domain.Campaign, error) {
	return uc.campaignRepo.GetCampaignsByStoreID(ctx, storeID)
}

func (uc *campaignUseCase) GetCampaignsByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*domain.Campaign, error) {
	return uc.campaignRepo.GetCampaignsByOwnerID(ctx, ownerID)
}

func (uc *campaignUseCase) UpdateCampaignStatus(ctx context.Context, id, storeID uuid.UUID, isActive bool) error {
	return uc.campaignRepo.UpdateCampaignStatus(ctx, id, storeID, isActive)
}
