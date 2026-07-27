package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/auth/internal/domain"
	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"
)

type AdminRoleUsecase interface {
	CreateRole(ctx context.Context, req CreateRoleRequest) (*RoleSummary, error)
	GetRole(ctx context.Context, roleID uuid.UUID) (*RoleSummary, error)
	ListRoles(ctx context.Context) ([]RoleSummary, error)
	UpdateRole(ctx context.Context, roleID uuid.UUID, req UpdateRoleRequest) (*RoleSummary, error)
	DeleteRole(ctx context.Context, roleID uuid.UUID) error

	AssignRoleToUser(ctx context.Context, req AssignRoleRequest) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]RoleSummary, error)
}

type adminRoleUsecase struct {
	roleRepo domain.RoleRepository
	userRepo domain.UserRepository
	idGen    IDGenerator
}


func NewAdminRoleUsecase(roleRepo domain.RoleRepository, userRepo domain.UserRepository, idGen IDGenerator) AdminRoleUsecase {
	return &adminRoleUsecase{roleRepo: roleRepo, userRepo: userRepo, idGen: idGen}
}

// ── Role CRUD ─────────────────────────────────────────────────────

func (u *adminRoleUsecase) CreateRole(ctx context.Context, req CreateRoleRequest) (*RoleSummary, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, apperrors.ErrInvalidInput
	}

	existing, err := u.roleRepo.GetByName(ctx, name)
	if err != nil && err != apperrors.ErrNotFound {
		return nil, fmt.Errorf("check existing role: %w", err)
	}
	if existing != nil {
		return nil, apperrors.ErrAlreadyExists
	}

	roleID, err := u.idGen.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate role id: %w", err)
	}

	role := &domain.Role{
		ID:          roleID,
		Name:        name,
		Description: req.Description,
	}

	if err := u.roleRepo.Create(ctx, role); err != nil {
		return nil, fmt.Errorf("create role: %w", err)
	}

	return toRoleSummary(*role), nil
}

func (u *adminRoleUsecase) GetRole(ctx context.Context, roleID uuid.UUID) (*RoleSummary, error) {
	role, err := u.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	return toRoleSummary(*role), nil
}

func (u *adminRoleUsecase) ListRoles(ctx context.Context) ([]RoleSummary, error) {
	roles, err := u.roleRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list roles: %w", err)
	}

	summaries := make([]RoleSummary, len(roles))
	for i, r := range roles {
		summaries[i] = *toRoleSummary(r)
	}
	return summaries, nil
}

func (u *adminRoleUsecase) UpdateRole(ctx context.Context, roleID uuid.UUID, req UpdateRoleRequest) (*RoleSummary, error) {
	if req.Name == nil && req.Description == nil {
		return nil, apperrors.ErrInvalidInput
	}

	role, err := u.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if trimmed == "" {
			return nil, apperrors.ErrInvalidInput
		}
		role.Name = trimmed
	}
	if req.Description != nil {
		role.Description = req.Description
	}

	if err := u.roleRepo.Update(ctx, role); err != nil {
		return nil, fmt.Errorf("update role: %w", err)
	}

	return toRoleSummary(*role), nil
}

func (u *adminRoleUsecase) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	role, err := u.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return err
	}
	if isProtectedRole(role.Name) {
		return apperrors.NewAppError(403, "ROLE_PROTECTED", fmt.Sprintf("the %q role cannot be deleted", role.Name), nil)
	}

	if err := u.roleRepo.Delete(ctx, roleID); err != nil {
		return fmt.Errorf("delete role: %w", err)
	}
	return nil
}

// ── User-Role Assignment ─────────────────────────────────────────

func (u *adminRoleUsecase) AssignRoleToUser(ctx context.Context, req AssignRoleRequest) error {
	if _, err := u.userRepo.GetByID(ctx, req.UserID); err != nil {
		return err
	}
	if _, err := u.roleRepo.GetByID(ctx, req.RoleID); err != nil {
		return err
	}

	if err := u.roleRepo.AssignRoleToUser(ctx, req.UserID, req.RoleID, req.GrantedBy); err != nil {
		return fmt.Errorf("assign role to user: %w", err)
	}
	return nil
}

func (u *adminRoleUsecase) RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error {
	if err := u.roleRepo.RemoveRoleFromUser(ctx, userID, roleID); err != nil {
		return fmt.Errorf("remove role from user: %w", err)
	}
	return nil
}

func (u *adminRoleUsecase) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]RoleSummary, error) {
	roles, err := u.roleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user roles: %w", err)
	}

	summaries := make([]RoleSummary, len(roles))
	for i, r := range roles {
		summaries[i] = *toRoleSummary(r)
	}
	return summaries, nil
}

// ── Helpers ──────────────────────────────────────────────────────

var protectedRoles = map[string]bool{
	"user":     true,
	"merchant": true,
	"admin":    true,
}

func isProtectedRole(name string) bool {
	return protectedRoles[name]
}

func toRoleSummary(r domain.Role) *RoleSummary {
	s := &RoleSummary{
		ID:        r.ID.String(),
		Name:      r.Name,
		CreatedAt: r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: r.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if r.Description != nil {
		s.Description = *r.Description
	}
	return s
}