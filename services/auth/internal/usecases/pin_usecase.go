package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/auth/internal/domain"
	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"
)

const (
	pinMinLength   = 4
	pinMaxLength   = 6
	maxPINAttempts = 5
	pinLockoutTime = 15 * time.Minute
)

type PinUsecase interface {
	SetPIN(ctx context.Context, userID uuid.UUID, pin string) error
	VerifyPIN(ctx context.Context, userID uuid.UUID, pin string) error
}

type pinUsecase struct {
	pins   domain.UserPinRepository
	hasher PasswordHasher
}

func NewPinUsecase(pins domain.UserPinRepository, hasher PasswordHasher) PinUsecase {
	return &pinUsecase{pins: pins, hasher: hasher}
}

func validPIN(pin string) bool {
	if len(pin) < pinMinLength || len(pin) > pinMaxLength {
		return false
	}
	for _, r := range pin {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// SetPIN creates the PIN on first use, or re-hashes it on reset.
func (u *pinUsecase) SetPIN(ctx context.Context, userID uuid.UUID, pin string) error {
	if !validPIN(pin) {
		return apperrors.ErrInvalidInput
	}
	hash, err := u.hasher.Hash(pin)
	if err != nil {
		return err
	}
	exists, err := u.pins.Exists(ctx, userID)
	if err != nil {
		return err
	}
	if exists {
		return u.pins.UpdateHash(ctx, userID, hash)
	}
	return u.pins.Create(ctx, &domain.UserPin{ID: uuid.New(), UserID: userID, PinHash: hash})
}

// VerifyPIN checks the PIN and enforces the 5-attempt / 15-minute lockout.
func (u *pinUsecase) VerifyPIN(ctx context.Context, userID uuid.UUID, pin string) error {
	record, err := u.pins.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.ErrPinNotSet
		}
		return err
	}

	if record.LockedUntil != nil && time.Now().Before(*record.LockedUntil) {
		return apperrors.ErrPinLocked
	}

	ok, err := u.hasher.Verify(record.PinHash, pin)
	if err != nil {
		return err
	}
	if !ok {
		attempts, err := u.pins.IncrementFailedAttempts(ctx, userID)
		if err != nil {
			return err
		}
		if attempts >= maxPINAttempts {
			until := time.Now().Add(pinLockoutTime)
			if err := u.pins.Lock(ctx, userID, until); err != nil {
				return err
			}
		}
		return apperrors.ErrInvalidPIN
	}

	return u.pins.ResetFailedAttempts(ctx, userID)
}
