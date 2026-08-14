package usecases

import (
	"context"

	"github.com/google/uuid"
	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"

	"github.com/moneymate-2026/moneymate-backend/services/notification/internal/domain"
)

type DeviceUsecase struct {
	deviceRepo domain.DeviceTokenRepository
}

func NewDeviceUsecase(r domain.DeviceTokenRepository) *DeviceUsecase {
	return &DeviceUsecase{deviceRepo: r}
}

type RegisterDeviceInput struct {
	RecipientType domain.RecipientType
	RecipientID   uuid.UUID
	DeviceID      string
	Token         string
	Platform      string
	AppVersion    string
}

func (uc *DeviceUsecase) Register(ctx context.Context, in RegisterDeviceInput) (*domain.DeviceToken, error) {
	if in.DeviceID == "" || in.Token == "" {
		return nil, apperrors.ErrInvalidInput
	}
	switch in.Platform {
	case "ios", "android", "web":
	default:
		return nil, apperrors.ErrInvalidInput
	}
	return uc.deviceRepo.Upsert(ctx, &domain.DeviceToken{
		RecipientType: in.RecipientType,
		RecipientID:   in.RecipientID,
		DeviceID:      in.DeviceID,
		Token:         in.Token,
		Platform:      in.Platform,
		AppVersion:    in.AppVersion,
	})
}

func (uc *DeviceUsecase) Revoke(ctx context.Context, recipientType domain.RecipientType, recipientID uuid.UUID, deviceID string) error {
	if deviceID == "" {
		return apperrors.ErrInvalidInput
	}
	return uc.deviceRepo.DeactivateByDevice(ctx, recipientType, recipientID, deviceID)
}
