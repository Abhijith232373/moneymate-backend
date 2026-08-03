package http

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/services/merchant/internal/usecases"
)

// SubscriptionHandler exposes REST endpoints for viewing pricing tiers, checking active billing standing, and changing plans.
type SubscriptionHandler struct {
	// usecase orchestrates pricing catalog retrieval, active subscription lookups, and compliance-validated tier changes.
	usecase usecases.SubscriptionUseCase
}

// NewSubscriptionHandler initializes and returns a new SubscriptionHandler instance with injected dependencies.
func NewSubscriptionHandler(uc usecases.SubscriptionUseCase) *SubscriptionHandler {
	return &SubscriptionHandler{usecase: uc}
}

// GetSubscriptionPlans handles GET /merchant/:store_id/subscriptions/plans and returns the catalog of pricing tiers.
func (h *SubscriptionHandler) GetSubscriptionPlans(c fiber.Ctx) error {
	storeIDStr := resolveMerchantID(c)
	if storeIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "store_id parameter is required"})
	}
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid store_id UUID format"})
	}

	plans, err := h.usecase.GetPlans(c.Context(), storeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var response []SubscriptionPlanResponse
	for _, p := range plans {
		formattedPrice := fmt.Sprintf("$%.0f", p.Price)
		if p.Price != float64(int64(p.Price)) {
			formattedPrice = fmt.Sprintf("$%.2f", p.Price)
		}

		formattedCycle := "/mo"
		if strings.ToLower(p.BillingCycle) == "annual" {
			formattedCycle = "/yr"
		}

		var featuresResp []PlanFeatureResponse
		for _, f := range p.Features {
			featuresResp = append(featuresResp, PlanFeatureResponse{
				Name:     f.Name,
				Included: f.Included,
			})
		}
		if featuresResp == nil {
			featuresResp = []PlanFeatureResponse{}
		}

		response = append(response, SubscriptionPlanResponse{
			ID:                 p.ID.String(),
			PlanCode:           p.PlanCode,
			Name:               p.Name,
			Price:              p.Price,
			FormattedPrice:     formattedPrice,
			BillingCycle:       p.BillingCycle,
			FormattedCycle:     formattedCycle,
			Description:        p.Description,
			MaxActiveCampaigns: p.MaxActiveCampaigns,
			IsMostPopular:      p.IsMostPopular,
			Features:           featuresResp,
			IsCurrent:          p.IsCurrent,
		})
	}

	if response == nil {
		response = []SubscriptionPlanResponse{}
	}

	return c.JSON(response)
}

// GetCurrentSubscription handles GET /merchant/:store_id/subscriptions/current and returns the store's active billing timeline.
func (h *SubscriptionHandler) GetCurrentSubscription(c fiber.Ctx) error {
	storeIDStr := resolveMerchantID(c)
	if storeIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "store_id parameter is required"})
	}
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid store_id UUID format"})
	}

	sub, err := h.usecase.GetCurrentSubscription(c.Context(), storeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	planName := strings.Title(sub.PlanCode)
	formattedRenewal := sub.CurrentPeriodEnd.Format("January 02, 2006")

	return c.JSON(CurrentSubscriptionResponse{
		ID:                   sub.ID.String(),
		StoreID:              sub.StoreID.String(),
		PlanCode:             sub.PlanCode,
		PlanName:             planName,
		Status:               sub.Status,
		BillingCycle:         sub.BillingCycle,
		CurrentPeriodStart:   sub.CurrentPeriodStart.Format(time.RFC3339),
		CurrentPeriodEnd:     sub.CurrentPeriodEnd.Format(time.RFC3339),
		FormattedRenewalDate: formattedRenewal,
		AutoRenew:            sub.AutoRenew,
	})
}

// ChangeSubscriptionPlan handles POST /merchant/:store_id/subscriptions/change to upgrade or downgrade a store's tier.
func (h *SubscriptionHandler) ChangeSubscriptionPlan(c fiber.Ctx) error {
	storeIDStr := resolveMerchantID(c)
	if storeIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "store_id parameter is required"})
	}
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid store_id UUID format"})
	}

	var req ChangePlanRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON request payload"})
	}
	if strings.TrimSpace(req.PlanCode) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "plan_code is required"})
	}

	sub, err := h.usecase.ChangePlan(c.Context(), storeID, req.PlanCode)
	if err != nil {
		// Return 400 Bad Request for business rule violations (e.g. campaign limits exceeded)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	planName := strings.Title(sub.PlanCode)
	msg := fmt.Sprintf("Successfully updated your subscription plan to %s.", planName)

	return c.Status(fiber.StatusOK).JSON(ChangePlanResponse{
		ID:       sub.ID.String(),
		StoreID:  sub.StoreID.String(),
		PlanCode: sub.PlanCode,
		PlanName: planName,
		Status:   sub.Status,
		Message:  msg,
	})
}

// InitiateUpgrade handles POST /merchant/:store_id/subscriptions/upgrade/initiate
func (h *SubscriptionHandler) InitiateUpgrade(c fiber.Ctx) error {
	storeIDStr := resolveMerchantID(c)
	if storeIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "store_id parameter is required"})
	}
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid store_id UUID format"})
	}

	var req struct {
		PlanCode string `json:"plan_code"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON request payload"})
	}

	orderID, err := h.usecase.CreateUpgradeOrder(c.Context(), storeID, req.PlanCode)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"order_id": orderID,
	})
}

// VerifyUpgrade handles POST /merchant/:store_id/subscriptions/upgrade/verify
func (h *SubscriptionHandler) VerifyUpgrade(c fiber.Ctx) error {
	storeIDStr := resolveMerchantID(c)
	if storeIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "store_id parameter is required"})
	}
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid store_id UUID format"})
	}

	var req struct {
		PlanCode  string `json:"plan_code"`
		OrderID   string `json:"razorpay_order_id"`
		PaymentID string `json:"razorpay_payment_id"`
		Signature string `json:"razorpay_signature"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON request payload"})
	}

	sub, err := h.usecase.VerifyUpgrade(c.Context(), storeID, req.PlanCode, req.OrderID, req.PaymentID, req.Signature)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Subscription upgraded successfully",
		"plan":    sub.PlanCode,
	})
}
