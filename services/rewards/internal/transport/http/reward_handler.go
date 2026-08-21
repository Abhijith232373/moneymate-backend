package http

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/moneymate-2026/moneymate-backend/services/rewards/internal/domain"
	"github.com/moneymate-2026/moneymate-backend/shared/pkg/money"
	response "github.com/moneymate-2026/moneymate-backend/shared/pkg/responses"
)

type RewardRuleRequest struct {
	Name                      string `json:"name"`
	MinPercentageBPS          int32  `json:"min_percentage_bps"`
	MaxPercentageBPS          int32  `json:"max_percentage_bps"`
	MinTransactionAmountPaise int64  `json:"min_transaction_amount_paise"`
	MaxPayoutAmountPaise      int64  `json:"max_payout_amount_paise"`
	Active                    bool   `json:"active"`
}

type RewardRuleResponse struct {
	ID                        string `json:"id"`
	Name                      string `json:"name"`
	MinPercentageBPS          int32  `json:"min_percentage_bps"`
	MaxPercentageBPS          int32  `json:"max_percentage_bps"`
	MinTransactionAmountPaise int64  `json:"min_transaction_amount_paise"`
	MaxPayoutAmountPaise      int64  `json:"max_payout_amount_paise"`
	Active                    bool   `json:"active"`
	CreatedAt                 string `json:"created_at"`
	UpdatedAt                 string `json:"updated_at"`
}

type RewardPayoutResponse struct {
	ID                   string  `json:"id"`
	TransactionID        string  `json:"transaction_id"`
	RewardAmountPaise    int64   `json:"reward_amount_paise"`
	RewardAmount         string  `json:"reward_amount"`
	Status               string  `json:"status"`
	PaymentTransactionID *string `json:"payment_transaction_id,omitempty"`
	CreatedAt            string  `json:"created_at"`
}

func (h *RewardHandler) CreateRule(c fiber.Ctx) error {
	var req RewardRuleRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request body")
	}

	rule := domain.RewardRule{
		Name:                      req.Name,
		MinPercentageBPS:          req.MinPercentageBPS,
		MaxPercentageBPS:          req.MaxPercentageBPS,
		MinTransactionAmountPaise: req.MinTransactionAmountPaise,
		MaxPayoutAmountPaise:      req.MaxPayoutAmountPaise,
		Active:                    req.Active,
	}

	created, err := h.rewardUC.CreateRule(c.Context(), rule)
	if err != nil {
		return response.BadRequest(c, nil, err.Error())
	}

	return response.Created(c, "reward rule created", toRuleResponse(created))
}

func (h *RewardHandler) ListRules(c fiber.Ctx) error {
	limit := parseQueryInt(c, "limit", 50)
	offset := parseQueryInt(c, "offset", 0)

	rules, err := h.rewardUC.ListRules(c.Context(), int32(limit), int32(offset))
	if err != nil {
		return response.InternalServerError(c)
	}

	var resp []*RewardRuleResponse
	for _, r := range rules {
		resp = append(resp, toRuleResponse(r))
	}
	return response.OK(c, "reward rules", resp)
}

func (h *RewardHandler) GetRule(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid rule id")
	}

	rule, err := h.rewardUC.GetRule(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "rule not found")
	}

	return response.OK(c, "reward rule", toRuleResponse(rule))
}

func (h *RewardHandler) UpdateRule(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid rule id")
	}

	var req RewardRuleRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, nil, "invalid request body")
	}

	rule := domain.RewardRule{
		ID:                        id,
		Name:                      req.Name,
		MinPercentageBPS:          req.MinPercentageBPS,
		MaxPercentageBPS:          req.MaxPercentageBPS,
		MinTransactionAmountPaise: req.MinTransactionAmountPaise,
		MaxPayoutAmountPaise:      req.MaxPayoutAmountPaise,
		Active:                    req.Active,
	}

	updated, err := h.rewardUC.UpdateRule(c.Context(), rule)
	if err != nil {
		return response.BadRequest(c, nil, err.Error())
	}

	return response.OK(c, "reward rule updated", toRuleResponse(updated))
}

func (h *RewardHandler) DeactivateRule(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid rule id")
	}

	rule, err := h.rewardUC.DeactivateRule(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "rule not found")
	}

	return response.OK(c, "reward rule deactivated", toRuleResponse(rule))
}

func (h *RewardHandler) ListMyPayouts(c fiber.Ctx) error {
	userIDStr := userIDFromLocals(c)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Unauthorized(c, "invalid user id")
	}

	var status *domain.PayoutStatus
	if s := c.Query("status"); s != "" {
		ps := domain.PayoutStatus(s)
		status = &ps
	}

	limit := parseQueryInt(c, "limit", 50)
	offset := parseQueryInt(c, "offset", 0)

	payouts, err := h.rewardUC.ListMyPayouts(c.Context(), userID, status, int32(limit), int32(offset))
	if err != nil {
		return response.InternalServerError(c)
	}

	var resp []*RewardPayoutResponse
	for _, p := range payouts {
		resp = append(resp, toPayoutResponse(p))
	}
	return response.OK(c, "reward payouts", resp)
}

func (h *RewardHandler) ListPayoutsByTransaction(c fiber.Ctx) error {
	txIDStr := c.Query("transaction_id")
	if txIDStr == "" {
		return response.BadRequest(c, nil, "transaction_id query param required")
	}

	txID, err := uuid.Parse(txIDStr)
	if err != nil {
		return response.BadRequest(c, nil, "invalid transaction_id")
	}

	payouts, err := h.rewardUC.ListPayoutsByTransaction(c.Context(), txID)
	if err != nil {
		return response.InternalServerError(c)
	}

	var resp []*RewardPayoutResponse
	for _, p := range payouts {
		resp = append(resp, toPayoutResponse(p))
	}
	return response.OK(c, "reward payouts", resp)
}

func (h *RewardHandler) GetPayoutByID(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, nil, "invalid payout id")
	}

	payout, err := h.rewardUC.GetPayoutByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "payout not found")
	}

	return response.OK(c, "reward payout", toPayoutResponse(payout))
}

func (h *RewardHandler) ReplayFailed(c fiber.Ctx) error {
	if err := h.rewardUC.ReplayFailed(c.Context()); err != nil {
		return response.InternalServerError(c)
	}
	return response.OK(c, "failed payouts replayed", nil)
}

func (h *RewardHandler) FakePaymentEvent(c fiber.Ctx) error {
	raw := c.Body()
	if len(raw) == 0 {
		return response.BadRequest(c, nil, "invalid request body")
	}

	if err := h.rewardUC.ProcessPaymentCompletedEvent(c.Context(), raw); err != nil {
		return response.BadRequest(c, nil, err.Error())
	}
	return response.OK(c, "fake payment event processed", nil)
}

func toRuleResponse(r *domain.RewardRule) *RewardRuleResponse {
	return &RewardRuleResponse{
		ID:                        r.ID.String(),
		Name:                      r.Name,
		MinPercentageBPS:          r.MinPercentageBPS,
		MaxPercentageBPS:          r.MaxPercentageBPS,
		MinTransactionAmountPaise: r.MinTransactionAmountPaise,
		MaxPayoutAmountPaise:      r.MaxPayoutAmountPaise,
		Active:                    r.Active,
		CreatedAt:                 r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:                 r.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toPayoutResponse(p *domain.RewardPayout) *RewardPayoutResponse {
	resp := &RewardPayoutResponse{
		ID:                p.ID.String(),
		TransactionID:     p.TransactionID.String(),
		RewardAmountPaise: p.RewardAmountPaise,
		RewardAmount:      money.FormatPaise(p.RewardAmountPaise),
		Status:            string(p.Status),
		CreatedAt:         p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if p.PaymentTransactionID != nil {
		s := p.PaymentTransactionID.String()
		resp.PaymentTransactionID = &s
	}
	return resp
}

func parseQueryInt(c fiber.Ctx, key string, defaultVal int) int {
	s := c.Query(key)
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}
