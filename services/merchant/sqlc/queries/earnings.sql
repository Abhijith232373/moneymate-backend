-- name: GetEarningsStats :one
SELECT store_id, total_scans, total_earned, updated_at
FROM earnings_stats
WHERE store_id = $1 LIMIT 1;

-- name: UpsertEarningsStats :one
INSERT INTO earnings_stats (
    store_id, total_scans, total_earned
) VALUES (
    $1, $2, $3
) ON CONFLICT (store_id) DO UPDATE SET
    total_scans = EXCLUDED.total_scans,
    total_earned = EXCLUDED.total_earned,
    updated_at = NOW()
RETURNING store_id, total_scans, total_earned, updated_at;

-- name: GetRequestedMilestones :many
SELECT milestone_scans
FROM earnings_payout_requests
WHERE store_id = $1;

-- name: CreatePayoutRequest :one
INSERT INTO earnings_payout_requests (
    store_id, milestone_scans, reward_amount, status
) VALUES (
    $1, $2, $3, 'requested'
) RETURNING id, store_id, milestone_scans, reward_amount, status, created_at;
