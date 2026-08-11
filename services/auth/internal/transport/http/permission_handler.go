package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	usecase "github.com/moneymate-2026/moneymate-backend/auth/internal/usecases"
	response "github.com/moneymate-2026/moneymate-backend/shared/pkg/responses"
)

type PermissionHandler struct {
	usecase usecase.PermissionUsecase
}

func NewPermissionHandler(u usecase.PermissionUsecase) *PermissionHandler {
	return &PermissionHandler{usecase: u}
}



func (h *PermissionHandler) Create(c fiber.Ctx) error {
	var req createPermissionRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return response.BadRequest(c, formatValidationErrors(err), "validation failed")
	}

	perm, err := h.usecase.CreatePermission(c.Context(), usecase.CreatePermissionRequest{
		Name: req.Name, Description: req.Description,
	})
	if err != nil {
		return handleError(c, err)
	}
	return response.Created(c, "permission created", perm)
}

func (h *PermissionHandler) List(c fiber.Ctx) error {
	perms, err := h.usecase.ListPermissions(c.Context())
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "permissions listed", perms)
}

func (h *PermissionHandler) Get(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid permission id")
	}
	perm, err := h.usecase.GetPermission(c.Context(), id)
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "permission found", perm)
}

func (h *PermissionHandler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid permission id")
	}
	var req updatePermissionRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request body")
	}
	perm, err := h.usecase.UpdatePermission(c.Context(), id, usecase.UpdatePermissionRequest{
		Name: req.Name, Description: req.Description,
	})
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "permission updated", perm)
}

func (h *PermissionHandler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid permission id")
	}
	if err := h.usecase.DeletePermission(c.Context(), id); err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "permission deleted", nil)
}

func (h *PermissionHandler) AssignToRole(c fiber.Ctx) error {
	var req assignPermissionRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return response.BadRequest(c, formatValidationErrors(err), "validation failed")
	}
	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return response.BadRequest(c, nil, "invalid role id")
	}
	permID, err := uuid.Parse(req.PermissionID)
	if err != nil {
		return response.BadRequest(c, nil, "invalid permission id")
	}
	if err := h.usecase.AssignPermissionToRole(c.Context(), roleID, permID); err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "permission assigned to role", nil)
}

func (h *PermissionHandler) RemoveFromRole(c fiber.Ctx) error {
	roleID, err := uuid.Parse(c.Params("roleId"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid role id")
	}
	permID, err := uuid.Parse(c.Params("permissionId"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid permission id")
	}
	if err := h.usecase.RemovePermissionFromRole(c.Context(), roleID, permID); err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "permission removed from role", nil)
}

func (h *PermissionHandler) GetRolePermissions(c fiber.Ctx) error {
	roleID, err := uuid.Parse(c.Params("roleId"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid role id")
	}
	perms, err := h.usecase.GetRolePermissions(c.Context(), roleID)
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "role permissions listed", perms)
}