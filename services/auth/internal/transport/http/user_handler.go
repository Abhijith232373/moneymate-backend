package http

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	usecase "github.com/moneymate-2026/moneymate-backend/auth/internal/usecases"
	response "github.com/moneymate-2026/moneymate-backend/shared/pkg/responses"
)

type UserHandler struct {
	adminUserUsecase usecase.AdminUserUsecase
}

func NewUserHandler(adminUserUsecase usecase.AdminUserUsecase) *UserHandler {
	return &UserHandler{adminUserUsecase: adminUserUsecase}
}


func (h *UserHandler) ListUsers(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	sortDesc := c.Query("sort_desc", "true") == "true"

	resp, err := h.adminUserUsecase.ListUsers(c.Context(), usecase.ListUsersRequest{
		Status:   c.Query("status"),
		Search:   c.Query("search"),
		SortBy:   c.Query("sort_by", "created_at"),
		SortDesc: sortDesc,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "users fetched", resp)
}

func (h *UserHandler) GetUser(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid user id")
	}

	resp, err := h.adminUserUsecase.GetUser(c.Context(), userID)
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "user fetched", resp)
}

func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid user id")
	}

	var req updateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return response.BadRequest(c, formatValidationErrors(err), "validation failed")
	}

	resp, err := h.adminUserUsecase.UpdateUser(c.Context(), userID, usecase.UpdateUserRequest{
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "user updated", resp)
}

func (h *UserHandler) UpdateUserStatus(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid user id")
	}

	var req updateUserStatusRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return response.BadRequest(c, formatValidationErrors(err), "validation failed")
	}

	if err := h.adminUserUsecase.UpdateUserStatus(c.Context(), userID, req.Status); err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "user status updated", nil)
}

func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid user id")
	}

	if err := h.adminUserUsecase.DeleteUser(c.Context(), userID); err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "user deleted", nil)
}