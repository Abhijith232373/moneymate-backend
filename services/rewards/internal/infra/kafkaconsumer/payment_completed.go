package kafkaconsumer

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

type PaymentCompletedEvent struct {
	EventID            string    `json:"event_id"`
	EventType          string    `json:"event_type"`
	TransactionID      uuid.UUID `json:"transaction_id"`
	RecipientID        uuid.UUID `json:"recipient_id"`
	RecipientAccountID uuid.UUID `json:"recipient_account_id"`
	RecipientType      string    `json:"recipient_type"`
	AmountPaise        int64     `json:"amount_paise"`
	OccurredAt         time.Time `json:"occurred_at"`
}

func ParsePaymentCompletedEvent(payload []byte) (*PaymentCompletedEvent, error) {
	var event PaymentCompletedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, err
	}
	log.Printf("parsed payment completed event: id=%s txn=%s", event.EventID, event.TransactionID)
	return &event, nil
}
