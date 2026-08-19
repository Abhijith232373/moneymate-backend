package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/auth/internal/domain"
	s3util "github.com/moneymate-2026/moneymate-backend/shared/pkg/S3"
)

type ProfilePictureUsecase interface {
	PresignUpload(ctx context.Context, userID uuid.UUID, contentType string) (*s3util.PresignedUpload, error)
	SetProfilePicture(ctx context.Context, userID uuid.UUID, url string) (*domain.User, error)
}

type profilePictureUsecase struct {
	userRepo domain.UserRepository
	s3       *s3util.Client
}

func NewProfilePictureUsecase(userRepo domain.UserRepository, s3 *s3util.Client) ProfilePictureUsecase {
	return &profilePictureUsecase{userRepo: userRepo, s3: s3}
}

func (u *profilePictureUsecase) PresignUpload(ctx context.Context, userID uuid.UUID, contentType string) (*s3util.PresignedUpload, error) {
	upload, err := u.s3.PresignProfilePictureUpload(ctx, "users", userID, contentType)
	if err != nil {
		return nil, fmt.Errorf("presign profile picture upload: %w", err)
	}
	return upload, nil
}

func (u *profilePictureUsecase) SetProfilePicture(ctx context.Context, userID uuid.UUID, url string) (*domain.User, error) {
	if !u.s3.IsOwnedURL(url, "users", userID) {
		return nil, fmt.Errorf("url does not belong to this user")
	}
	user, err := u.userRepo.UpdateProfilePicture(ctx, userID, url)
	if err != nil {
		return nil, fmt.Errorf("set profile picture: %w", err)
	}
	return user, nil
}