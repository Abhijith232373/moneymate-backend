package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/auth/internal/domain"
	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"
)



type AdminUserUsecase interface {
	ListUsers(ctx context.Context, req ListUsersRequest) (*ListUsersResponse, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*UserDetail, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, req UpdateUserRequest) (*UserDetail, error)
	UpdateUserStatus(ctx context.Context, userID uuid.UUID, status string) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type adminUserUsecase struct {
	userRepo domain.UserRepository
	roleRepo domain.RoleRepository
	hasher   PasswordHasher
}

func NewAdminUserUsecase(userRepo domain.UserRepository, roleRepo domain.RoleRepository, hasher PasswordHasher) AdminUserUsecase {
	return &adminUserUsecase{userRepo: userRepo, roleRepo: roleRepo, hasher: hasher}
}

// ── List ─────────────────────────────────────────────────────────
func (u *adminUserUsecase) ListUsers(ctx context.Context, req ListUsersRequest) (*ListUsersResponse, error) {
	result, err := u.userRepo.ListUsers(ctx, domain.ListUsersFilter{
		Status:   req.Status,
		Search:   req.Search,
		SortBy:   req.SortBy,
		SortDesc: req.SortDesc,
	}, domain.Pagination{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	summaries := make([]AdminUserSummary, len(result.Users))
	for i, us := range result.Users {
		summaries[i] = toAdminUserSummary(us)
	}

	return &ListUsersResponse{
		Users:      summaries,
		TotalCount: result.TotalCount,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}

// ── Get single user (with roles) ────────────────────────────────────

func (u *adminUserUsecase) GetUser(ctx context.Context, userID uuid.UUID) (*UserDetail, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	roles, err := u.roleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user roles: %w", err)
	}

	roleNames := make([]string, len(roles))
	for i, r := range roles {
		roleNames[i] = r.Name
	}

	return &UserDetail{
		AdminUserSummary: toAdminUserSummary(*user),
		Roles:       roleNames,
	}, nil
}

// ── Update ───────────────────────────────────────────────────────

func (u *adminUserUsecase) UpdateUser(ctx context.Context, userID uuid.UUID, req UpdateUserRequest) (*UserDetail, error) {
	if req.FullName == nil && req.Email == nil && req.Phone == nil && req.Password == nil {
		return nil, apperrors.ErrInvalidInput
	}

	var email *string
	if req.Email != nil {
		normalized := normalizeEmail(*req.Email)
		if normalized == "" || !strings.Contains(normalized, "@") {
			return nil, apperrors.ErrInvalidInput
		}
		email = &normalized
	}

	var passwordHash *string
	if req.Password != nil {
		if err := validatePassword(*req.Password); err != nil {
			return nil, err
		}
		hashed, err := u.hasher.Hash(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("hash password: %w", err)
		}
		passwordHash = &hashed
	}

	updated, err := u.userRepo.AdminUpdate(ctx, userID, req.FullName, email, req.Phone, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	roles, err := u.roleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user roles: %w", err)
	}
	roleNames := make([]string, len(roles))
	for i, r := range roles {
		roleNames[i] = r.Name
	}

	return &UserDetail{
		AdminUserSummary: toAdminUserSummary(*updated),
		Roles:       roleNames,
	}, nil
}

// ── Status change ────────────────────────────────────────────────

var validStatuses = map[string]bool{
	string(domain.UserStatusPending):   true,
	string(domain.UserStatusActive):    true,
	string(domain.UserStatusSuspended): true,
}

func (u *adminUserUsecase) UpdateUserStatus(ctx context.Context, userID uuid.UUID, status string) error {
	if !validStatuses[status] {
		return apperrors.ErrInvalidInput
	}
	if err := u.userRepo.UpdateStatus(ctx, userID, domain.UserStatus(status)); err != nil {
		return fmt.Errorf("update user status: %w", err)
	}
	return nil
}

// ── Delete (soft) ────────────────────────────────────────────────

func (u *adminUserUsecase) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if err := u.userRepo.SoftDelete(ctx, userID); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

// ── Helpers ──────────────────────────────────────────────────────

func toAdminUserSummary(u domain.User) AdminUserSummary {
	s := AdminUserSummary{
		ID:              u.ID.String(),
		Email:           u.Email,
		FullName:        u.FullName,
		Handle:          u.Handle,
		Status:          string(u.Status),
		IsEmailVerified: u.IsEmailVerified,
		IsPhoneVerified: u.IsPhoneVerified,
		CreatedAt:       u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if u.Phone != nil {
		s.Phone = *u.Phone
	}
	return s
}