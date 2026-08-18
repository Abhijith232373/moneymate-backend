package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/auth/internal/domain"
	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"
)

type StaffUsecase interface {
	CreateStaff(ctx context.Context, req CreateUserRequest) (*UserDetail, error)
	ListStaff(ctx context.Context, req ListUsersRequest) (*ListUsersResponse, error)
	GetStaff(ctx context.Context, staffID uuid.UUID) (*UserDetail, error)
	UpdateStaff(ctx context.Context, staffID uuid.UUID, req UpdateUserRequest) (*UserDetail, error)
	UpdateStaffStatus(ctx context.Context, staffID uuid.UUID, status string) error
	DeleteStaff(ctx context.Context, staffID uuid.UUID) error
}

type staffUsecase struct {
	staffRepo domain.StaffRepository
	roleRepo  domain.RoleRepository
	idGen     IDGenerator
	hasher    PasswordHasher
}

func NewStaffUsecase(staffRepo domain.StaffRepository, roleRepo domain.RoleRepository, hasher PasswordHasher, idGen IDGenerator) StaffUsecase {
	return &staffUsecase{staffRepo: staffRepo, roleRepo: roleRepo, hasher: hasher, idGen: idGen}
}

func (u *staffUsecase) CreateStaff(ctx context.Context, req CreateUserRequest) (*UserDetail, error) {
	email := normalizeEmail(req.Email)
	if email == "" || !strings.Contains(email, "@") {
		return nil, apperrors.ErrInvalidInput
	}
	if err := validatePassword(req.Password); err != nil {
		return nil, err
	}

	exists, err := u.staffRepo.EmailExists(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("check email exists: %w", err)
	}
	if exists {
		return nil, apperrors.ErrEmailAlreadyTaken
	}

	passwordHash, err := u.hasher.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	role, err := u.roleRepo.GetByName(ctx, strings.ToLower(req.Role))
	if err != nil {
		return nil, fmt.Errorf("resolve role %q: %w", req.Role, err)
	}

	staffID, err := u.idGen.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate staff id: %w", err)
	}

	staff := &domain.Staff{
		ID:           staffID,
		Email:        email,
		FullName:     strings.TrimSpace(req.FullName),
		PasswordHash: passwordHash,
		Status:       domain.UserStatusActive,
	}

	if err := u.staffRepo.Create(ctx, staff); err != nil {
		return nil, fmt.Errorf("create staff: %w", err)
	}

	if err := u.staffRepo.AssignRole(ctx, staffID, role.ID); err != nil {
		return nil, fmt.Errorf("assign role to staff: %w", err)
	}

	return &UserDetail{
		AdminUserSummary: toStaffSummary(*staff),
		Roles:            []string{req.Role},
	}, nil
}

func (u *staffUsecase) ListStaff(ctx context.Context, req ListUsersRequest) (*ListUsersResponse, error) {
	result, err := u.staffRepo.ListStaff(ctx, domain.ListUsersFilter{
		Status:   req.Status,
		Search:   req.Search,
		SortBy:   req.SortBy,
		SortDesc: req.SortDesc,
	}, domain.Pagination{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("list staff: %w", err)
	}

	summaries := make([]AdminUserSummary, len(result.Staff))
	for i, s := range result.Staff {
		summary := toStaffSummary(s)
		roles, err := u.staffRepo.GetRoles(ctx, s.ID)
		if err == nil && len(roles) > 0 {
			summary.Role = roles[0].Name
		} else {
			summary.Role = "staff"
		}
		summaries[i] = summary
	}

	return &ListUsersResponse{
		Users:      summaries,
		TotalCount: result.TotalCount,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}

func (u *staffUsecase) GetStaff(ctx context.Context, staffID uuid.UUID) (*UserDetail, error) {
	staff, err := u.staffRepo.GetByID(ctx, staffID)
	if err != nil {
		return nil, err
	}

	roles, err := u.staffRepo.GetRoles(ctx, staffID)
	if err != nil {
		return nil, fmt.Errorf("get staff roles: %w", err)
	}

	roleNames := make([]string, len(roles))
	for i, r := range roles {
		roleNames[i] = r.Name
	}

	return &UserDetail{
		AdminUserSummary: toStaffSummary(*staff),
		Roles:            roleNames,
	}, nil
}

func (u *staffUsecase) UpdateStaff(ctx context.Context, staffID uuid.UUID, req UpdateUserRequest) (*UserDetail, error) {
	if req.FullName == nil && req.Email == nil && req.Password == nil {
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

	updated, err := u.staffRepo.Update(ctx, staffID, req.FullName, email, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("update staff: %w", err)
	}

	roles, err := u.staffRepo.GetRoles(ctx, staffID)
	if err != nil {
		return nil, fmt.Errorf("get staff roles: %w", err)
	}
	roleNames := make([]string, len(roles))
	for i, r := range roles {
		roleNames[i] = r.Name
	}

	return &UserDetail{
		AdminUserSummary: toStaffSummary(*updated),
		Roles:            roleNames,
	}, nil
}

func (u *staffUsecase) UpdateStaffStatus(ctx context.Context, staffID uuid.UUID, status string) error {
	if !validStatuses[status] {
		return apperrors.ErrInvalidInput
	}
	if err := u.staffRepo.UpdateStatus(ctx, staffID, domain.UserStatus(status)); err != nil {
		return fmt.Errorf("update staff status: %w", err)
	}
	return nil
}

func (u *staffUsecase) DeleteStaff(ctx context.Context, staffID uuid.UUID) error {
	if err := u.staffRepo.SoftDelete(ctx, staffID); err != nil {
		return fmt.Errorf("delete staff: %w", err)
	}
	return nil
}

func toStaffSummary(s domain.Staff) AdminUserSummary {
	return AdminUserSummary{
		ID:              s.ID.String(),
		Email:           s.Email,
		FullName:        s.FullName,
		Status:          string(s.Status),
		IsEmailVerified: true, // Internal staff are verified by default
		IsPhoneVerified: false,
		CreatedAt:       s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
