package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type InboxMessage struct {
	ID            uuid.UUID
	RecipientType RecipientType
	RecipientID   uuid.UUID
	Category      Category
	Title         string
	Body          string
	Data          map[string]any
	EventID       string
	Status        string // pending | sent | failed
	ReadAt        *time.Time
	SentAt        *time.Time
	CreatedAt     time.Time
}

type InboxRepository interface {
	// Insert returns uuid.Nil (no error) when the event is a duplicate.
	Insert(ctx context.Context, m *InboxMessage) (uuid.UUID, error)
	List(ctx context.Context, recipientType RecipientType, recipientID uuid.UUID, limit, offset int32) ([]*InboxMessage, error)
	Get(ctx context.Context, id uuid.UUID, recipientType RecipientType, recipientID uuid.UUID) (*InboxMessage, error)
	MarkRead(ctx context.Context, id uuid.UUID, recipientType RecipientType, recipientID uuid.UUID) error
	MarkSent(ctx context.Context, id uuid.UUID) error
}
