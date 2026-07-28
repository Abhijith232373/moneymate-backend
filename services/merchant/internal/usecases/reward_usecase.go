package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/domain"
)

// RewardUseCase defines the high-level business logic and validation contract for the Rewards Center.
// It orchestrates financial summaries, transaction history querying, and bank withdrawal requests.
type RewardUseCase interface {
	// GetRewardSummary fetches the current balance, scan counters, and growth percentage for a merchant store.
	GetRewardSummary(ctx context.Context, storeID uuid.UUID) (*domain.RewardSummary, error)
	// GetRewardHistory retrieves paginated reward transactions with filtering by time period and keyword search.
	GetRewardHistory(ctx context.Context, storeID uuid.UUID, filter string, searchQuery string, limit, offset int) ([]*domain.RewardTransaction, error)
	// RedeemRewards validates authorization and fund availability, executes a payout deduction, and returns the withdrawal request and updated balance summary.
	RedeemRewards(ctx context.Context, storeID uuid.UUID, amount float64, bankAuthorized bool) (*domain.RedemptionRequest, *domain.RewardSummary, error)
}

// rewardUseCase implements the RewardUseCase interface using domain repositories.
type rewardUseCase struct {
	// rewardRepo provides persistent data access to reward balances and ledger tables.
	rewardRepo domain.RewardRepository
	// storeRepo verifies merchant store existence and compliance status before financial operations.
	storeRepo domain.MerchantRepository
}

// NewRewardUseCase constructs and returns a new rewardUseCase instance with injected dependencies.
func NewRewardUseCase(rr domain.RewardRepository, sr domain.MerchantRepository) RewardUseCase {
	return &rewardUseCase{
		rewardRepo: rr,
		storeRepo:  sr,
	}
}

// GetRewardSummary validates the store identifier and retrieves the financial metrics from the repository layer.
func (uc *rewardUseCase) GetRewardSummary(ctx context.Context, storeID uuid.UUID) (*domain.RewardSummary, error) {
	if storeID == uuid.Nil {
		return nil, errors.New("invalid store ID")
	}
	return uc.rewardRepo.GetSummaryByStoreID(ctx, storeID)
}

// GetRewardHistory validates input pagination bounds and delegates dynamic transaction search to the repository layer.
func (uc *rewardUseCase) GetRewardHistory(ctx context.Context, storeID uuid.UUID, filter string, searchQuery string, limit, offset int) ([]*domain.RewardTransaction, error) {
	if storeID == uuid.Nil {
		return nil, errors.New("invalid store ID")
	}
	if limit <= 0 || limit > 100 {
		limit = 50 // Default safe page size for high-throughput production loads
	}
	if offset < 0 {
		offset = 0
	}
	return uc.rewardRepo.GetTransactionsByStoreID(ctx, storeID, filter, searchQuery, limit, offset)
}

// RedeemRewards validates bank transfer authorization, ensures requested amount is positive, executes the atomic payout transaction,
// and retrieves the updated financial summary to reflect the new available balance immediately.
func (uc *rewardUseCase) RedeemRewards(ctx context.Context, storeID uuid.UUID, amount float64, bankAuthorized bool) (*domain.RedemptionRequest, *domain.RewardSummary, error) {
	if storeID == uuid.Nil {
		return nil, nil, errors.New("invalid store ID")
	}
	if !bankAuthorized {
		return nil, nil, errors.New("bank transfer authorization is mandatory before redeeming funds")
	}

	// If amount is 0 or not provided, let's fetch the available balance to redeem the full balance!
	if amount <= 0 {
		summary, err := uc.rewardRepo.GetSummaryByStoreID(ctx, storeID)
		if err != nil {
			return nil, nil, fmt.Errorf("fetch balance for full redemption: %w", err)
		}
		if summary.AvailableBalance <= 0 {
			return nil, nil, errors.New("available balance is zero; nothing to redeem")
		}
		amount = summary.AvailableBalance
	}

	req, err := uc.rewardRepo.RedeemBalance(ctx, storeID, amount, bankAuthorized)
	if err != nil {
		return nil, nil, fmt.Errorf("process redemption: %w", err)
	}

	// Retrieve updated summary after deduction
	updatedSummary, err := uc.rewardRepo.GetSummaryByStoreID(ctx, storeID)
	if err != nil {
		return req, nil, fmt.Errorf("redemption succeeded but failed to load updated summary: %w", err)
	}

	return req, updatedSummary, nil
}
