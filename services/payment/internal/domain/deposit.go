// services/payment/internal/domain/deposit.go
package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type DepositStatus string

const (
	DepositStatusCreated DepositStatus = "created"
	DepositStatusPaid    DepositStatus = "paid"
	DepositStatusFailed  DepositStatus = "failed"
)

type Deposit struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	AccountID         uuid.UUID
	RazorpayOrderID   string
	RazorpayPaymentID string
	Amount            int64
	Status            DepositStatus
	CreatedAt         time.Time
	CompletedAt       *time.Time
}

type DepositRepository interface {
	Create(ctx context.Context, d *Deposit) error
	GetByOrderID(ctx context.Context, orderID string) (*Deposit, error)
	ConfirmPayment(ctx context.Context, orderID, paymentID string, externalSettlementAccountID uuid.UUID) (*Deposit, bool, error)
	MarkFailed(ctx context.Context, orderID, paymentID string) error
	List(ctx context.Context, status *DepositStatus, userID *uuid.UUID, limit, offset int32) ([]*Deposit, int64, error)
}