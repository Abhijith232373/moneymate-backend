package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Preference struct {
	ID            uuid.UUID
	RecipientType RecipientType
	RecipientID   uuid.UUID
	Category      Category
	Enabled       bool
	UpdatedAt     time.Time
}

type PreferenceRepository interface {
	// Get returns enabled=true when no explicit row exists (opt-out by default off).
	Get(ctx context.Context, recipientType RecipientType, recipientID uuid.UUID, category Category) (bool, error)
	Upsert(ctx context.Context, recipientType RecipientType, recipientID uuid.UUID, category Category, enabled bool) (*Preference, error)
	List(ctx context.Context, recipientType RecipientType, recipientID uuid.UUID) ([]*Preference, error)
}
