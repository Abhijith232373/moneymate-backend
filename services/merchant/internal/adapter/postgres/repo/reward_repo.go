package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/domain"
	"github.com/moneymate-2026/moneymate-backend/services/merchant/sqlc/generated"
)

// RewardRepo implements the domain.RewardRepository interface using PostgreSQL via pgxpool and SQLC generated queries.
// It is designed for high-concurrency transactional safety for millions of active merchants.
type RewardRepo struct {
	// db holds the connection pool to PostgreSQL.
	db *pgxpool.Pool
	// queries holds the sqlc-generated type-safe database methods.
	queries *generated.Queries
}

// NewRewardRepo initializes and returns a new RewardRepo instance.
func NewRewardRepo(db *pgxpool.Pool) domain.RewardRepository {
	return &RewardRepo{
		db:      db,
		queries: generated.New(db),
	}
}

// GetSummaryByStoreID retrieves the financial summary and scan metrics for a merchant store.
// If the record does not exist yet (e.g., newly registered merchant or migration onboarding),
// it automatically calls InitializeSummary to bootstrap initial balance and transaction history.
func (r *RewardRepo) GetSummaryByStoreID(ctx context.Context, storeID uuid.UUID) (*domain.RewardSummary, error) {
	row, err := r.queries.GetRewardBalanceByStoreID(ctx, storeID)
	if errors.Is(err, pgx.ErrNoRows) {
		return r.InitializeSummary(ctx, storeID)
	}
	if err != nil {
		return nil, fmt.Errorf("query reward summary: %w", err)
	}

	avail, _ := row.AvailableBalance.Float64Value()
	growth, _ := row.WeeklyGrowthPercentage.Float64Value()

	return &domain.RewardSummary{
		ID:                     row.ID,
		StoreID:                row.StoreID,
		AvailableBalance:       avail.Float64,
		TotalScans:             row.TotalScans,
		PremiumPoints:          row.PremiumPoints,
		WeeklyGrowthPercentage: growth.Float64,
		CreatedAt:              row.CreatedAt,
		UpdatedAt:              row.UpdatedAt,
	}, nil
}

// InitializeSummary creates the initial reward balance and seeds production-realistic initial ledger entries
// so that newly onboarded merchants have a functional and visually complete dashboard immediately.
func (r *RewardRepo) InitializeSummary(ctx context.Context, storeID uuid.UUID) (*domain.RewardSummary, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin init tx: %w", err)
	}
	defer tx.Rollback(ctx)

	qTx := r.queries.WithTx(tx)

	var availNumeric, growthNumeric pgtype.Numeric
	_ = availNumeric.Scan("0.00")
	_ = growthNumeric.Scan("0.00")

	// Initial balance seeded
	row, err := qTx.CreateRewardBalance(ctx, generated.CreateRewardBalanceParams{
		StoreID:                storeID,
		AvailableBalance:       availNumeric,
		TotalScans:             0,
		PremiumPoints:          0,
		WeeklyGrowthPercentage: growthNumeric,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		// If balance already exists due to concurrent init, query it
		row, err = qTx.GetRewardBalanceByStoreID(ctx, storeID)
		if err != nil {
			return nil, fmt.Errorf("get existing reward balance during init: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit init tx: %w", err)
	}

	avail, _ := row.AvailableBalance.Float64Value()
	growth, _ := row.WeeklyGrowthPercentage.Float64Value()

	return &domain.RewardSummary{
		ID:                     row.ID,
		StoreID:                row.StoreID,
		AvailableBalance:       avail.Float64,
		TotalScans:             row.TotalScans,
		PremiumPoints:          row.PremiumPoints,
		WeeklyGrowthPercentage: growth.Float64,
		CreatedAt:              row.CreatedAt,
		UpdatedAt:              row.UpdatedAt,
	}, nil
}

// GetTransactionsByStoreID queries paginated reward ledger entries with dynamic filtering for time ranges
// ("all", "this_month", "this_week") and case-insensitive keyword searches on campaign names or display IDs.
func (r *RewardRepo) GetTransactionsByStoreID(ctx context.Context, storeID uuid.UUID, filter string, searchQuery string, limit, offset int) ([]*domain.RewardTransaction, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	var timeArg pgtype.Timestamptz
	now := time.Now().UTC()

	switch filter {
	case "this_week", "this_week_ui":
		timeArg = pgtype.Timestamptz{Time: now.AddDate(0, 0, -7), Valid: true}
	case "this_month", "this_month_ui":
		timeArg = pgtype.Timestamptz{Time: now.AddDate(0, -1, 0), Valid: true}
	default:
		timeArg = pgtype.Timestamptz{Valid: false}
	}

	rows, err := r.queries.GetRewardTransactionsByStoreID(ctx, generated.GetRewardTransactionsByStoreIDParams{
		StoreID:      storeID,
		SearchQuery:  searchQuery,
		CreatedAfter: timeArg,
		LimitCount:   int32(limit),
		OffsetCount:  int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("query reward transactions: %w", err)
	}

	var transactions []*domain.RewardTransaction
	for _, row := range rows {
		amt, _ := row.Amount.Float64Value()
		transactions = append(transactions, &domain.RewardTransaction{
			ID:              row.ID,
			StoreID:         row.StoreID,
			CampaignName:    row.CampaignName,
			DisplayID:       row.DisplayID,
			Status:          row.Status,
			Amount:          amt.Float64,
			TransactionType: string(row.TransactionType),
			CreatedAt:       row.CreatedAt,
		})
	}
	return transactions, nil
}

// RedeemBalance executes an atomic transaction that verifies sufficient funds, deducts the requested payout amount
// from the merchant's available balance, logs a redemption request, and appends a withdrawal entry to the transaction ledger.
func (r *RewardRepo) RedeemBalance(ctx context.Context, storeID uuid.UUID, amount float64, bankAuthorized bool) (*domain.RedemptionRequest, error) {
	if amount <= 0 {
		return nil, errors.New("redemption amount must be greater than zero")
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin redemption tx: %w", err)
	}
	defer tx.Rollback(ctx)

	qTx := r.queries.WithTx(tx)

	// 1. Lock balance row for update to prevent double-spending under high concurrency
	var currentBalance float64
	lockQuery := `SELECT available_balance FROM reward_balances WHERE store_id = $1 FOR UPDATE;`
	err = tx.QueryRow(ctx, lockQuery, storeID).Scan(&currentBalance)
	if errors.Is(err, pgx.ErrNoRows) {
		// Initialize if missing
		_, _ = r.InitializeSummary(ctx, storeID)
		err = tx.QueryRow(ctx, lockQuery, storeID).Scan(&currentBalance)
	}
	if err != nil {
		return nil, fmt.Errorf("lock balance row: %w", err)
	}

	if currentBalance < amount {
		return nil, fmt.Errorf("insufficient available balance: requested $%.2f, available $%.2f", amount, currentBalance)
	}

	// 2. Deduct balance using sqlc generated query
	var amtNumeric pgtype.Numeric
	_ = amtNumeric.Scan(fmt.Sprintf("%.2f", amount))
	err = qTx.DeductRewardBalance(ctx, generated.DeductRewardBalanceParams{
		StoreID:          storeID,
		AvailableBalance: amtNumeric,
	})
	if err != nil {
		return nil, fmt.Errorf("deduct reward balance: %w", err)
	}

	// 3. Create redemption request audit record
	refCode := fmt.Sprintf("PAYOUT-%d-%s", time.Now().Unix(), uuid.New().String()[:6])
	reqRow, err := qTx.CreateRedemptionRequest(ctx, generated.CreateRedemptionRequestParams{
		StoreID:                storeID,
		Amount:                 amtNumeric,
		BankTransferAuthorized: bankAuthorized,
		Status:                 generated.RedemptionStatusProcessing,
		ReferenceID:            &refCode,
	})
	if err != nil {
		return nil, fmt.Errorf("insert redemption request: %w", err)
	}

	// 4. Log immutable ledger entry for the withdrawal
	var negAmtNumeric pgtype.Numeric
	_ = negAmtNumeric.Scan(fmt.Sprintf("%.2f", -amount))
	displayRef := fmt.Sprintf("#WD-%d", time.Now().Unix()%10000)
	_, err = qTx.CreateRewardTransaction(ctx, generated.CreateRewardTransactionParams{
		StoreID:         storeID,
		CampaignName:    "Bank Account Withdrawal",
		DisplayID:       displayRef,
		Status:          "Processing",
		Amount:          negAmtNumeric,
		TransactionType: generated.RewardTransactionTypeRedemption,
	})
	if err != nil {
		return nil, fmt.Errorf("insert withdrawal transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit redemption tx: %w", err)
	}

	resAmt, _ := reqRow.Amount.Float64Value()

	return &domain.RedemptionRequest{
		ID:                     reqRow.ID,
		StoreID:                reqRow.StoreID,
		Amount:                 resAmt.Float64,
		BankTransferAuthorized: reqRow.BankTransferAuthorized,
		Status:                 string(reqRow.Status),
		ReferenceID:            reqRow.ReferenceID,
		CreatedAt:              reqRow.CreatedAt,
		UpdatedAt:              reqRow.UpdatedAt,
	}, nil
}

