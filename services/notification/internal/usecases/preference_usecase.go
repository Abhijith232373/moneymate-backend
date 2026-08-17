package usecases

import (
	"context"

	"github.com/google/uuid"
	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"

	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/domain"
)

type PreferenceUsecase struct {
	preferenceRepo domain.PreferenceRepository
}

func NewPreferenceUsecase(r domain.PreferenceRepository) *PreferenceUsecase {
	return &PreferenceUsecase{preferenceRepo: r}
}

func (uc *PreferenceUsecase) List(ctx context.Context, recipientType domain.RecipientType, recipientID uuid.UUID) ([]*domain.Preference, error) {
	return uc.preferenceRepo.List(ctx, recipientType, recipientID)
}

func (uc *PreferenceUsecase) Upsert(ctx context.Context, recipientType domain.RecipientType, recipientID uuid.UUID, category domain.Category, enabled bool) (*domain.Preference, error) {
	if !validCategory(category) {
		return nil, apperrors.ErrInvalidInput
	}
	return uc.preferenceRepo.Upsert(ctx, recipientType, recipientID, category, enabled)
}

func validCategory(c domain.Category) bool {
	switch c {
	case domain.CategoryBillDue, domain.CategoryDebt, domain.CategoryTransfer,
		domain.CategoryMerchant, domain.CategoryCampaign, domain.CategoryOffer,
		domain.CategoryPromo, domain.CategorySystem:
		return true
	}
	return false
}
