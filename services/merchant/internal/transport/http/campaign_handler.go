package http

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/domain"
	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/usecases"
)

type CampaignHandler struct {
	usecase usecases.CampaignUseCase
}

func NewCampaignHandler(uc usecases.CampaignUseCase) *CampaignHandler {
	return &CampaignHandler{usecase: uc}
}

func (h *CampaignHandler) CreateCampaign(c fiber.Ctx) error {
	var req CreateCampaignRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	storeIDStr := resolveMerchantID(c)
	if storeIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "store_id is required"})
	}
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid store_id format"})
	}

	// Try parsing with seconds first, then without, then as a simple date
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
		if err.Error() == "active campaign limit reached for Growth plan (max 5)" || err.Error() == "active campaign limit reached for Essential plan (max 1)" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(CampaignResponse{
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
	})
}

func (h *CampaignHandler) GetCampaigns(c fiber.Ctx) error {
	storeIDStr := resolveMerchantID(c)
	if storeIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "store_id is required"})
	}
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid store_id format"})
	}

	campaigns, err := h.usecase.GetCampaignsByStoreID(c.Context(), storeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var response []CampaignResponse
	for _, cam := range campaigns {
		response = append(response, CampaignResponse{
			ID:              cam.ID.String(),
			StoreID:         cam.StoreID.String(),
			Name:            cam.Name,
			RedeemCode:      getStringOrEmpty(cam.RedeemCode),
			OfferCategory:   getStringOrEmpty(cam.OfferCategory),
			OfferType:       cam.OfferType,
			RewardValue:     cam.RewardValue,
			MinBillAmount:   cam.MinBillAmount,
			RedemptionLimit: cam.RedemptionLimit,
			StartDate:       cam.StartDate.Format(time.RFC3339),
			EndDate:         cam.EndDate.Format(time.RFC3339),
			BannerURL:       getStringOrEmpty(cam.BannerURL),
			Status:          cam.Status,
			CreatedAt:       cam.CreatedAt.Format(time.RFC3339),
		})
	}

	return c.JSON(response)
}

func (h *CampaignHandler) GetAllPublicCampaigns(c fiber.Ctx) error {
	limit := 50
	offset := 0
	campaigns, err := h.usecase.GetPublicCampaigns(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var response []CampaignResponse
	for _, cam := range campaigns {
		response = append(response, CampaignResponse{
			ID:              cam.ID.String(),
			StoreID:         cam.StoreID.String(),
			Name:            cam.Name,
			RedeemCode:      getStringOrEmpty(cam.RedeemCode),
			OfferCategory:   getStringOrEmpty(cam.OfferCategory),
			OfferType:       cam.OfferType,
			RewardValue:     cam.RewardValue,
			MinBillAmount:   cam.MinBillAmount,
			RedemptionLimit: cam.RedemptionLimit,
			StartDate:       cam.StartDate.Format(time.RFC3339),
			EndDate:         cam.EndDate.Format(time.RFC3339),
			BannerURL:       getStringOrEmpty(cam.BannerURL),
			Status:          cam.Status,
			CreatedAt:       cam.CreatedAt.Format(time.RFC3339),
		})
	}

	return c.JSON(response)
}

func (h *CampaignHandler) UpdateCampaignStatus(c fiber.Ctx) error {
	storeIDStr := resolveMerchantID(c)
	if storeIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "store_id is required"})
	}
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid store_id format"})
	}

	campaignIDStr := c.Params("id")
	if campaignIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "campaign id is required"})
	}
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid campaign id format"})
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	status := "paused"
	if req.IsActive {
		status = "active"
	}

	err = h.usecase.UpdateCampaignStatus(c.Context(), campaignID, storeID, status)
	if err != nil {
		if err.Error() == "active campaign limit reached for Growth plan (max 5)" || err.Error() == "active campaign limit reached for Essential plan (max 1)" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Campaign status updated successfully"})
}

func getStringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
