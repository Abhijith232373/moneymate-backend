-- name: CreateRewardRule :one
INSERT INTO reward_rules (
    name, min_percentage_bps, max_percentage_bps,
    min_transaction_amount_paise, max_payout_amount_paise,
    active, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: ListRewardRules :many
SELECT * FROM reward_rules
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetRewardRuleByID :one
SELECT * FROM reward_rules
WHERE id = $1;

-- name: GetActiveRewardRule :one
SELECT * FROM reward_rules
WHERE active = TRUE
ORDER BY created_at DESC
LIMIT 1;

-- name: UpdateRewardRule :one
UPDATE reward_rules
SET name = $2,
    min_percentage_bps = $3,
    max_percentage_bps = $4,
    min_transaction_amount_paise = $5,
    max_payout_amount_paise = $6,
    active = $7,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeactivateRewardRule :one
UPDATE reward_rules
SET active = FALSE,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
