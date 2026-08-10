package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/moneymate-2026/moneymate-backend/auth/internal/domain"
	apperrors "github.com/moneymate-2026/moneymate-backend/shared/pkg/errors"
)

type PermissionUsecase interface {
	CreatePermission(ctx context.Context, req CreatePermissionRequest) (*domain.Permission, error)
	GetPermission(ctx context.Context, id uuid.UUID) (*domain.Permission, error)
	ListPermissions(ctx context.Context) ([]domain.Permission, error)
	UpdatePermission(ctx context.Context, id uuid.UUID, req UpdatePermissionRequest) (*domain.Permission, error)
	DeletePermission(ctx context.Context, id uuid.UUID) error

	AssignPermissionToRole(ctx context.Context, roleID, permissionID uuid.UUID) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID uuid.UUID) error
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]domain.Permission, error)
}

type CreatePermissionRequest struct {
	Name        string
	Description string
}

type UpdatePermissionRequest struct {
	Name        string
	Description string
}

type permissionUsecase struct {
	permRepo domain.PermissionRepository
	roleRepo domain.RoleRepository
	idGen    IDGenerator
}

func NewPermissionUsecase(permRepo domain.PermissionRepository, roleRepo domain.RoleRepository, idGen IDGenerator) PermissionUsecase {
	return &permissionUsecase{permRepo: permRepo, roleRepo: roleRepo, idGen: idGen}
}

func (u *permissionUsecase) CreatePermission(ctx context.Context, req CreatePermissionRequest) (*domain.Permission, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, apperrors.ErrInvalidInput
	}

	existing, err := u.permRepo.GetByName(ctx, name)
	if err != nil && err != apperrors.ErrNotFound {
		return nil, fmt.Errorf("check existing permission: %w", err)
	}
	if existing != nil {
		return nil, apperrors.ErrAlreadyExists
	}

	id, err := u.idGen.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate permission id: %w", err)
	}

	perm := &domain.Permission{ID: id, Name: name}
	if desc := strings.TrimSpace(req.Description); desc != "" {
		perm.Description = &desc
	}

	if err := u.permRepo.Create(ctx, perm); err != nil {
		return nil, fmt.Errorf("create permission: %w", err)
	}
	return perm, nil
}

func (u *permissionUsecase) GetPermission(ctx context.Context, id uuid.UUID) (*domain.Permission, error) {
	perm, err := u.permRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get permission: %w", err)
	}
	return perm, nil
}

func (u *permissionUsecase) ListPermissions(ctx context.Context) ([]domain.Permission, error) {
	perms, err := u.permRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list permissions: %w", err)
	}
	return perms, nil
}

func (u *permissionUsecase) UpdatePermission(ctx context.Context, id uuid.UUID, req UpdatePermissionRequest) (*domain.Permission, error) {
	perm, err := u.permRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get permission: %w", err)
	}

	if name := strings.TrimSpace(req.Name); name != "" {
		perm.Name = name
	}
	if desc := strings.TrimSpace(req.Description); desc != "" {
		perm.Description = &desc
	}

	if err := u.permRepo.Update(ctx, perm); err != nil {
		return nil, fmt.Errorf("update permission: %w", err)
	}
	return perm, nil
}

func (u *permissionUsecase) DeletePermission(ctx context.Context, id uuid.UUID) error {
	if _, err := u.permRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("get permission: %w", err)
	}
	if err := u.permRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete permission: %w", err)
	}
	return nil
}

func (u *permissionUsecase) AssignPermissionToRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	if _, err := u.roleRepo.GetByID(ctx, roleID); err != nil {
		return fmt.Errorf("get role: %w", err)
	}
	if _, err := u.permRepo.GetByID(ctx, permissionID); err != nil {
		return fmt.Errorf("get permission: %w", err)
	}
	if err := u.permRepo.AssignPermissionToRole(ctx, roleID, permissionID); err != nil {
		return fmt.Errorf("assign permission to role: %w", err)
	}
	return nil
}

func (u *permissionUsecase) RemovePermissionFromRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	if err := u.permRepo.RemovePermissionFromRole(ctx, roleID, permissionID); err != nil {
		return fmt.Errorf("remove permission from role: %w", err)
	}
	return nil
}

func (u *permissionUsecase) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]domain.Permission, error) {
	perms, err := u.permRepo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("get role permissions: %w", err)
	}
	return perms, nil
}