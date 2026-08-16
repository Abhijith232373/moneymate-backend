package http

import (
	"github.com/gofiber/fiber/v3"

	"github.com/moneymate-2026/moneymate-backend/services/payment/internal/usecases"
	response "github.com/moneymate-2026/moneymate-backend/shared/pkg/responses"
)

type CategoryHandler struct {
	categories usecases.CategoryUsecase
}

func NewCategoryHandler(categories usecases.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{categories: categories}
}

type createCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

func (h *CategoryHandler) Create(c fiber.Ctx) error {
	userID := userIDFromLocals(c)
	if userID == "" {
		return response.Unauthorized(c, "authentication required")
	}
	var req createCategoryRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request body")
	}
	if err := validate.Struct(req); err != nil {
		return response.BadRequest(c, nil, "validation failed")
	}
	cat, err := h.categories.Create(c.Context(), userID, req.Name)
	if err != nil {
		return handleError(c, err)
	}
	return response.Created(c, "category created", cat)
}

func (h *CategoryHandler) List(c fiber.Ctx) error {
	userID := userIDFromLocals(c)
	if userID == "" {
		return response.Unauthorized(c, "authentication required")
	}
	cats, err := h.categories.List(c.Context(), userID)
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "categories listed", cats)
}

func (h *CategoryHandler) Update(c fiber.Ctx) error {
	userID := userIDFromLocals(c)
	if userID == "" {
		return response.Unauthorized(c, "authentication required")
	}
	var req createCategoryRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request body")
	}
	cat, err := h.categories.Update(c.Context(), c.Params("id"), userID, req.Name)
	if err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "category updated", cat)
}

func (h *CategoryHandler) Delete(c fiber.Ctx) error {
	userID := userIDFromLocals(c)
	if userID == "" {
		return response.Unauthorized(c, "authentication required")
	}
	if err := h.categories.Delete(c.Context(), c.Params("id"), userID); err != nil {
		return handleError(c, err)
	}
	return response.OK(c, "category deleted", nil)
}