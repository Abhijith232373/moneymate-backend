package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/domain"
	"github.com/moneymate-2026/moneymate-backend/services/merchant/sqlc/generated"
)

type EarningsRepo struct {
	db      *pgxpool.Pool
	queries *generated.Queries
}

func NewEarningsRepo(db *pgxpool.Pool) domain.EarningsRepository {
	return &EarningsRepo{
		db:      db,
		queries: generated.New(db),
	}
}

func (r *EarningsRepo) GetEarningsStats(ctx context.Context, storeID uuid.UUID) (*domain.EarningsStats, error) {
	row, err := r.queries.GetEarningsStats(ctx, storeID)
	if errors.Is(err, pgx.ErrNoRows) {
		// Auto-initialize if it doesn't exist
		return r.UpsertEarningsStats(ctx, storeID, 0, 0)
	}
	if err != nil {
		return nil, fmt.Errorf("get earnings stats: %w", err)
	}

	totalEarn, _ := row.TotalEarned.Float64Value()

	return &domain.EarningsStats{
		StoreID:     row.StoreID,
		TotalScans:  row.TotalScans,
		TotalEarned: totalEarn.Float64,
		UpdatedAt:   row.UpdatedAt,
	}, nil
}

func (r *EarningsRepo) UpsertEarningsStats(ctx context.Context, storeID uuid.UUID, totalScans int64, totalEarned float64) (*domain.EarningsStats, error) {
	var te pgtype.Numeric
	_ = te.Scan(fmt.Sprintf("%f", totalEarned))

	row, err := r.queries.UpsertEarningsStats(ctx, generated.UpsertEarningsStatsParams{
		StoreID:     storeID,
		TotalScans:  totalScans,
		TotalEarned: te,
	})
	if err != nil {
		return nil, fmt.Errorf("upsert earnings stats: %w", err)
	}

	earnVal, _ := row.TotalEarned.Float64Value()

	return &domain.EarningsStats{
		StoreID:     row.StoreID,
		TotalScans:  row.TotalScans,
		TotalEarned: earnVal.Float64,
		UpdatedAt:   row.UpdatedAt,
	}, nil
}

func (r *EarningsRepo) GetRequestedMilestones(ctx context.Context, storeID uuid.UUID) ([]int32, error) {
	rows, err := r.queries.GetRequestedMilestones(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("get requested milestones: %w", err)
	}
	return rows, nil
}

func (r *EarningsRepo) CreatePayoutRequest(ctx context.Context, storeID uuid.UUID, milestoneScans int32, rewardAmount float64) (*domain.EarningsPayoutRequest, error) {
	var amt pgtype.Numeric
	_ = amt.Scan(fmt.Sprintf("%f", rewardAmount))

	row, err := r.queries.CreatePayoutRequest(ctx, generated.CreatePayoutRequestParams{
		StoreID:        storeID,
		MilestoneScans: milestoneScans,
		RewardAmount:   amt,
	})
	if err != nil {
		return nil, fmt.Errorf("create payout request: %w", err)
	}

	amtVal, _ := row.RewardAmount.Float64Value()
	return &domain.EarningsPayoutRequest{
		ID:             row.ID,
		StoreID:        row.StoreID,
		MilestoneScans: row.MilestoneScans,
		RewardAmount:   amtVal.Float64,
		Status:         row.Status,
		CreatedAt:      row.CreatedAt,
	}, nil
}
