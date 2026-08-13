package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ChatMessage struct {
	ID           uuid.UUID `json:"id"`
	SenderID     uuid.UUID `json:"sender_id"`
	SenderType   string    `json:"sender_type"`
	ReceiverID   uuid.UUID `json:"receiver_id"`
	ReceiverType string    `json:"receiver_type"`
	Message      string    `json:"message"`
	IsRead       bool      `json:"is_read"`
	CreatedAt    time.Time `json:"created_at"`
}

type ChatRepository interface {
	CreateChatMessage(ctx context.Context, msg *ChatMessage) (*ChatMessage, error)
	GetChatHistory(ctx context.Context, userID, adminID uuid.UUID) ([]*ChatMessage, error)
	GetAdminChatHistory(ctx context.Context, adminID uuid.UUID) ([]*ChatMessage, error)
	MarkMessagesAsRead(ctx context.Context, senderID, receiverID uuid.UUID) error
}
