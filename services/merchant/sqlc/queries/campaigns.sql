-- name: CreateCampaign :one
INSERT INTO campaigns (
    store_id, name, offer_type, reward_value, min_bill_amount, target_audience, start_date, end_date, is_active, banner_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING id, store_id, name, offer_type, reward_value, min_bill_amount, target_audience, start_date, end_date, is_active, banner_url, created_at, updated_at;

-- name: GetCampaignsByStoreID :many
SELECT 
    id, store_id, name, offer_type, reward_value, min_bill_amount, target_audience, start_date, end_date, is_active, banner_url, created_at, updated_at
FROM campaigns
WHERE store_id = $1
ORDER BY created_at DESC;


-- name: GetCampaignByID :one
SELECT 
    id, store_id, name, offer_type, reward_value, min_bill_amount, target_audience, start_date, end_date, is_active, banner_url, created_at, updated_at
FROM campaigns
WHERE id = $1 LIMIT 1;

-- name: UpdateCampaignStatus :exec
UPDATE campaigns
SET is_active = $2, updated_at = NOW()
WHERE id = $1 AND store_id = $3;
