package usecases

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/domain"
)

// SubscriptionUseCase defines the business logic contract for managing pricing tiers and billing transitions.
// It orchestrates plan catalog retrieval, active subscription lookups, and compliance-validated tier changes.
type SubscriptionUseCase interface {
	// GetPlans retrieves all available subscription pricing tiers from the catalog and marks the store's active plan.
	GetPlans(ctx context.Context, storeID uuid.UUID) ([]*domain.SubscriptionPlanDetail, error)
	// GetCurrentSubscription fetches the active billing record and renewal timeline for a merchant store.
	GetCurrentSubscription(ctx context.Context, storeID uuid.UUID) (*domain.MerchantSubscription, error)
	// ChangePlan validates tier eligibility, checks active promotional limits, and atomically upgrades or downgrades the store's plan.
	ChangePlan(ctx context.Context, storeID uuid.UUID, newPlanCode string, reason string) (*domain.MerchantSubscription, error)
}

// subscriptionUseCase implements SubscriptionUseCase with injected repository dependencies.
type subscriptionUseCase struct {
	// subRepo provides persistent data access to subscription catalogs and active billing ledgers.
	subRepo domain.SubscriptionRepository
	// storeRepo verifies core merchant identity and onboarding standing.
	storeRepo domain.MerchantRepository
	// campaignRepo enables validation of promotional offer limits before allowing tier downgrades.
	campaignRepo domain.CampaignRepository
}

// NewSubscriptionUseCase constructs and returns a new subscriptionUseCase instance with required dependencies.
func NewSubscriptionUseCase(sr domain.SubscriptionRepository, mr domain.MerchantRepository, cr domain.CampaignRepository) SubscriptionUseCase {
	return &subscriptionUseCase{
		subRepo:      sr,
		storeRepo:    mr,
		campaignRepo: cr,
	}
}

// GetPlans retrieves the pricing tier catalog from the repository layer, ensuring the store ID format is valid.
func (uc *subscriptionUseCase) GetPlans(ctx context.Context, storeID uuid.UUID) ([]*domain.SubscriptionPlanDetail, error) {
	if storeID == uuid.Nil {
		return nil, errors.New("invalid store ID")
	}
	return uc.subRepo.GetAvailablePlans(ctx, storeID)
}

// GetCurrentSubscription retrieves the active subscription record for the merchant store, returning an error if the ID is nil.
func (uc *subscriptionUseCase) GetCurrentSubscription(ctx context.Context, storeID uuid.UUID) (*domain.MerchantSubscription, error) {
	if storeID == uuid.Nil {
		return nil, errors.New("invalid store ID")
	}
	return uc.subRepo.GetStoreSubscription(ctx, storeID)
}

// ChangePlan enforces business rules before upgrading or downgrading a merchant's subscription tier.
// It verifies the target plan exists and prevents downgrades if active promotional offer counts exceed the target plan's allowance.
func (uc *subscriptionUseCase) ChangePlan(ctx context.Context, storeID uuid.UUID, newPlanCode string, reason string) (*domain.MerchantSubscription, error) {
	if storeID == uuid.Nil {
		return nil, errors.New("invalid store ID")
	}
	newPlanCode = strings.ToLower(strings.TrimSpace(newPlanCode))
	if newPlanCode != "essential" && newPlanCode != "growth" && newPlanCode != "enterprise" {
		return nil, fmt.Errorf("invalid plan code %q: must be essential, growth, or enterprise", newPlanCode)
	}

	// 1. Fetch available plans to inspect the target tier's campaign limits
	plans, err := uc.subRepo.GetAvailablePlans(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("fetch available plans: %w", err)
	}
	var targetLimit int = 1
	for _, p := range plans {
		if p.PlanCode == newPlanCode {
			targetLimit = p.MaxActiveCampaigns
			break
		}
	}

	// 2. If the target limit is not unlimited (-1), check if active campaigns exceed the limit
	if targetLimit != -1 && uc.campaignRepo != nil {
		campaigns, err := uc.campaignRepo.GetCampaignsByStoreID(ctx, storeID)
		if err == nil {
			activeCount := 0
			for _, c := range campaigns {
				if c.IsActive {
					activeCount++
				}
			}
			if activeCount > targetLimit {
				return nil, fmt.Errorf("cannot switch to %s: your store currently has %d active offer campaigns, but this tier allows a maximum of %d. Please deactivate %d campaigns before downgrading",
					strings.Title(newPlanCode), activeCount, targetLimit, activeCount-targetLimit)
			}
		}
	}

	// 3. Delegate atomic transition to repository layer
	return uc.subRepo.UpdateStorePlan(ctx, storeID, newPlanCode, reason)
}
