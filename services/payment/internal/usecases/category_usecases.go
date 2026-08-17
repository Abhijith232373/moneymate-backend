package usecases

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/services/payment/internal/domain"
	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"
)

type CategoryUsecase interface {
	Create(ctx context.Context, userID, name string) (*domain.Category, error)
	List(ctx context.Context, userID string) ([]domain.Category, error)
	Update(ctx context.Context, id, userID, name string) (*domain.Category, error)
	Delete(ctx context.Context, id, userID string) error
}

type categoryUsecase struct {
	categories domain.CategoryRepository
}

func NewCategoryUsecase(categories domain.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{categories: categories}
}

func (u *categoryUsecase) Create(ctx context.Context, userID, name string) (*domain.Category, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperrors.ErrInvalidInput
	}
	name = strings.TrimSpace(name)
	if name == "" || strings.EqualFold(name, "other") {
		// "Other" is reserved — it's the implicit fallback, not a real category.
		return nil, apperrors.ErrInvalidInput
	}
	cat, err := u.categories.Create(ctx, uid, name)
	if err != nil {
		return nil, fmt.Errorf("create category: %w", err)
	}
	return cat, nil
}

func (u *categoryUsecase) List(ctx context.Context, userID string) ([]domain.Category, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperrors.ErrInvalidInput
	}
	return u.categories.ListByUser(ctx, uid)
}

func (u *categoryUsecase) Update(ctx context.Context, id, userID, name string) (*domain.Category, error) {
	cid, err := uuid.Parse(id)
	if err != nil {
		return nil, apperrors.ErrInvalidInput
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperrors.ErrInvalidInput
	}
	name = strings.TrimSpace(name)
	if name == "" || strings.EqualFold(name, "other") {
		return nil, apperrors.ErrInvalidInput
	}
	return u.categories.Update(ctx, cid, uid, name)
}

func (u *categoryUsecase) Delete(ctx context.Context, id, userID string) error {
	cid, err := uuid.Parse(id)
	if err != nil {
		return apperrors.ErrInvalidInput
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return apperrors.ErrInvalidInput
	}
	return u.categories.Delete(ctx, cid, uid)
}