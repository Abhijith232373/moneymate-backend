package usecases

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/domain"
	"github.com/moneymate-2026/moneymate-backend/shared/pkg/qrcode"
)

// StoreUseCase orchestrates merchant workflows.
type StoreUseCase struct {
	repo domain.MerchantRepository
}

// NewStoreUseCase constructs a usecase with repository dependencies.
func NewStoreUseCase(repo domain.MerchantRepository) *StoreUseCase {
	return &StoreUseCase{repo: repo}
}

type RegisterStoreInput struct {
	OwnerID           string
	OwnerName         string
	ContactEmail      string
	MobileNumber      string
	LegalName         string
	DBAName           *string
	Type              string
	TaxID             *string
	RegisteredAddress string
	AadhaarNumber     string
	AadhaarDocURL     string
	ShopLicenseURL    string
}

// ProcessRegistration applies validation and executes state persistence.
func (uc *StoreUseCase) ProcessRegistration(ctx context.Context, in RegisterStoreInput) (*domain.Store, error) {
	ownerUUID, err := uuid.Parse(in.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("invalid owner UUID format: %w", err)
	}

	displayID, err := generateDisplayID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate secure display ID: %w", err)
	}

	vpa := generateVPA(in.ContactEmail)
	qrCodeBase64, err := qrcode.GenerateBase64(vpa)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	store := &domain.Store{
		OwnerID:           ownerUUID,
		OwnerName:         strings.TrimSpace(in.OwnerName),
		ContactEmail:      strings.ToLower(strings.TrimSpace(in.ContactEmail)),
		MobileNumber:      strings.TrimSpace(in.MobileNumber),
		LegalName:         strings.TrimSpace(in.LegalName),
		DBAName:           in.DBAName,
		Type:              in.Type,
		TaxID:             in.TaxID,
		RegisteredAddress: strings.TrimSpace(in.RegisteredAddress),
		DisplayID:         displayID,
		VPA:               vpa,
		QRCodeBase64:      qrCodeBase64,
	}

	createdStore, err := uc.repo.RegisterStore(ctx, store)
	if err != nil {
		return nil, fmt.Errorf("failed to register store: %w", err)
	}

	kyc := &domain.KYCDocument{
		StoreID:        createdStore.ID,
		AadhaarNumber:  in.AadhaarNumber,
		AadhaarDocURL:  in.AadhaarDocURL,
		ShopLicenseURL: in.ShopLicenseURL,
	}

	if err := uc.repo.SubmitKYC(ctx, kyc); err != nil {
		return nil, fmt.Errorf("failed to submit KYC documents: %w", err)
	}

	return createdStore, nil
}

// GetStore retrieves a store by owner ID.
func (uc *StoreUseCase) GetStore(ctx context.Context, ownerID string) (*domain.Store, error) {
	ownerUUID, err := uuid.Parse(ownerID)
	if err != nil {
		return nil, fmt.Errorf("invalid owner UUID format: %w", err)
	}

	return uc.repo.GetStoreByOwnerID(ctx, ownerUUID)
}

// generateDisplayID yields a collision-resistant MM-XXXX-XX identifier.
func generateDisplayID() (string, error) {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	hexStr := strings.ToUpper(hex.EncodeToString(b))
	return fmt.Sprintf("MM-%s-%s", hexStr[:4], hexStr[4:]), nil
}

// generateVPA creates a unique VPA like emailprefix+random@moneymate
func generateVPA(email string) string {
	parts := strings.Split(email, "@")
	prefix := parts[0]
	if len(prefix) > 10 {
		prefix = prefix[:10]
	}
	
	b := make([]byte, 2)
	rand.Read(b)
	hexStr := strings.ToLower(hex.EncodeToString(b))
	
	return fmt.Sprintf("%s%s@moneymate", prefix, hexStr)
}

// GetPendingStores retrieves all merchants in the pending_kyc status.
func (uc *StoreUseCase) GetPendingStores(ctx context.Context) ([]*domain.Store, error) {
	return uc.repo.GetPendingStores(ctx)
}

type UpdateProfileInput struct {
	StoreID           string
	OwnerID           string
	BusinessName      string
	DBAName           string
	Address           string
	BusinessType      string
	TaxID             string
	OwnerName         string
	Email             string
	Mobile            string
	ProfileImage      string
}

// GetProfile retrieves the store profile by store UUID or owner UUID.
func (uc *StoreUseCase) GetProfile(ctx context.Context, id string) (*domain.Store, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID format for profile lookup: %w", err)
	}

	store, err := uc.repo.GetStoreProfileByStoreID(ctx, parsedUUID)
	if err == nil && store != nil && store.ID != uuid.Nil {
		return store, nil
	}

	return uc.repo.GetStoreProfileByOwnerID(ctx, parsedUUID)
}

// UpdateProfile updates store profile fields.
func (uc *StoreUseCase) UpdateProfile(ctx context.Context, in UpdateProfileInput) (*domain.Store, error) {
	var storeID, ownerID uuid.UUID
	if in.StoreID != "" {
		storeID, _ = uuid.Parse(in.StoreID)
	}
	if in.OwnerID != "" {
		ownerID, _ = uuid.Parse(in.OwnerID)
	}

	if storeID == uuid.Nil && ownerID == uuid.Nil {
		return nil, fmt.Errorf("either StoreID or OwnerID must be provided")
	}

	var dba *string
	if in.DBAName != "" {
		dba = &in.DBAName
	}
	var tax *string
	if in.TaxID != "" {
		tax = &in.TaxID
	}

	store := &domain.Store{
		ID:                storeID,
		OwnerID:           ownerID,
		LegalName:         strings.TrimSpace(in.BusinessName),
		DBAName:           dba,
		RegisteredAddress: strings.TrimSpace(in.Address),
		Type:              in.BusinessType,
		TaxID:             tax,
		OwnerName:         strings.TrimSpace(in.OwnerName),
		ContactEmail:      strings.ToLower(strings.TrimSpace(in.Email)),
		MobileNumber:      strings.TrimSpace(in.Mobile),
		LogoURL:           in.ProfileImage,
	}

	if storeID != uuid.Nil {
		res, err := uc.repo.UpdateStoreProfileByStoreID(ctx, store)
		if err == nil && res != nil && res.ID != uuid.Nil {
			return res, nil
		}
	}

	if ownerID != uuid.Nil {
		return uc.repo.UpdateStoreProfileByOwnerID(ctx, store)
	}

	return nil, fmt.Errorf("failed to update store profile")
}