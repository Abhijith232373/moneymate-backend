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

type WalletRepo struct {
	db      *pgxpool.Pool
	queries *generated.Queries
}

func NewWalletRepo(db *pgxpool.Pool) domain.WalletRepository {
	return &WalletRepo{
		db:      db,
		queries: generated.New(db),
	}
}

func (r *WalletRepo) GetWalletByStoreID(ctx context.Context, storeID uuid.UUID) (*domain.Wallet, error) {
	row, err := r.queries.GetWalletByStoreID(ctx, storeID)
	if errors.Is(err, pgx.ErrNoRows) {
		// Auto-initialize wallet if it doesn't exist
		return r.UpsertWallet(ctx, storeID, 0, 0, 0)
	}
	if err != nil {
		return nil, fmt.Errorf("get wallet by store id: %w", err)
	}

	availBal, _ := row.AvailableBalance.Float64Value()
	totalEarn, _ := row.TotalEarnings.Float64Value()
	totalRed, _ := row.TotalRedeemed.Float64Value()

	return &domain.Wallet{
		ID:               row.ID,
		StoreID:          row.StoreID,
		AvailableBalance: availBal.Float64,
		TotalEarnings:    totalEarn.Float64,
		TotalRedeemed:    totalRed.Float64,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}, nil
}

func (r *WalletRepo) UpsertWallet(ctx context.Context, storeID uuid.UUID, availableBalance, totalEarnings, totalRedeemed float64) (*domain.Wallet, error) {
	var ab, te, tr pgtype.Numeric
	_ = ab.Scan(fmt.Sprintf("%f", availableBalance))
	_ = te.Scan(fmt.Sprintf("%f", totalEarnings))
	_ = tr.Scan(fmt.Sprintf("%f", totalRedeemed))

	row, err := r.queries.UpsertWallet(ctx, generated.UpsertWalletParams{
		StoreID:          storeID,
		AvailableBalance: ab,
		TotalEarnings:    te,
		TotalRedeemed:    tr,
	})
	if err != nil {
		return nil, fmt.Errorf("upsert wallet: %w", err)
	}

	availBal, _ := row.AvailableBalance.Float64Value()
	totalEarn, _ := row.TotalEarnings.Float64Value()
	totalRed, _ := row.TotalRedeemed.Float64Value()

	return &domain.Wallet{
		ID:               row.ID,
		StoreID:          row.StoreID,
		AvailableBalance: availBal.Float64,
		TotalEarnings:    totalEarn.Float64,
		TotalRedeemed:    totalRed.Float64,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}, nil
}

func (r *WalletRepo) GetWalletTransactions(ctx context.Context, storeID uuid.UUID) ([]*domain.WalletTransaction, error) {
	rows, err := r.queries.GetWalletTransactions(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("get wallet transactions: %w", err)
	}

	var txns []*domain.WalletTransaction
	for _, row := range rows {
		amt, _ := row.Amount.Float64Value()
		txns = append(txns, &domain.WalletTransaction{
			ID:              row.ID,
			StoreID:         row.StoreID,
			TransactionID:   row.TransactionID,
			Title:           row.Title,
			Subtitle:        row.Subtitle,
			Amount:          amt.Float64,
			TransactionType: string(row.TxnType),
			CreatedAt:       row.CreatedAt,
		})
	}
	return txns, nil
}

func (r *WalletRepo) GetWalletTransactionsByType(ctx context.Context, storeID uuid.UUID, txnType string) ([]*domain.WalletTransaction, error) {
	rows, err := r.queries.GetWalletTransactionsByType(ctx, generated.GetWalletTransactionsByTypeParams{
		StoreID: storeID,
		TxnType: generated.WalletTxnType(txnType),
	})
	if err != nil {
		return nil, fmt.Errorf("get wallet transactions by type: %w", err)
	}

	var txns []*domain.WalletTransaction
	for _, row := range rows {
		amt, _ := row.Amount.Float64Value()
		txns = append(txns, &domain.WalletTransaction{
			ID:              row.ID,
			StoreID:         row.StoreID,
			TransactionID:   row.TransactionID,
			Title:           row.Title,
			Subtitle:        row.Subtitle,
			Amount:          amt.Float64,
			TransactionType: string(row.TxnType),
			CreatedAt:       row.CreatedAt,
		})
	}
	return txns, nil
}

func (r *WalletRepo) CreateWalletTransaction(ctx context.Context, storeID uuid.UUID, transactionID, title, subtitle, txnType string, amount float64) (*domain.WalletTransaction, error) {
	var amt pgtype.Numeric
	_ = amt.Scan(fmt.Sprintf("%f", amount))

	row, err := r.queries.CreateWalletTransaction(ctx, generated.CreateWalletTransactionParams{
		StoreID:       storeID,
		TransactionID: transactionID,
		Title:         title,
		Subtitle:      subtitle,
		Amount:        amt,
		TxnType:       generated.WalletTxnType(txnType),
	})
	if err != nil {
		return nil, fmt.Errorf("create wallet transaction: %w", err)
	}

	amtVal, _ := row.Amount.Float64Value()
	return &domain.WalletTransaction{
		ID:              row.ID,
		StoreID:         row.StoreID,
		TransactionID:   row.TransactionID,
		Title:           row.Title,
		Subtitle:        row.Subtitle,
		Amount:          amtVal.Float64,
		TransactionType: string(row.TxnType),
		CreatedAt:       row.CreatedAt,
	}, nil
}
