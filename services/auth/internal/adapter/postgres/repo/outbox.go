// auth/internal/adapter/postgres/repo/outbox_repo.go
package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moneymate-2026/moneymate-backend/auth/internal/domain"
	db "github.com/moneymate-2026/moneymate-backend/auth/sqlc/generated"
	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"
	"github.com/moneymate-2026/moneymate-backend/shared/pkg/pgxtx"
)

type outboxRepo struct{ q *db.Queries }

func NewOutboxRepo(pool *pgxpool.Pool) domain.OutboxRepository {
	return &outboxRepo{q: db.New(pool)}
}

func (r *outboxRepo) queries(ctx context.Context) *db.Queries {
	if tx, ok := pgxtx.FromContext(ctx); ok {
		return r.q.WithTx(tx)
	}
	return r.q
}

func (r *outboxRepo) Insert(ctx context.Context, e *domain.OutboxEvent) error {
	err := r.queries(ctx).InsertOutboxEvent(ctx, db.InsertOutboxEventParams{
		ID:      uuidToPgtype(e.ID),
		Topic:   e.Topic,
		Payload: e.Payload,
	})
	return apperrors.MapDBErrors(err)
}

func (r *outboxRepo) FetchUnpublished(ctx context.Context, limit int32) ([]*domain.OutboxEvent, error) {
	rows, err := r.q.FetchUnpublishedOutboxEvents(ctx, limit)
	if err != nil {
		return nil, apperrors.MapDBErrors(err)
	}
	events := make([]*domain.OutboxEvent, len(rows))
	for i, row := range rows {
		events[i] = &domain.OutboxEvent{
			ID:        uuid.UUID(row.ID.Bytes),
			Topic:     row.Topic,
			Payload:   row.Payload,
			CreatedAt: row.CreatedAt.Time,
		}
	}
	return events, nil
}

func (r *outboxRepo) MarkPublished(ctx context.Context, id uuid.UUID) error {
	return apperrors.MapDBErrors(r.q.MarkOutboxEventPublished(ctx, uuidToPgtype(id)))
}