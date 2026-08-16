package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/domain"
	"github.com/moneymate-2026/moneymate-backend/services/notification/sqlc/generated"
)

type InboxRepo struct {
	q *generated.Queries
}

func NewInboxRepo(pool *pgxpool.Pool) *InboxRepo {
	return &InboxRepo{q: generated.New(pool)}
}

func (r *InboxRepo) Insert(ctx context.Context, m *domain.InboxMessage) (uuid.UUID, error) {
	data, err := json.Marshal(m.Data)
	if err != nil {
		return uuid.Nil, fmt.Errorf("marshal inbox data: %w", err)
	}
	id, err := r.q.InsertInbox(ctx, generated.InsertInboxParams{
		RecipientType: string(m.RecipientType),
		RecipientID:   m.RecipientID,
		Category:      generated.NotificationCategory(m.Category),
		Title:         m.Title,
		Body:          m.Body,
		Data:          data,
		EventID:       m.EventID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, nil // duplicate event — already processed
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("insert inbox: %w", err)
	}
	return id, nil
}

func (r *InboxRepo) List(ctx context.Context, recipientType domain.RecipientType, recipientID uuid.UUID, limit, offset int32) ([]*domain.InboxMessage, error) {
	rows, err := r.q.ListInbox(ctx, generated.ListInboxParams{
		RecipientType: string(recipientType),
		RecipientID:   recipientID,
		Limit:         limit,
		Offset:        offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list inbox: %w", err)
	}
	out := make([]*domain.InboxMessage, 0, len(rows))
	for _, row := range rows {
		out = append(out, rowToInbox(row))
	}
	return out, nil
}

func (r *InboxRepo) Get(ctx context.Context, id uuid.UUID, recipientType domain.RecipientType, recipientID uuid.UUID) (*domain.InboxMessage, error) {
	row, err := r.q.GetInbox(ctx, generated.GetInboxParams{
		ID:            id,
		RecipientType: string(recipientType),
		RecipientID:   recipientID,
	})
	if err != nil {
		return nil, err
	}
	return rowToInbox(row), nil
}

func (r *InboxRepo) MarkRead(ctx context.Context, id uuid.UUID, recipientType domain.RecipientType, recipientID uuid.UUID) error {
	return r.q.MarkInboxRead(ctx, generated.MarkInboxReadParams{
		ID:            id,
		RecipientType: string(recipientType),
		RecipientID:   recipientID,
	})
}

func (r *InboxRepo) MarkSent(ctx context.Context, id uuid.UUID) error {
	return r.q.MarkInboxSent(ctx, id)
}

func rowToInbox(row generated.NotificationInbox) *domain.InboxMessage {
	var data map[string]any
	_ = json.Unmarshal(row.Data, &data)
	return &domain.InboxMessage{
		ID:            row.ID,
		RecipientType: domain.RecipientType(row.RecipientType),
		RecipientID:   row.RecipientID,
		Category:      domain.Category(row.Category),
		Title:         row.Title,
		Body:          row.Body,
		Data:          data,
		EventID:       row.EventID,
		Status:        row.Status,
		ReadAt:        timeFromPgtype(row.ReadAt),
		SentAt:        timeFromPgtype(row.SentAt),
		CreatedAt:     row.CreatedAt,
	}
}
