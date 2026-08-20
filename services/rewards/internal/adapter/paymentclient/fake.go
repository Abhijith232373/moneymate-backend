package paymentclient

import (
	"context"
	"log"

	"github.com/google/uuid"
)

type FakeClient struct{}

func NewFakeClient() *FakeClient {
	return &FakeClient{}
}

func (c *FakeClient) ExecuteRewardPayout(ctx context.Context, recipientAccountID uuid.UUID, amountPaise int64) (uuid.UUID, error) {
	txID := uuid.New()
	log.Printf("fake reward payout: recipient_account_id=%s amount_paise=%d tx_id=%s", recipientAccountID, amountPaise, txID)
	return txID, nil
}
