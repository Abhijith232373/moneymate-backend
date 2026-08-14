package kafkaconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/domain"
	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/usecases"
)

// HandleEvent is the per-message callback wired into shared/pkg/kafka Consumer.Run.
func HandleEvent(ctx context.Context, uc *usecases.NotificationUsecase, payload []byte) error {
	var env domain.EventEnvelope
	if err := json.Unmarshal(payload, &env); err != nil {
		// Already-committed message — return an error just to log it; the
		// consumer moves on either way.
		return fmt.Errorf("unmarshal event envelope: %w", err)
	}
	if err := uc.HandleEvent(ctx, env); err != nil {
		log.Printf("notification: handle event %s failed: %v", env.EventType, err)
		return err
	}
	return nil
}
