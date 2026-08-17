package usecases

import (
	"context"

	"github.com/google/uuid"
	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"

	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/domain"
)

type InboxUsecase struct {
	inboxRepo domain.InboxRepository
}

func NewInboxUsecase(r domain.InboxRepository) *InboxUsecase {
	return &InboxUsecase{inboxRepo: r}
}

func (uc *InboxUsecase) List(ctx context.Context, recipientType domain.RecipientType, recipientID uuid.UUID, page, pageSize int) ([]*domain.InboxMessage, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	return uc.inboxRepo.List(ctx, recipientType, recipientID, int32(pageSize), int32((page-1)*pageSize))
}

// Get is ownership-scoped: the query filters by recipient, so a foreign id returns ErrNotFound.
func (uc *InboxUsecase) Get(ctx context.Context, id uuid.UUID, recipientType domain.RecipientType, recipientID uuid.UUID) (*domain.InboxMessage, error) {
	msg, err := uc.inboxRepo.Get(ctx, id, recipientType, recipientID)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	return msg, nil
}

func (uc *InboxUsecase) MarkRead(ctx context.Context, id uuid.UUID, recipientType domain.RecipientType, recipientID uuid.UUID) error {
	return uc.inboxRepo.MarkRead(ctx, id, recipientType, recipientID)
}
