package domain

import (
	"context"
	"time"
	"github.com/google/uuid"
)

type Staff struct {
	ID           uuid.UUID
	Email        string
	FullName     string
	PasswordHash string
	Status       UserStatus
	TokenVersion int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ListStaffResult struct {
	Staff      []Staff
	TotalCount int64
}

type StaffRepository interface {
	Create(ctx context.Context, staff *Staff) error
	GetByID(ctx context.Context, id uuid.UUID) (*Staff, error)
	GetByEmail(ctx context.Context, email string) (*Staff, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	UpdateStatus(ctx context.Context, staffID uuid.UUID, status UserStatus) error
	GetTokenVersion(ctx context.Context, staffID uuid.UUID) (int64, error)
	IncrementTokenVersion(ctx context.Context, staffID uuid.UUID) (int64, error)
	SoftDelete(ctx context.Context, staffID uuid.UUID) error
	Update(ctx context.Context, staffID uuid.UUID, fullName, email, passwordHash *string) (*Staff, error)
	ListStaff(ctx context.Context, filter ListUsersFilter, page Pagination) (*ListStaffResult, error)

	// Role management
	AssignRole(ctx context.Context, staffID, roleID uuid.UUID) error
	RemoveRole(ctx context.Context, staffID, roleID uuid.UUID) error
	GetRoles(ctx context.Context, staffID uuid.UUID) ([]Role, error)
}
