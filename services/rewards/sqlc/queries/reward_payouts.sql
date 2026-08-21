-- name: InsertRewardPayout :one
INSERT INTO reward_payouts (
    transaction_id, recipient_id, recipient_account_id, recipient_type,
    rule_id, original_amount_paise, reward_percentage_bps,
    reward_amount_paise, status, event_payload
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, 'pending', $9
) RETURNING *;

-- name: GetRewardPayoutByID :one
SELECT * FROM reward_payouts
WHERE id = $1;

-- name: GetRewardPayoutByOriginalTransaction :many
SELECT * FROM reward_payouts
WHERE transaction_id = $1
ORDER BY created_at DESC;

-- name: ListRewardPayoutsByRecipient :many
SELECT * FROM reward_payouts
WHERE recipient_id = $1
  AND (sqlc.narg('status')::reward_payout_status IS NULL OR status = sqlc.narg('status')::reward_payout_status)
ORDER BY created_at DESC
LIMIT sqlc.arg('limit_count') OFFSET sqlc.arg('offset_count');

-- name: MarkRewardPayoutCompleted :one
UPDATE reward_payouts
SET status = 'completed',
    payment_transaction_id = $2,
    completed_at = NOW(),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: MarkRewardPayoutFailed :one
UPDATE reward_payouts
SET status = 'failed',
    failure_reason = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
