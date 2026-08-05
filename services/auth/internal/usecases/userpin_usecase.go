package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/domain"
	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"
)

const maxPinAttempts = 5
const pinLockDuration = 15 * time.Minute

type UserPinUsecase interface {
	SetPIN(ctx context.Context, userID uuid.UUID, req SetPINRequest) error
	UpdatePIN(ctx context.Context, userID uuid.UUID, req UpdatePINRequest) error
	VerifyPIN(ctx context.Context, userID uuid.UUID, req VerifyPINRequest) error
}

type userPinUsecase struct {
	pinRepo domain.UserPinRepository
	hasher  PasswordHasher
	idGen   IDGenerator
}

func NewUserPinUsecase(pinRepo domain.UserPinRepository, hasher PasswordHasher, idGen IDGenerator) UserPinUsecase {
	return &userPinUsecase{
		pinRepo: pinRepo,
		hasher:  hasher,
		idGen:   idGen,
	}
}

func (u *userPinUsecase) SetPIN(ctx context.Context, userID uuid.UUID, req SetPINRequest) error {
	if len(req.PIN) != 6 {
		return apperrors.ErrInvalidInput
	}
	exists, err := u.pinRepo.Exists(ctx, userID)
	if err != nil {
		return fmt.Errorf("check pin existence: %w", err)
	}
	if exists {
		return apperrors.ErrAlreadyExists
	}

	pinHash, err := u.hasher.Hash(req.PIN)
	if err != nil {
		return fmt.Errorf("hash pin: %w", err)
	}

	pinID, err := u.idGen.NewV7()
	if err != nil {
		return fmt.Errorf("generate pin id: %w", err)
	}

	pin := &domain.UserPin{
		ID:      pinID,
		UserID:  userID,
		PinHash: pinHash,
	}

	if err := u.pinRepo.Create(ctx, pin); err != nil {
		return fmt.Errorf("create pin: %w", err)
	}

	return nil
}

func (u *userPinUsecase) UpdatePIN(ctx context.Context, userID uuid.UUID, req UpdatePINRequest) error {
	if len(req.NewPIN) != 6 {
		return apperrors.ErrInvalidInput
	}
	pin, err := u.pinRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get pin: %w", err)
	}

	if pin.LockedUntil != nil && pin.LockedUntil.After(time.Now()) {
		return apperrors.ErrForbidden
	}

	ok, err := u.hasher.Verify(pin.PinHash, req.OldPIN)
	if err != nil {
		return fmt.Errorf("verify old pin: %w", err)
	}
	if !ok {
		return apperrors.ErrInvalidInput
	}

	newPinHash, err := u.hasher.Hash(req.NewPIN)
	if err != nil {
		return fmt.Errorf("hash new pin: %w", err)
	}

	if err := u.pinRepo.UpdateHash(ctx, userID, newPinHash); err != nil {
		return fmt.Errorf("update pin hash: %w", err)
	}

	return nil
}

func (u *userPinUsecase) VerifyPIN(ctx context.Context, userID uuid.UUID, req VerifyPINRequest) error {
	pin, err := u.pinRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get pin: %w", err)
	}

	if pin.LockedUntil != nil && pin.LockedUntil.After(time.Now()) {
		return apperrors.ErrForbidden
	}

	ok, err := u.hasher.Verify(pin.PinHash, req.PIN)
	if err != nil {
		return fmt.Errorf("verify pin hash: %w", err)
	}

	if !ok {
		attempts, err := u.pinRepo.IncrementFailedAttempts(ctx, userID)
		if err != nil {
			return fmt.Errorf("increment failed attempts: %w", err)
		}
		if attempts >= maxPinAttempts {
			lockUntil := time.Now().Add(pinLockDuration)
			if err := u.pinRepo.Lock(ctx, userID, lockUntil); err != nil {
				return fmt.Errorf("lock pin: %w", err)
			}
			return apperrors.ErrForbidden
		}
		return apperrors.ErrInvalidInput
	}

	if pin.FailedAttempts > 0 {
		if err := u.pinRepo.ResetFailedAttempts(ctx, userID); err != nil {
			return fmt.Errorf("reset failed attempts: %w", err)
		}
	}

	return nil
}
