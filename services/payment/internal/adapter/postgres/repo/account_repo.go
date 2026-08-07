package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moneymate-2026/moneymate-backend/services/payment/internal/domain"
	"github.com/moneymate-2026/moneymate-backend/services/payment/sqlc/generated"
	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"
)

type AccountRepo struct {
	pool *pgxpool.Pool
	q    *generated.Queries
}

func NewAccountRepo(pool *pgxpool.Pool) *AccountRepo {
	return &AccountRepo{pool: pool, q: generated.New(pool)}
}

func (r *AccountRepo) Create(ctx context.Context, a *domain.Account) (*domain.Account, error) {
	row, err := r.q.CreateAccount(ctx, generated.CreateAccountParams{
		UserID:     uuidPtrToPgtype(a.UserID),
		MerchantID: uuidPtrToPgtype(a.MerchantID),
		Column3:    generated.PaymentAccountType(a.Type),
		Currency:   a.Currency,
	})
	if err != nil {
		return nil, mapDBErr(err)
	}
	return rowToAccount(generated.GetAccountByIDRow(row)), nil
}

func (r *AccountRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	row, err := r.q.GetAccountByID(ctx, id)
	if err != nil {
		return nil, mapDBErr(err)
	}
	return rowToAccount(row), nil
}

func (r *AccountRepo) GetWalletByUserID(ctx context.Context, userID uuid.UUID) (*domain.Account, error) {
	row, err := r.q.GetWalletByUserID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, mapDBErr(err)
	}
	return rowToAccount(generated.GetAccountByIDRow(row)), nil
}

func (r *AccountRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error) {
	rows, err := r.q.ListAccountsByUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, mapDBErr(err)
	}
	accounts := make([]*domain.Account, 0, len(rows))
	for _, row := range rows {
		accounts = append(accounts, rowToAccount(generated.GetAccountByIDRow(row)))
	}
	return accounts, nil
}

func rowToAccount(row generated.GetAccountByIDRow) *domain.Account {
	var userID, merchantID *uuid.UUID
	if row.UserID.Valid {
		id := uuid.UUID(row.UserID.Bytes)
		userID = &id
	}
	if row.MerchantID.Valid {
		id := uuid.UUID(row.MerchantID.Bytes)
		merchantID = &id
	}
	return &domain.Account{
		ID:         row.ID,
		UserID:     userID,
		MerchantID: merchantID,
		Type:       domain.AccountType(row.Type),
		Currency:   row.Currency,
		Balance:    row.Balance,
		Version:    row.Version,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}
}

func mapDBErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, apperrors.ErrNotFound) {
		return err
	}
	mapped := apperrors.MapDBErrors(err)
	if mapped != nil {
		return mapped
	}
	return fmt.Errorf("db: %w", err)
}
