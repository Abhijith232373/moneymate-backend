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

	storeIDStr := c.Params("store_id")
	if storeIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "store_id is required"})
	}
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid store_id format"})
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid start_date format, expected YYYY-MM-DD"})
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid end_date format, expected YYYY-MM-DD"})
	}

	campaign := &domain.Campaign{
		StoreID:        storeID,
		Name:           req.Name,
		OfferType:      req.OfferType,
		RewardValue:    req.RewardValue,
		MinBillAmount:  req.MinBillAmount,
		StartDate:      startDate,
		EndDate:        endDate,
	}

	if req.BannerURL != "" {
		campaign.BannerURL = &req.BannerURL
	}

	created, err := h.usecase.CreateCampaign(c.Context(), campaign)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(CampaignResponse{
		ID:            created.ID.String(),
		StoreID:       created.StoreID.String(),
		Name:          created.Name,
		OfferType:     created.OfferType,
		RewardValue:   created.RewardValue,
		MinBillAmount: created.MinBillAmount,
		StartDate:     created.StartDate.Format("2006-01-02"),
		EndDate:       created.EndDate.Format("2006-01-02"),
		BannerURL:     getStringOrEmpty(created.BannerURL),
		IsActive:      created.IsActive,
		CreatedAt:     created.CreatedAt.Format(time.RFC3339),
	})
}

func (h *CampaignHandler) GetCampaigns(c fiber.Ctx) error {
	storeIDStr := c.Params("store_id")
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
			ID:            cam.ID.String(),
			StoreID:       cam.StoreID.String(),
			Name:          cam.Name,
			OfferType:     cam.OfferType,
			RewardValue:   cam.RewardValue,
			MinBillAmount: cam.MinBillAmount,
			StartDate:     cam.StartDate.Format("2006-01-02"),
			EndDate:       cam.EndDate.Format("2006-01-02"),
			BannerURL:     getStringOrEmpty(cam.BannerURL),
			IsActive:      cam.IsActive,
			CreatedAt:     cam.CreatedAt.Format(time.RFC3339),
		})
	}

	return c.JSON(response)
}

func getStringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
