package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moneymate-2026/moneymate-backend/services/rewards/internal/domain"
)

type RewardRepo struct {
	pool *pgxpool.Pool
}

func NewRewardRepo(pool *pgxpool.Pool) *RewardRepo {
	return &RewardRepo{pool: pool}
}

func (r *RewardRepo) CreateRule(ctx context.Context, rule domain.RewardRule) (*domain.RewardRule, error) {
	var createdBy interface{} = rule.CreatedBy
	if rule.CreatedBy == nil {
		createdBy = nil
	}
	row := r.pool.QueryRow(ctx, `
		INSERT INTO reward_rules (
			name, min_percentage_bps, max_percentage_bps,
			min_transaction_amount_paise, max_payout_amount_paise,
			active, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, min_percentage_bps, max_percentage_bps,
			min_transaction_amount_paise, max_payout_amount_paise,
			active, created_by, created_at, updated_at
	`, rule.Name, rule.MinPercentageBPS, rule.MaxPercentageBPS,
		rule.MinTransactionAmountPaise, rule.MaxPayoutAmountPaise,
		rule.Active, createdBy)

	return scanRule(row)
}

func (r *RewardRepo) ListRules(ctx context.Context, limit, offset int32) ([]*domain.RewardRule, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, min_percentage_bps, max_percentage_bps,
			min_transaction_amount_paise, max_payout_amount_paise,
			active, created_by, created_at, updated_at
		FROM reward_rules
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list rules: %w", err)
	}
	defer rows.Close()

	var rules []*domain.RewardRule
	for rows.Next() {
		rule, err := scanRule(rows)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, rows.Err()
}

func (r *RewardRepo) GetRule(ctx context.Context, id uuid.UUID) (*domain.RewardRule, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, min_percentage_bps, max_percentage_bps,
			min_transaction_amount_paise, max_payout_amount_paise,
			active, created_by, created_at, updated_at
		FROM reward_rules WHERE id = $1
	`, id)
	return scanRule(row)
}

func (r *RewardRepo) GetActiveRule(ctx context.Context) (*domain.RewardRule, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, min_percentage_bps, max_percentage_bps,
			min_transaction_amount_paise, max_payout_amount_paise,
			active, created_by, created_at, updated_at
		FROM reward_rules
		WHERE active = TRUE
		ORDER BY created_at DESC
		LIMIT 1
	`)
	return scanRule(row)
}

func (r *RewardRepo) UpdateRule(ctx context.Context, rule domain.RewardRule) (*domain.RewardRule, error) {
	row := r.pool.QueryRow(ctx, `
		UPDATE reward_rules
		SET name = $2, min_percentage_bps = $3, max_percentage_bps = $4,
			min_transaction_amount_paise = $5, max_payout_amount_paise = $6,
			active = $7, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, min_percentage_bps, max_percentage_bps,
			min_transaction_amount_paise, max_payout_amount_paise,
			active, created_by, created_at, updated_at
	`, rule.ID, rule.Name, rule.MinPercentageBPS, rule.MaxPercentageBPS,
		rule.MinTransactionAmountPaise, rule.MaxPayoutAmountPaise, rule.Active)
	return scanRule(row)
}

func (r *RewardRepo) DeactivateRule(ctx context.Context, id uuid.UUID) (*domain.RewardRule, error) {
	row := r.pool.QueryRow(ctx, `
		UPDATE reward_rules
		SET active = FALSE, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, min_percentage_bps, max_percentage_bps,
			min_transaction_amount_paise, max_payout_amount_paise,
			active, created_by, created_at, updated_at
	`, id)
	return scanRule(row)
}

func (r *RewardRepo) InsertPayout(ctx context.Context, payout domain.RewardPayout) (*domain.RewardPayout, error) {
	var ruleID interface{} = payout.RuleID
	if payout.RuleID == nil {
		ruleID = nil
	}

	eventPayload := payout.EventPayload
	if eventPayload == nil {
		eventPayload = []byte("{}")
	}

	row := r.pool.QueryRow(ctx, `
		INSERT INTO reward_payouts (
			transaction_id, recipient_id, recipient_account_id, recipient_type,
			rule_id, original_amount_paise, reward_percentage_bps,
			reward_amount_paise, status, event_payload
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'pending', $9)
		RETURNING id, transaction_id, recipient_id, recipient_account_id, recipient_type,
			rule_id, original_amount_paise, reward_percentage_bps,
			reward_amount_paise, status, payment_transaction_id, failure_reason,
			event_payload, created_at, updated_at, completed_at
	`, payout.TransactionID, payout.RecipientID, payout.RecipientAccountID,
		payout.RecipientType, ruleID, payout.OriginalAmountPaise,
		payout.RewardPercentageBPS, payout.RewardAmountPaise, eventPayload)
	return scanPayout(row)
}

func (r *RewardRepo) GetPayoutByID(ctx context.Context, id uuid.UUID) (*domain.RewardPayout, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, transaction_id, recipient_id, recipient_account_id, recipient_type,
			rule_id, original_amount_paise, reward_percentage_bps,
			reward_amount_paise, status, payment_transaction_id, failure_reason,
			event_payload, created_at, updated_at, completed_at
		FROM reward_payouts WHERE id = $1
	`, id)
	return scanPayout(row)
}

func (r *RewardRepo) ListPayoutsByRecipient(ctx context.Context, recipientID uuid.UUID, status *domain.PayoutStatus, limit, offset int32) ([]*domain.RewardPayout, error) {
	var statusVal interface{} = nil
	if status != nil {
		statusVal = string(*status)
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, transaction_id, recipient_id, recipient_account_id, recipient_type,
			rule_id, original_amount_paise, reward_percentage_bps,
			reward_amount_paise, status, payment_transaction_id, failure_reason,
			event_payload, created_at, updated_at, completed_at
		FROM reward_payouts
		WHERE recipient_id = $1
		  AND ($2::reward_payout_status IS NULL OR status = $2::reward_payout_status)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`, recipientID, statusVal, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list payouts by recipient: %w", err)
	}
	defer rows.Close()

	var payouts []*domain.RewardPayout
	for rows.Next() {
		p, err := scanPayout(rows)
		if err != nil {
			return nil, err
		}
		payouts = append(payouts, p)
	}
	return payouts, rows.Err()
}

func (r *RewardRepo) ListPayoutsByTransaction(ctx context.Context, transactionID uuid.UUID) ([]*domain.RewardPayout, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, transaction_id, recipient_id, recipient_account_id, recipient_type,
			rule_id, original_amount_paise, reward_percentage_bps,
			reward_amount_paise, status, payment_transaction_id, failure_reason,
			event_payload, created_at, updated_at, completed_at
		FROM reward_payouts
		WHERE transaction_id = $1
		ORDER BY created_at DESC
	`, transactionID)
	if err != nil {
		return nil, fmt.Errorf("list payouts by transaction: %w", err)
	}
	defer rows.Close()

	var payouts []*domain.RewardPayout
	for rows.Next() {
		p, err := scanPayout(rows)
		if err != nil {
			return nil, err
		}
		payouts = append(payouts, p)
	}
	return payouts, rows.Err()
}

func (r *RewardRepo) MarkCompleted(ctx context.Context, payoutID, paymentTransactionID uuid.UUID) (*domain.RewardPayout, error) {
	row := r.pool.QueryRow(ctx, `
		UPDATE reward_payouts
		SET status = 'completed',
			payment_transaction_id = $2,
			completed_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, transaction_id, recipient_id, recipient_account_id, recipient_type,
			rule_id, original_amount_paise, reward_percentage_bps,
			reward_amount_paise, status, payment_transaction_id, failure_reason,
			event_payload, created_at, updated_at, completed_at
	`, payoutID, paymentTransactionID)
	return scanPayout(row)
}

func (r *RewardRepo) MarkFailed(ctx context.Context, payoutID uuid.UUID, reason string) (*domain.RewardPayout, error) {
	row := r.pool.QueryRow(ctx, `
		UPDATE reward_payouts
		SET status = 'failed',
			failure_reason = $2,
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, transaction_id, recipient_id, recipient_account_id, recipient_type,
			rule_id, original_amount_paise, reward_percentage_bps,
			reward_amount_paise, status, payment_transaction_id, failure_reason,
			event_payload, created_at, updated_at, completed_at
	`, payoutID, reason)
	return scanPayout(row)
}

func (r *RewardRepo) ListFailedPayouts(ctx context.Context, limit int32) ([]*domain.RewardPayout, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, transaction_id, recipient_id, recipient_account_id, recipient_type,
			rule_id, original_amount_paise, reward_percentage_bps,
			reward_amount_paise, status, payment_transaction_id, failure_reason,
			event_payload, created_at, updated_at, completed_at
		FROM reward_payouts
		WHERE status = 'failed'
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("list failed payouts: %w", err)
	}
	defer rows.Close()

	var payouts []*domain.RewardPayout
	for rows.Next() {
		p, err := scanPayout(rows)
		if err != nil {
			return nil, err
		}
		payouts = append(payouts, p)
	}
	return payouts, rows.Err()
}

type scannable interface {
	Scan(dest ...any) error
}

func scanRule(row scannable) (*domain.RewardRule, error) {
	var r domain.RewardRule
	err := row.Scan(
		&r.ID, &r.Name, &r.MinPercentageBPS, &r.MaxPercentageBPS,
		&r.MinTransactionAmountPaise, &r.MaxPayoutAmountPaise,
		&r.Active, &r.CreatedBy, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("not found")
		}
		return nil, fmt.Errorf("scan rule: %w", err)
	}
	return &r, nil
}

func scanPayout(row scannable) (*domain.RewardPayout, error) {
	var p domain.RewardPayout
	var eventPayload []byte
	err := row.Scan(
		&p.ID, &p.TransactionID, &p.RecipientID, &p.RecipientAccountID,
		&p.RecipientType, &p.RuleID, &p.OriginalAmountPaise,
		&p.RewardPercentageBPS, &p.RewardAmountPaise, &p.Status,
		&p.PaymentTransactionID, &p.FailureReason, &eventPayload,
		&p.CreatedAt, &p.UpdatedAt, &p.CompletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("not found")
		}
		return nil, fmt.Errorf("scan payout: %w", err)
	}
	if eventPayload != nil {
		p.EventPayload = eventPayload
	}
	return &p, nil
}

// Suppress unused import
var _ = json.Marshal
var _ = time.Now
