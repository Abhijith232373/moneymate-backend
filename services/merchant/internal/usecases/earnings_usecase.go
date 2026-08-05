package usecases

import (
	"context"
	"errors"
	"strings"
	"fmt"

	"github.com/google/uuid"
	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/domain"
)

type EarningsUseCase interface {
	GetEarningsData(ctx context.Context, storeID uuid.UUID) (*domain.EarningsStats, map[int32]bool, error)
	RequestPayout(ctx context.Context, storeID uuid.UUID, milestoneScans int32, rewardAmount float64) (*domain.EarningsPayoutRequest, error)
}

type earningsUseCase struct {
	earningsRepo domain.EarningsRepository
	walletRepo   domain.WalletRepository
}

func NewEarningsUseCase(er domain.EarningsRepository, wr domain.WalletRepository) EarningsUseCase {
	return &earningsUseCase{
		earningsRepo: er,
		walletRepo:   wr,
	}
}

func (uc *earningsUseCase) GetEarningsData(ctx context.Context, storeID uuid.UUID) (*domain.EarningsStats, map[int32]bool, error) {
	if storeID == uuid.Nil {
		return nil, nil, errors.New("invalid store ID")
	}
	
	stats, err := uc.earningsRepo.GetEarningsStats(ctx, storeID)
	if err != nil {
		return nil, nil, err
	}

	requestedList, err := uc.earningsRepo.GetRequestedMilestones(ctx, storeID)
	if err != nil {
		return nil, nil, err
	}

	requestedMap := make(map[int32]bool)
	for _, m := range requestedList {
		requestedMap[m] = true
	}

	return stats, requestedMap, nil
}

func (uc *earningsUseCase) RequestPayout(ctx context.Context, storeID uuid.UUID, milestoneScans int32, rewardAmount float64) (*domain.EarningsPayoutRequest, error) {
	if storeID == uuid.Nil {
		return nil, errors.New("invalid store ID")
	}
	
	// Create request (this acts as a record that they redeemed this milestone)
	req, err := uc.earningsRepo.CreatePayoutRequest(ctx, storeID, milestoneScans, rewardAmount)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") || strings.Contains(err.Error(), "unique constraint") {
			return nil, errors.New("reward already redeemed for this milestone")
		}
		return nil, err
	}

	// Update total earned stats in earnings
	stats, err := uc.earningsRepo.GetEarningsStats(ctx, storeID)
	if err == nil && stats != nil {
		newTotalEarned := stats.TotalEarned + rewardAmount
		_, _ = uc.earningsRepo.UpsertEarningsStats(ctx, storeID, stats.TotalScans, newTotalEarned)
	}

	// Instantly add to wallet
	if uc.walletRepo != nil {
		wallet, err := uc.walletRepo.GetWalletByStoreID(ctx, storeID)
		if err != nil {
			// If wallet doesn't exist, start fresh
			wallet = &domain.Wallet{
				AvailableBalance: 0,
				TotalEarnings:    0,
				TotalRedeemed:    0,
			}
		}

		newBalance := wallet.AvailableBalance + rewardAmount
		newEarnings := wallet.TotalEarnings + rewardAmount

		_, err = uc.walletRepo.UpsertWallet(ctx, storeID, newBalance, newEarnings, wallet.TotalRedeemed)
		if err == nil {
			txnID := fmt.Sprintf("RWD-%d-%s", milestoneScans, uuid.New().String()[:8])
			title := fmt.Sprintf("%d Scans Reward", milestoneScans)
			subtitle := "Redeem"
			_, _ = uc.walletRepo.CreateWalletTransaction(ctx, storeID, txnID, title, subtitle, "redeem", rewardAmount)
		}
	}

	return req, nil
}
