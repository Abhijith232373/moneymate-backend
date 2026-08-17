package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type DeliveryLog struct {
	ID                uuid.UUID
	InboxID           uuid.UUID
	DeviceTokenID     *uuid.UUID
	Provider          string
	ProviderMessageID string
	Status            string // sent | failed | dropped
	ErrorCode         string
	AttemptCount      int32
	LastAttemptAt     *time.Time
	CreatedAt         time.Time
}

type DeliveryRepository interface {
	Insert(ctx context.Context, l *DeliveryLog) error
}
