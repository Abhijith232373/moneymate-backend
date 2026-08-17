package http

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	usecase "github.com/moneymate-2026/moneymate-backend/auth/internal/usecases"
	response "github.com/moneymate-2026/moneymate-backend/shared/pkg/responses"
)

type StaffHandler struct {
	staffUsecase usecase.StaffUsecase
}

func NewStaffHandler(staffUsecase usecase.StaffUsecase) *StaffHandler {
	return &StaffHandler{staffUsecase: staffUsecase}
}

func (h *StaffHandler) CreateStaff(c fiber.Ctx) error {
	var req createUserRequest // Reusing user creation request dto
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return response.BadRequest(c, formatValidationErrors(err), "validation failed")
	}

	resp, err := h.staffUsecase.CreateStaff(c.Context(), usecase.CreateUserRequest{
		Email:    req.Email,
		Phone:    req.Phone,
		FullName: req.FullName,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		return handleError(c, err)
	}
	return response.Created(c, "staff created", resp)
}

func (h *StaffHandler) ListStaff(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	sortDesc := c.Query("sort_desc", "true") == "true"

	resp, err := h.staffUsecase.ListStaff(c.Context(), usecase.ListUsersRequest{
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
	return response.OK(c, "staff fetched", resp)
}

func (h *StaffHandler) GetStaff(c fiber.Ctx) error {
	staffID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid staff id")
	}

	resp, err := h.staffUsecase.GetStaff(c.Context(), staffID)
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "staff fetched", resp)
}

func (h *StaffHandler) UpdateStaff(c fiber.Ctx) error {
	staffID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid staff id")
	}

	var req updateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return response.BadRequest(c, formatValidationErrors(err), "validation failed")
	}

	resp, err := h.staffUsecase.UpdateStaff(c.Context(), staffID, usecase.UpdateUserRequest{
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "staff updated", resp)
}

func (h *StaffHandler) UpdateStaffStatus(c fiber.Ctx) error {
	staffID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid staff id")
	}

	var req updateUserStatusRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return response.BadRequest(c, formatValidationErrors(err), "validation failed")
	}

	if err := h.staffUsecase.UpdateStaffStatus(c.Context(), staffID, req.Status); err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "staff status updated", nil)
}

func (h *StaffHandler) DeleteStaff(c fiber.Ctx) error {
	staffID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid staff id")
	}

	if err := h.staffUsecase.DeleteStaff(c.Context(), staffID); err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "staff deleted", nil)
}
