package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/abijith/moneymate-backend/services/support/internal/domain"
)

type ChatRepo struct {
	db      *sql.DB
	querier Querier
}

func NewChatRepo(db *sql.DB) *ChatRepo {
	return &ChatRepo{
		db:      db,
		querier: New(db),
	}
}

func (r *ChatRepo) CreateChatMessage(ctx context.Context, msg *domain.ChatMessage) (*domain.ChatMessage, error) {
	row, err := r.querier.CreateChatMessage(ctx, CreateChatMessageParams{
		SenderID:     msg.SenderID,
		SenderType:   msg.SenderType,
		ReceiverID:   msg.ReceiverID,
		ReceiverType: msg.ReceiverType,
		Message:      msg.Message,
	})
	if err != nil {
		return nil, err
	}

	return &domain.ChatMessage{
		ID:           row.ID,
		SenderID:     row.SenderID,
		SenderType:   row.SenderType,
		ReceiverID:   row.ReceiverID,
		ReceiverType: row.ReceiverType,
		Message:      row.Message,
		IsRead:       row.IsRead.Bool,
		CreatedAt:    row.CreatedAt.Time,
	}, nil
}

func (r *ChatRepo) GetChatHistory(ctx context.Context, userID, adminID uuid.UUID) ([]*domain.ChatMessage, error) {
	rows, err := r.querier.GetChatHistory(ctx, GetChatHistoryParams{
		SenderID:   userID,
		ReceiverID: adminID,
	})
	if err != nil {
		return nil, err
	}

	var res []*domain.ChatMessage
	for _, row := range rows {
		res = append(res, &domain.ChatMessage{
			ID:           row.ID,
			SenderID:     row.SenderID,
			SenderType:   row.SenderType,
			ReceiverID:   row.ReceiverID,
			ReceiverType: row.ReceiverType,
			Message:      row.Message,
			IsRead:       row.IsRead.Bool,
			CreatedAt:    row.CreatedAt.Time,
		})
	}
	return res, nil
}

func (r *ChatRepo) GetAdminChatHistory(ctx context.Context, adminID uuid.UUID) ([]*domain.ChatMessage, error) {
	rows, err := r.querier.GetAdminChatHistory(ctx, adminID)
	if err != nil {
		return nil, err
	}

	var res []*domain.ChatMessage
	for _, row := range rows {
		res = append(res, &domain.ChatMessage{
			ID:           row.ID,
			SenderID:     row.SenderID,
			SenderType:   row.SenderType,
			ReceiverID:   row.ReceiverID,
			ReceiverType: row.ReceiverType,
			Message:      row.Message,
			IsRead:       row.IsRead.Bool,
			CreatedAt:    row.CreatedAt.Time,
		})
	}
	return res, nil
}

func (r *ChatRepo) MarkMessagesAsRead(ctx context.Context, senderID, receiverID uuid.UUID) error {
	return r.querier.MarkMessagesAsRead(ctx, MarkMessagesAsReadParams{
		SenderID:   senderID,
		ReceiverID: receiverID,
	})
}
