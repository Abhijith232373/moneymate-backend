package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/domain"
	"github.com/moneymate-2026/moneymate-backend/services/notification/sqlc/generated"
)

type PreferenceRepo struct {
	q *generated.Queries
}

func NewPreferenceRepo(pool *pgxpool.Pool) *PreferenceRepo {
	return &PreferenceRepo{q: generated.New(pool)}
}

func (r *PreferenceRepo) Get(ctx context.Context, recipientType domain.RecipientType, recipientID uuid.UUID, category domain.Category) (bool, error) {
	enabled, err := r.q.GetPreference(ctx, generated.GetPreferenceParams{
		RecipientType: string(recipientType),
		RecipientID:   recipientID,
		Category:      generated.NotificationCategory(category),
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return true, nil // default: enabled
	}
	if err != nil {
		return false, fmt.Errorf("get preference: %w", err)
	}
	return enabled, nil
}

func (r *PreferenceRepo) Upsert(ctx context.Context, recipientType domain.RecipientType, recipientID uuid.UUID, category domain.Category, enabled bool) (*domain.Preference, error) {
	row, err := r.q.UpsertPreference(ctx, generated.UpsertPreferenceParams{
		RecipientType: string(recipientType),
		RecipientID:   recipientID,
		Category:      generated.NotificationCategory(category),
		Enabled:       enabled,
	})
	if err != nil {
		return nil, fmt.Errorf("upsert preference: %w", err)
	}
	return &domain.Preference{
		ID:            row.ID,
		RecipientType: domain.RecipientType(row.RecipientType),
		RecipientID:   row.RecipientID,
		Category:      domain.Category(row.Category),
		Enabled:       row.Enabled,
		UpdatedAt:     row.UpdatedAt,
	}, nil
}

func (r *PreferenceRepo) List(ctx context.Context, recipientType domain.RecipientType, recipientID uuid.UUID) ([]*domain.Preference, error) {
	rows, err := r.q.ListPreferences(ctx, generated.ListPreferencesParams{
		RecipientType: string(recipientType),
		RecipientID:   recipientID,
	})
	if err != nil {
		return nil, fmt.Errorf("list preferences: %w", err)
	}
	out := make([]*domain.Preference, 0, len(rows))
	for _, row := range rows {
		out = append(out, &domain.Preference{
			ID:            row.ID,
			RecipientType: domain.RecipientType(row.RecipientType),
			RecipientID:   row.RecipientID,
			Category:      domain.Category(row.Category),
			Enabled:       row.Enabled,
			UpdatedAt:     row.UpdatedAt,
		})
	}
	return out, nil
}
