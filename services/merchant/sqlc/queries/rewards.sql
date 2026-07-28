-- name: CreateRewardBalance :one
-- CreateRewardBalance initializes the financial and scan stats record for a newly registered store.
INSERT INTO reward_balances (
    store_id, available_balance, total_scans, premium_points, weekly_growth_percentage
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id, store_id, available_balance, total_scans, premium_points, weekly_growth_percentage, created_at, updated_at;

-- name: GetRewardBalanceByStoreID :one
-- GetRewardBalanceByStoreID fetches the current rewards center summary metrics for a specific merchant store.
SELECT 
    id, store_id, available_balance, total_scans, premium_points, weekly_growth_percentage, created_at, updated_at
FROM reward_balances
WHERE store_id = $1 LIMIT 1;

-- name: UpdateRewardBalance :exec
-- UpdateRewardBalance modifies the balance, scan counts, and growth stats in a high-concurrency safe manner.
UPDATE reward_balances
SET available_balance = $2, total_scans = $3, premium_points = $4, weekly_growth_percentage = $5, updated_at = NOW()
WHERE store_id = $1;

-- name: DeductRewardBalance :exec
-- DeductRewardBalance atomically subtracts redeemed funds from the available balance if sufficient funds exist.
UPDATE reward_balances
SET available_balance = available_balance - $2, updated_at = NOW()
WHERE store_id = $1 AND available_balance >= $2;

-- name: CreateRewardTransaction :one
-- CreateRewardTransaction inserts an immutable ledger entry for a reward earning or redemption event.
INSERT INTO reward_transactions (
    store_id, campaign_name, display_id, status, amount, transaction_type
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING id, store_id, campaign_name, display_id, status, amount, transaction_type, created_at;

-- name: GetRewardTransactionsByStoreID :many
-- GetRewardTransactionsByStoreID retrieves paginated transaction history with optional time-range and keyword search filtering.
SELECT 
    id, store_id, campaign_name, display_id, status, amount, transaction_type, created_at
FROM reward_transactions
WHERE store_id = sqlc.arg('store_id')
  AND (sqlc.arg('search_query')::text = '' OR campaign_name ILIKE '%' || sqlc.arg('search_query') || '%' OR display_id ILIKE '%' || sqlc.arg('search_query') || '%')
  AND (sqlc.narg('created_after')::timestamptz IS NULL OR created_at >= sqlc.narg('created_after'))
ORDER BY created_at DESC
LIMIT sqlc.arg('limit_count') OFFSET sqlc.arg('offset_count');


-- name: CreateRedemptionRequest :one
-- CreateRedemptionRequest logs a merchant's withdrawal request to their bank account.
INSERT INTO redemption_requests (
    store_id, amount, bank_transfer_authorized, status, reference_id
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id, store_id, amount, bank_transfer_authorized, status, reference_id, created_at, updated_at;
