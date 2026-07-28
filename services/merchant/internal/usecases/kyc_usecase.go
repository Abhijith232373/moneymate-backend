package usecases

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/domain"
)

// KYCUseCase defines the business logic contract for reviewing compliance documentation and processing document updates.
type KYCUseCase interface {
	// GetStatus retrieves the verification standing and document records for a specific merchant store.
	GetStatus(ctx context.Context, storeID uuid.UUID) (*domain.KYCStatusDetail, error)
	// UpdateDocuments validates document URL formats and identification strings before submitting them for compliance re-review.
	UpdateDocuments(ctx context.Context, storeID uuid.UUID, aadhaarNumber, aadhaarURL, licenseURL string) (*domain.KYCStatusDetail, error)
}

// kycUseCase implements KYCUseCase with injected repository dependencies.
type kycUseCase struct {
	// kycRepo provides persistent data access to compliance document ledgers.
	kycRepo domain.KYCRepository
}

// NewKYCUseCase constructs and returns a new kycUseCase instance with required dependencies.
func NewKYCUseCase(kr domain.KYCRepository) KYCUseCase {
	return &kycUseCase{kycRepo: kr}
}

// GetStatus validates the merchant store UUID format and retrieves compliance records from the repository layer.
func (uc *kycUseCase) GetStatus(ctx context.Context, storeID uuid.UUID) (*domain.KYCStatusDetail, error) {
	if storeID == uuid.Nil {
		return nil, errors.New("invalid store ID")
	}
	return uc.kycRepo.GetKYCStatusByStoreID(ctx, storeID)
}

// UpdateDocuments applies compliance validation rules on uploaded document URIs and Aadhaar identification numbers
// before delegating atomic state updates to the repository layer.
func (uc *kycUseCase) UpdateDocuments(ctx context.Context, storeID uuid.UUID, aadhaarNumber, aadhaarURL, licenseURL string) (*domain.KYCStatusDetail, error) {
	if storeID == uuid.Nil {
		return nil, errors.New("invalid store ID")
	}

	aadhaarNumber = strings.TrimSpace(aadhaarNumber)
	aadhaarURL = strings.TrimSpace(aadhaarURL)
	licenseURL = strings.TrimSpace(licenseURL)

	if aadhaarNumber != "" && len(aadhaarNumber) != 12 {
		return nil, fmt.Errorf("aadhaar number must be exactly 12 digits, got %d", len(aadhaarNumber))
	}
	if aadhaarURL != "" && (!strings.HasPrefix(aadhaarURL, "http://") && !strings.HasPrefix(aadhaarURL, "https://")) {
		return nil, errors.New("aadhaar document url must be a valid HTTP/HTTPS URI")
	}
	if licenseURL != "" && (!strings.HasPrefix(licenseURL, "http://") && !strings.HasPrefix(licenseURL, "https://")) {
		return nil, errors.New("shop license document url must be a valid HTTP/HTTPS URI")
	}

	return uc.kycRepo.UpdateKYCDocuments(ctx, storeID, aadhaarNumber, aadhaarURL, licenseURL)
}
