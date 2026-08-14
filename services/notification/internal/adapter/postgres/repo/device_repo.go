package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/domain"
	"github.com/moneymate-2026/moneymate-backend/services/notification/sqlc/generated"
)

type DeviceRepo struct {
	q *generated.Queries
}

func NewDeviceRepo(pool *pgxpool.Pool) *DeviceRepo {
	return &DeviceRepo{q: generated.New(pool)}
}

func (r *DeviceRepo) Upsert(ctx context.Context, t *domain.DeviceToken) (*domain.DeviceToken, error) {
	row, err := r.q.UpsertDeviceToken(ctx, generated.UpsertDeviceTokenParams{
		RecipientType: string(t.RecipientType),
		RecipientID:   t.RecipientID,
		DeviceID:      t.DeviceID,
		Token:         t.Token,
		Platform:      t.Platform,
		AppVersion:    strPtr(t.AppVersion),
	})
	if err != nil {
		return nil, fmt.Errorf("upsert device token: %w", err)
	}
	return &domain.DeviceToken{
		ID:            row.ID,
		RecipientType: domain.RecipientType(row.RecipientType),
		RecipientID:   row.RecipientID,
		DeviceID:      row.DeviceID,
		Token:         row.Token,
		Platform:      row.Platform,
		AppVersion:    derefStr(row.AppVersion),
		IsActive:      row.IsActive,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}, nil
}

func (r *DeviceRepo) ListActiveByRecipient(ctx context.Context, recipientType domain.RecipientType, recipientID uuid.UUID) ([]*domain.DeviceToken, error) {
	rows, err := r.q.ListActiveTokensByRecipient(ctx, generated.ListActiveTokensByRecipientParams{
		RecipientType: string(recipientType),
		RecipientID:   recipientID,
	})
	if err != nil {
		return nil, fmt.Errorf("list active tokens: %w", err)
	}
	out := make([]*domain.DeviceToken, 0, len(rows))
	for _, row := range rows {
		out = append(out, &domain.DeviceToken{
			ID:            row.ID,
			RecipientType: domain.RecipientType(row.RecipientType),
			RecipientID:   row.RecipientID,
			DeviceID:      row.DeviceID,
			Token:         row.Token,
			Platform:      row.Platform,
			AppVersion:    derefStr(row.AppVersion),
			IsActive:      row.IsActive,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
		})
	}
	return out, nil
}

func (r *DeviceRepo) Deactivate(ctx context.Context, id uuid.UUID) error {
	return r.q.DeactivateDeviceToken(ctx, id)
}

func (r *DeviceRepo) DeactivateByDevice(ctx context.Context, recipientType domain.RecipientType, recipientID uuid.UUID, deviceID string) error {
	return r.q.DeactivateTokensByDevice(ctx, generated.DeactivateTokensByDeviceParams{
		RecipientType: string(recipientType),
		RecipientID:   recipientID,
		DeviceID:      deviceID,
	})
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
