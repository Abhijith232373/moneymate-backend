-- name: CreateQRTransaction :one
INSERT INTO qr_transactions (
    store_id, customer_display_id, bill_amount, reward_issued
) VALUES (
    $1, $2, $3, $4
) RETURNING id, store_id, customer_display_id, bill_amount, reward_issued, created_at;

-- name: GetQRTransactionsByStoreID :many
SELECT 
    id, store_id, customer_display_id, bill_amount, reward_issued, created_at
FROM qr_transactions
WHERE store_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetTodayQRScanCount :one
SELECT COUNT(*) 
FROM qr_transactions 
WHERE store_id = $1 AND created_at >= date_trunc('day', NOW());

-- name: GetTodayQRScanVolume :one
SELECT COALESCE(SUM(bill_amount), 0)::numeric
FROM qr_transactions 
WHERE store_id = $1 AND created_at >= date_trunc('day', NOW());
