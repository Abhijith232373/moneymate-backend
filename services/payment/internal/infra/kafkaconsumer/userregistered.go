package kafkaconsumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/moneymate-2026/moneymate-backend/services/payment/internal/usecases"
)

type UserRegisteredEvent struct {
	UserID string `json:"user_id"`
	Handle string `json:"handle"`
}

func HandleUserRegistered(ctx context.Context, wallets usecases.WalletUsecase, payload []byte) error {
	var evt UserRegisteredEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		return err
	}
	if _, err := wallets.CreateWallet(ctx, evt.UserID, evt.Handle); err != nil {
		log.Printf("create wallet for user %s failed: %v", evt.UserID, err)
		return err
	}
	return nil
}