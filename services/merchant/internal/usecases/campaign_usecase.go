package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/domain"
)

type CampaignUseCase interface {
	CreateCampaign(ctx context.Context, campaign *domain.Campaign) (*domain.Campaign, error)
	GetCampaignsByStoreID(ctx context.Context, storeID uuid.UUID) ([]*domain.Campaign, error)

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
	// Fetch store to get the current plan
	store, err := uc.storeRepo.GetStoreProfileByStoreID(ctx, c.StoreID)
	if err != nil {
		return nil, err
	}

	// Fetch existing campaigns
	campaigns, err := uc.campaignRepo.GetCampaignsByStoreID(ctx, c.StoreID)
	if err != nil {
		return nil, err
	}

	activeCount := 0
	for _, camp := range campaigns {
		if camp.IsActive {
			activeCount++
		}
	}

	// Enforce limits based on plan
	// Default to Essential limit if plan is empty or unknown
	switch store.Plan {
	case "Enterprise":
		// Unlimited
	case "Growth":
		if activeCount >= 5 {
			return nil, fmt.Errorf("active campaign limit reached for Growth plan (max 5)")
		}
	default:
		// Essential or empty plan
		if activeCount >= 1 {
			return nil, fmt.Errorf("active campaign limit reached for Essential plan (max 1)")
		}
	}

	c.IsActive = true
	return uc.campaignRepo.CreateCampaign(ctx, c)
}

func (uc *campaignUseCase) GetCampaignsByStoreID(ctx context.Context, storeID uuid.UUID) ([]*domain.Campaign, error) {
	return uc.campaignRepo.GetCampaignsByStoreID(ctx, storeID)
}


func (uc *campaignUseCase) UpdateCampaignStatus(ctx context.Context, id, storeID uuid.UUID, isActive bool) error {
	if isActive {
		// Fetch store to get the current plan
		store, err := uc.storeRepo.GetStoreProfileByStoreID(ctx, storeID)
		if err != nil {
			return err
		}

		// Fetch existing campaigns
		campaigns, err := uc.campaignRepo.GetCampaignsByStoreID(ctx, storeID)
		if err != nil {
			return err
		}

		activeCount := 0
		for _, camp := range campaigns {
			// Count active ones excluding the one we are updating, though if it's already active this won't change the count
			if camp.IsActive && camp.ID != id {
				activeCount++
			}
		}

		// Enforce limits
		switch store.Plan {
		case "Enterprise":
			// Unlimited
		case "Growth":
			if activeCount >= 5 {
				return fmt.Errorf("active campaign limit reached for Growth plan (max 5)")
			}
		default:
			if activeCount >= 1 {
				return fmt.Errorf("active campaign limit reached for Essential plan (max 1)")
			}
		}
	}

	return uc.campaignRepo.UpdateCampaignStatus(ctx, id, storeID, isActive)
}
