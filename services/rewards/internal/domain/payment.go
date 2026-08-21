package domain

import (
	"context"

	"github.com/google/uuid"
)

type PaymentClient interface {
	ExecuteRewardPayout(ctx context.Context, recipientAccountID uuid.UUID, amountPaise int64) (txID uuid.UUID, err error)
}
