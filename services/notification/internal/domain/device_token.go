package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type DeviceToken struct {
	ID            uuid.UUID
	RecipientType RecipientType
	RecipientID   uuid.UUID
	DeviceID      string
	Token         string
	Platform      string
	AppVersion    string
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type DeviceTokenRepository interface {
	Upsert(ctx context.Context, t *DeviceToken) (*DeviceToken, error)
	ListActiveByRecipient(ctx context.Context, recipientType RecipientType, recipientID uuid.UUID) ([]*DeviceToken, error)
	Deactivate(ctx context.Context, id uuid.UUID) error
	DeactivateByDevice(ctx context.Context, recipientType RecipientType, recipientID uuid.UUID, deviceID string) error
}
