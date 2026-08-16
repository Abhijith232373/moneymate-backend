package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moneymate-2026/moneymate-backend/services/payment/internal/domain"
	"github.com/moneymate-2026/moneymate-backend/services/payment/sqlc/generated"
)

type CategoryRepo struct {
	q *generated.Queries
}

func NewCategoryRepo(pool *pgxpool.Pool) *CategoryRepo {
	return &CategoryRepo{q: generated.New(pool)}
}

func (r *CategoryRepo) Create(ctx context.Context, userID uuid.UUID, name string) (*domain.Category, error) {
	row, err := r.q.CreateCategory(ctx, generated.CreateCategoryParams{UserID: userID, Name: name})
	if err != nil {
		return nil, mapDBErr(err)
	}
	return toDomainCategory(row), nil
}

func (r *CategoryRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Category, error) {
	rows, err := r.q.ListCategoriesByUser(ctx, userID)
	if err != nil {
		return nil, mapDBErr(err)
	}
	out := make([]domain.Category, len(rows))
	for i, row := range rows {
		out[i] = *toDomainCategory(row)
	}
	return out, nil
}

func (r *CategoryRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	row, err := r.q.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, mapDBErr(err)
	}
	return toDomainCategory(row), nil
}

func (r *CategoryRepo) Update(ctx context.Context, id, userID uuid.UUID, name string) (*domain.Category, error) {
	row, err := r.q.UpdateCategory(ctx, generated.UpdateCategoryParams{ID: id, Name: name, UserID: userID})
	if err != nil {
		return nil, mapDBErr(err)
	}
	return toDomainCategory(row), nil
}

func (r *CategoryRepo) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return mapDBErr(r.q.DeleteCategory(ctx, generated.DeleteCategoryParams{ID: id, UserID: userID}))
}

func toDomainCategory(row generated.PaymentCategory) *domain.Category {
	return &domain.Category{
		ID: row.ID, UserID: row.UserID, Name: row.Name,
		CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}
}