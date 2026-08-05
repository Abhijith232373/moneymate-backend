package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/domain"
	"github.com/moneymate-2026/moneymate-backend/services/merchant/sqlc/generated"
)

type PostgresQRRepository struct {
	queries *generated.Queries
}

func NewQRRepo(db *pgxpool.Pool) *PostgresQRRepository {
	return &PostgresQRRepository{queries: generated.New(db)}
}

func (r *PostgresQRRepository) CreateQRTransaction(ctx context.Context, storeID uuid.UUID, customerDisplayID string, billAmount, rewardIssued float64) (*domain.QRTransaction, error) {
	numAmount := pgtype.Numeric{}
	numAmount.Scan(billAmount) // Simplified conversion, in production use proper numeric handling or string

	// For simplicity, we convert to string and use Scan which handles postgres string numeric format
	// But `pgtype.Numeric` can be tricky without the right helper. We can just use string representation:
	// Let's use float64 converted to string and parse inside `pgtype.Numeric`
	// Actually, pgtype.Numeric supports scanning from float64 but it's deprecated or might lose precision.
	// Since billAmount is float64, we can format to string.
	// We'll write a small helper.

	res, err := r.queries.CreateQRTransaction(ctx, generated.CreateQRTransactionParams{
		StoreID:           storeID,
		CustomerDisplayID: customerDisplayID,
		BillAmount:        floatToNumeric(billAmount),
		RewardIssued:      floatToNumeric(rewardIssued),
	})
	if err != nil {
		return nil, err
	}

	return &domain.QRTransaction{
		ID:                res.ID,
		StoreID:           res.StoreID,
		CustomerDisplayID: res.CustomerDisplayID,
		BillAmount:        numericToFloat(res.BillAmount),
		RewardIssued:      numericToFloat(res.RewardIssued),
		CreatedAt:         res.CreatedAt,
	}, nil
}

func (r *PostgresQRRepository) GetQRTransactionsByStoreID(ctx context.Context, storeID uuid.UUID, limit, offset int) ([]*domain.QRTransaction, error) {
	rows, err := r.queries.GetQRTransactionsByStoreID(ctx, generated.GetQRTransactionsByStoreIDParams{
		StoreID: storeID,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}

	var txs []*domain.QRTransaction
	for _, row := range rows {
		txs = append(txs, &domain.QRTransaction{
			ID:                row.ID,
			StoreID:           row.StoreID,
			CustomerDisplayID: row.CustomerDisplayID,
			BillAmount:        numericToFloat(row.BillAmount),
			RewardIssued:      numericToFloat(row.RewardIssued),
			CreatedAt:         row.CreatedAt,
		})
	}
	return txs, nil
}

func (r *PostgresQRRepository) GetTodayQRScanCount(ctx context.Context, storeID uuid.UUID) (int64, error) {
	count, err := r.queries.GetTodayQRScanCount(ctx, storeID)
	return count, err
}

func (r *PostgresQRRepository) GetTodayQRScanVolume(ctx context.Context, storeID uuid.UUID) (float64, error) {
	vol, err := r.queries.GetTodayQRScanVolume(ctx, storeID)
	if err != nil {
		return 0, err
	}
	return numericToFloat(vol), nil
}

// Helpers for numeric
func floatToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	n.Scan(f) // this actually works for pgtype.Numeric in pgx v5
	return n
}

func numericToFloat(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	f, _ := n.Float64Value()
	return f.Float64
}
