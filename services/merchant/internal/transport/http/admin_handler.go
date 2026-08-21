package http

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/domain"
	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/usecases"
)

// AdminHandler handles REST requests from platform administrators for managing merchants across all modules.
type AdminHandler struct {
	usecase *usecases.AdminUseCase
}

func NewAdminHandler(uc *usecases.AdminUseCase) *AdminHandler {
	return &AdminHandler{usecase: uc}
}

// Stores
func (h *AdminHandler) GetAllStores(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	stores, err := h.usecase.GetAllStores(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": stores})
}

func (h *AdminHandler) GetStoreByID(c fiber.Ctx) error {
	storeID := c.Params("id")
	store, err := h.usecase.GetStoreByID(c.Context(), storeID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "store not found or invalid id"})
	}
	return c.JSON(fiber.Map{"success": true, "data": store})
}

func (h *AdminHandler) UpdateStoreStatus(c fiber.Ctx) error {
	storeID := c.Params("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.Bind().Body(&req); err != nil || req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "status is required"})
	}

	if err := h.usecase.UpdateStoreStatus(c.Context(), storeID, req.Status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "store status updated successfully"})
}

func (h *AdminHandler) DeleteStore(c fiber.Ctx) error {
	storeID := c.Params("id")
	if err := h.usecase.DeleteStore(c.Context(), storeID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "store deleted successfully"})
}

// Campaigns
func (h *AdminHandler) GetAllCampaigns(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	campaigns, err := h.usecase.GetAllCampaigns(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": campaigns})
}

func (h *AdminHandler) CreateCampaign(c fiber.Ctx) error {
	var req CreateCampaignRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	storeIDStr := c.Params("store_id")
	if storeIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "store_id is required in params"})
	}
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid store_id format"})
	}

	parseDate := func(dateStr string) (time.Time, error) {
		if t, err := time.Parse("2006-01-02T15:04:05", dateStr); err == nil {
			return t, nil
		}
		if t, err := time.Parse("2006-01-02T15:04", dateStr); err == nil {
			return t, nil
		}
		if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
			return t, nil
		}
		return time.Parse("2006-01-02", dateStr)
	}

	startDate, err := parseDate(req.StartDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid start_date format"})
	}

	endDate, err := parseDate(req.EndDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid end_date format"})
	}

	campaign := &domain.Campaign{
		StoreID:         storeID,
		Name:            req.Name,
		OfferType:       req.OfferType,
		RewardValue:     req.RewardValue,
		MinBillAmount:   req.MinBillAmount,
		RedemptionLimit: req.RedemptionLimit,
		TargetAudience:  req.TargetAudience,
		StartDate:       startDate,
		EndDate:         endDate,
	}

	if req.RedeemCode != "" {
		campaign.RedeemCode = &req.RedeemCode
	}
	if req.OfferCategory != "" {
		campaign.OfferCategory = &req.OfferCategory
	}
	if req.BannerURL != "" {
		campaign.BannerURL = &req.BannerURL
	}

	created, err := h.usecase.CreateCampaign(c.Context(), campaign)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": CampaignResponse{
			ID:              created.ID.String(),
			StoreID:         created.StoreID.String(),
			Name:            created.Name,
			RedeemCode:      getStringOrEmpty(created.RedeemCode),
			OfferCategory:   getStringOrEmpty(created.OfferCategory),
			OfferType:       created.OfferType,
			RewardValue:     created.RewardValue,
			MinBillAmount:   created.MinBillAmount,
			RedemptionLimit: created.RedemptionLimit,
			StartDate:       created.StartDate.Format(time.RFC3339),
			EndDate:         created.EndDate.Format(time.RFC3339),
			BannerURL:       getStringOrEmpty(created.BannerURL),
			Status:          created.Status,
			CreatedAt:       created.CreatedAt.Format(time.RFC3339),
		},
	})
}

func (h *AdminHandler) GetCampaignsByStoreID(c fiber.Ctx) error {
	storeID := c.Params("store_id")
	campaigns, err := h.usecase.GetCampaignsByStoreID(c.Context(), storeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": campaigns})
}

func (h *AdminHandler) UpdateCampaignStatus(c fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid payload"})
	}

	if err := h.usecase.UpdateCampaignStatus(c.Context(), id, req.IsActive); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "campaign status updated successfully"})
}

func (h *AdminHandler) DeleteCampaign(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.usecase.DeleteCampaign(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "campaign deleted successfully"})
}

// KYC Verification
func (h *AdminHandler) GetAllKYCDocuments(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	docs, err := h.usecase.GetAllKYCDocuments(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": docs})
}

func (h *AdminHandler) GetKYCByStoreID(c fiber.Ctx) error {
	storeID := c.Params("store_id")
	doc, err := h.usecase.GetKYCByStoreID(c.Context(), storeID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "kyc document not found"})
	}
	return c.JSON(fiber.Map{"success": true, "data": doc})
}

func (h *AdminHandler) VerifyKYCDocument(c fiber.Ctx) error {
	storeID := c.Params("store_id")
	var req struct {
		IsVerified bool   `json:"is_verified"`
		Status     string `json:"status"` // e.g. "active", "rejected"
	}
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid payload"})
	}
	if req.Status == "" {
		if req.IsVerified {
			req.Status = "active"
		} else {
			req.Status = "rejected"
		}
	}

	doc, err := h.usecase.VerifyKYCDocument(c.Context(), storeID, req.IsVerified, req.Status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "kyc verification status updated", "data": doc})
}

// Rewards
func (h *AdminHandler) GetAllRewardTransactions(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	txs, err := h.usecase.GetAllRewardTransactions(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": txs})
}

func (h *AdminHandler) GetPlatformRewardSummary(c fiber.Ctx) error {
	summary, err := h.usecase.GetPlatformRewardSummary(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": summary})
}

// Subscriptions
func (h *AdminHandler) GetAllSubscriptions(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	subs, err := h.usecase.GetAllSubscriptions(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": subs})
}

func (h *AdminHandler) UpdateStoreSubscriptionPlan(c fiber.Ctx) error {
	storeID := c.Params("store_id")
	var req struct {
		PlanCode string `json:"plan_code"`
	}
	if err := c.Bind().Body(&req); err != nil || req.PlanCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "plan_code is required"})
	}

	sub, err := h.usecase.UpdateStoreSubscriptionPlan(c.Context(), storeID, req.PlanCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "subscription plan updated successfully", "data": sub})
}

// Dashboard
func (h *AdminHandler) GetAdminDashboardStats(c fiber.Ctx) error {
	stats, err := h.usecase.GetAdminDashboardStats(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": stats})
}

