package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// KYCStatusDetail represents the full compliance standing and document verification record of a merchant store.
type KYCStatusDetail struct {
	ID             uuid.UUID
	StoreID        uuid.UUID
	AadhaarNumber  string
	AadhaarDocURL  string
	ShopLicenseURL string
	IsVerified     bool
	VerifiedAt     *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	StoreStatus    string
}

// KYCRepository defines the strict data access contract for querying compliance states and atomically updating verification documentation.
type KYCRepository interface {
	GetKYCStatusByStoreID(ctx context.Context, storeID uuid.UUID) (*KYCStatusDetail, error)
	UpdateKYCDocuments(ctx context.Context, storeID uuid.UUID, aadhaarNumber, aadhaarURL, licenseURL string) (*KYCStatusDetail, error)
}

