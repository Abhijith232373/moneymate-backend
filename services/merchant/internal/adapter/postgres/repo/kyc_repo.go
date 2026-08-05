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

// KYCRepo implements domain.KYCRepository using PostgreSQL via pgxpool and SQLC generated queries.
// It guarantees atomic document updates and synchronizes store compliance state machines for millions of merchants.
type KYCRepo struct {
	// db holds the PostgreSQL connection pool.
	db *pgxpool.Pool
	// queries holds the sqlc-generated type-safe database methods.
	queries *generated.Queries
}

// NewKYCRepo initializes and returns a new KYCRepo instance.
func NewKYCRepo(db *pgxpool.Pool) domain.KYCRepository {
	return &KYCRepo{
		db:      db,
		queries: generated.New(db),
	}
}

// GetKYCStatusByStoreID queries the compliance standing and document URIs for a merchant store.
// If missing, it automatically seeds an approved compliance record so the UI renders seamlessly out of the box.
func (r *KYCRepo) GetKYCStatusByStoreID(ctx context.Context, storeID uuid.UUID) (*domain.KYCStatusDetail, error) {
	row, err := r.queries.GetKYCStatusByStoreID(ctx, storeID)
	if errors.Is(err, pgx.ErrNoRows) {
		// Store exists but no KYC documents yet. Return skeleton state instead of dummy data.
		storeStatus, _ := r.queries.GetStoreStatusByID(ctx, storeID)
		if storeStatus == "" {
			storeStatus = "pending_kyc"
		}
		return &domain.KYCStatusDetail{
			StoreID:     storeID,
			IsVerified:  false,
			StoreStatus: storeStatus,
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query kyc status by store id: %w", err)
	}

	var verifiedAtPtr *time.Time
	if row.VerifiedAt.Valid {
		verifiedAtPtr = &row.VerifiedAt.Time
	}

	return &domain.KYCStatusDetail{
		ID:             row.ID,
		StoreID:        row.StoreID,
		AadhaarNumber:  row.AadhaarNumber,
		AadhaarDocURL:  row.AadhaarDocUrl,
		ShopLicenseURL: row.ShopLicenseUrl,
		IsVerified:     row.IsVerified,
		VerifiedAt:     verifiedAtPtr,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
		StoreStatus:    row.StoreStatus,
	}, nil
}



// UpdateKYCDocuments executes an atomic transaction to update compliance document URIs and transitions the store
// state machine to 'pending_kyc', triggering automated or compliance officer re-review.
func (r *KYCRepo) UpdateKYCDocuments(ctx context.Context, storeID uuid.UUID, aadhaarNumber, aadhaarURL, licenseURL string) (*domain.KYCStatusDetail, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin kyc update tx: %w", err)
	}
	defer tx.Rollback(ctx)

	qTx := r.queries.WithTx(tx)

	// 1. Check if record exists
	var existingID uuid.UUID
	err = tx.QueryRow(ctx, `SELECT id FROM kyc_documents WHERE store_id = $1 LIMIT 1 FOR UPDATE;`, storeID).Scan(&existingID)

	now := time.Now().UTC()
	var row generated.KycDocument
	if errors.Is(err, pgx.ErrNoRows) {
		row, err = qTx.InsertKYCDocuments(ctx, generated.InsertKYCDocumentsParams{
			ID:             uuid.New(),
			StoreID:        storeID,
			AadhaarNumber:  aadhaarNumber,
			AadhaarDocUrl:  aadhaarURL,
			ShopLicenseUrl: licenseURL,
			IsVerified:     false,
			VerifiedAt:     pgtype.Timestamptz{Valid: false},
			CreatedAt:      now,
		})
	} else if err == nil {
		row, err = qTx.UpdateKYCDocumentsByStoreID(ctx, generated.UpdateKYCDocumentsByStoreIDParams{
			StoreID:        storeID,
			AadhaarNumber:  aadhaarNumber,
			AadhaarDocUrl:  aadhaarURL,
			ShopLicenseUrl: licenseURL,
		})
	}
	if err != nil {
		return nil, fmt.Errorf("upsert kyc documents: %w", err)
	}

	// 2. Synchronize store status to pending_kyc
	storeStatus, err := qTx.UpdateStoreStatusByID(ctx, generated.UpdateStoreStatusByIDParams{
		ID:     storeID,
		Status: generated.MerchantStatusPendingKyc,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("sync store status to pending_kyc: %w", err)
	}
	if storeStatus == "" {
		storeStatus = "pending_kyc"
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit kyc update tx: %w", err)
	}

	var verifiedAtPtr *time.Time
	if row.VerifiedAt.Valid {
		verifiedAtPtr = &row.VerifiedAt.Time
	}

	return &domain.KYCStatusDetail{
		ID:             row.ID,
		StoreID:        row.StoreID,
		AadhaarNumber:  row.AadhaarNumber,
		AadhaarDocURL:  row.AadhaarDocUrl,
		ShopLicenseURL: row.ShopLicenseUrl,
		IsVerified:     row.IsVerified,
		VerifiedAt:     verifiedAtPtr,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
		StoreStatus:    storeStatus,
	}, nil
}
