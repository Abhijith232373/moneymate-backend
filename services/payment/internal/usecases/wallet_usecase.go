package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/services/payment/internal/domain"
	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"
)

type WalletUsecase interface {
	CreateWallet(ctx context.Context, userID string) (*domain.Account, error)
	GetWallet(ctx context.Context, userID string) (*domain.Account, error)
	GetByID(ctx context.Context, id string) (*domain.Account, error)
}

type walletUsecase struct {
	accounts domain.AccountRepository
}

func NewWalletUsecase(accounts domain.AccountRepository) WalletUsecase {
	return &walletUsecase{accounts: accounts}
}

// CreateWallet is idempotent: a user only ever gets one wallet.
func (u *walletUsecase) CreateWallet(ctx context.Context, userID string) (*domain.Account, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperrors.ErrInvalidInput
	}

	existing, err := u.accounts.GetWalletByUserID(ctx, uid)
	if err == nil && existing != nil {
		return existing, nil // already has a wallet — return it
	}
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, err
	}

	created, err := u.accounts.Create(ctx, &domain.Account{
		UserID:   &uid,
		Type:     domain.AccountTypeWallet,
		Currency: "INR",
	})
	if err != nil {
		if errors.Is(err, apperrors.ErrAlreadyExists) {
			// Lost a race — another request created it first. Return the winner.
			return u.accounts.GetWalletByUserID(ctx, uid)
		}
		return nil, fmt.Errorf("create wallet: %w", err)
	}
	return created, nil
}

func (u *walletUsecase) GetWallet(ctx context.Context, userID string) (*domain.Account, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperrors.ErrInvalidInput
	}
	return u.accounts.GetWalletByUserID(ctx, uid)
}

func (u *walletUsecase) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	accID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperrors.ErrInvalidInput
	}
	return u.accounts.GetByID(ctx, accID)
}