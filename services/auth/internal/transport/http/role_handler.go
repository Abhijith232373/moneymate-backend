package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	usecase "github.com/moneymate-2026/moneymate-backend/auth/internal/usecases"
	response "github.com/moneymate-2026/moneymate-backend/shared/pkg/responses"
)

type RoleHandler struct {
	adminRoleUsecase usecase.AdminRoleUsecase
}

func NewRoleHandler(adminRoleUsecase usecase.AdminRoleUsecase) *RoleHandler {
	return &RoleHandler{adminRoleUsecase: adminRoleUsecase}
}



func (h *RoleHandler) CreateRole(c fiber.Ctx) error {
	var req createRoleRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return response.BadRequest(c, formatValidationErrors(err), "validation failed")
	}

	resp, err := h.adminRoleUsecase.CreateRole(c.Context(), usecase.CreateRoleRequest{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return handleError(c, err)
	}
	return response.Created(c, "role created", resp)
}

func (h *RoleHandler) GetRole(c fiber.Ctx) error {
	roleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid role id")
	}

	resp, err := h.adminRoleUsecase.GetRole(c.Context(), roleID)
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "role fetched", resp)
}

func (h *RoleHandler) ListRoles(c fiber.Ctx) error {
	resp, err := h.adminRoleUsecase.ListRoles(c.Context())
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "roles fetched", resp)
}

func (h *RoleHandler) UpdateRole(c fiber.Ctx) error {
	roleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid role id")
	}

	var req updateRoleRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return response.BadRequest(c, formatValidationErrors(err), "validation failed")
	}

	resp, err := h.adminRoleUsecase.UpdateRole(c.Context(), roleID, usecase.UpdateRoleRequest{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "role updated", resp)
}

func (h *RoleHandler) DeleteRole(c fiber.Ctx) error {
	roleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid role id")
	}

	if err := h.adminRoleUsecase.DeleteRole(c.Context(), roleID); err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "role deleted", nil)
}

func (h *RoleHandler) AssignRoleToUser(c fiber.Ctx) error {
	var req assignRoleRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return response.BadRequest(c, formatValidationErrors(err), "validation failed")
	}

	userID, _ := uuid.Parse(req.UserID)
	roleID, _ := uuid.Parse(req.RoleID)

	var grantedBy *uuid.UUID
	if adminIDStr, ok := c.Locals("userID").(string); ok && adminIDStr != "" {
		if adminID, err := uuid.Parse(adminIDStr); err == nil {
			grantedBy = &adminID
		}
	}

	err := h.adminRoleUsecase.AssignRoleToUser(c.Context(), usecase.AssignRoleRequest{
		UserID:    userID,
		RoleID:    roleID,
		GrantedBy: grantedBy,
	})
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "role assigned", nil)
}

func (h *RoleHandler) RemoveRoleFromUser(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid user id")
	}
	roleID, err := uuid.Parse(c.Params("roleId"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid role id")
	}

	if err := h.adminRoleUsecase.RemoveRoleFromUser(c.Context(), userID, roleID); err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "role removed", nil)
}

func (h *RoleHandler) GetUserRoles(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid user id")
	}

	resp, err := h.adminRoleUsecase.GetUserRoles(c.Context(), userID)
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "user roles fetched", resp)
}