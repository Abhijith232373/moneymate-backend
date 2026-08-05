package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/domain"
)

type WalletUseCase interface {
	GetWalletData(ctx context.Context, storeID uuid.UUID) (*domain.Wallet, error)
	GetWalletTransactions(ctx context.Context, storeID uuid.UUID, filterType string) ([]*domain.WalletTransaction, error)
}

type walletUseCase struct {
	walletRepo domain.WalletRepository
}

func NewWalletUseCase(wr domain.WalletRepository) WalletUseCase {
	return &walletUseCase{
		walletRepo: wr,
	}
}

func (uc *walletUseCase) GetWalletData(ctx context.Context, storeID uuid.UUID) (*domain.Wallet, error) {
	if storeID == uuid.Nil {
		return nil, errors.New("invalid store ID")
	}
	return uc.walletRepo.GetWalletByStoreID(ctx, storeID)
}

func (uc *walletUseCase) GetWalletTransactions(ctx context.Context, storeID uuid.UUID, filterType string) ([]*domain.WalletTransaction, error) {
	if storeID == uuid.Nil {
		return nil, errors.New("invalid store ID")
	}

	if filterType == "all" || filterType == "" {
		return uc.walletRepo.GetWalletTransactions(ctx, storeID)
	}

	// Maps frontend filters to database enum values
	txnType := ""
	switch filterType {
	case "qr_scanned":
		txnType = "qr_scan"
	case "redeemed":
		txnType = "redeem"
	default:
		return nil, errors.New("invalid filter type")
	}

	return uc.walletRepo.GetWalletTransactionsByType(ctx, storeID, txnType)
}
