package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/moneymate-2026/moneymate-backend/services/rewards/config"
	"github.com/moneymate-2026/moneymate-backend/shared/pkg/kafka"
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

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	producer, err := kafka.NewProducer(kafka.Config{
		Brokers:  cfg.Kafka.Brokers,
		Username: cfg.Kafka.Username,
		Password: cfg.Kafka.Password,
		CACert:   cfg.Kafka.CACert,
	})
	if err != nil {
		log.Fatalf("failed to create kafka producer: %v", err)
	}
	defer producer.Close()

	event := PaymentCompletedEvent{
		EventID:            "local-test-001",
		EventType:          "moneymate.payment.completed",
		TransactionID:      uuid.New(),
		RecipientID:        uuid.New(),
		RecipientAccountID: uuid.New(),
		RecipientType:      "user",
		AmountPaise:        250000,
		OccurredAt:         time.Now(),
	}

	payload, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("failed to marshal event: %v", err)
	}

	err = producer.Publish(context.Background(), cfg.Rewards.PaymentCompletedTopic, []byte(event.RecipientID.String()), payload)
	if err != nil {
		log.Fatalf("failed to publish event: %v", err)
	}

	log.Printf("published fake payment completed event: %s", event.EventID)
}
