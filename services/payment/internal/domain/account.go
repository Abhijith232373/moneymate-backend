package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AccountType string

const (
	AccountTypeWallet             AccountType = "wallet"
	AccountTypePod                AccountType = "pod"
	AccountTypeMerchantSettlement AccountType = "merchant_settlement"
	AccountTypeMerchantPayout     AccountType = "merchant_payout"
	AccountTypePlatformCommission AccountType = "platform_commission_pool"
)

type Account struct {
	ID         uuid.UUID
	UserID     *uuid.UUID
	MerchantID *uuid.UUID
	Type       AccountType
	Currency   string
	Balance    int64 // paise
	Version    int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type AccountRepository interface {
	Create(ctx context.Context, a *Account) (*Account, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Account, error)
	GetWalletByUserID(ctx context.Context, userID uuid.UUID) (*Account, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*Account, error)
}