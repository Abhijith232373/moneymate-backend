-- name: CreateCampaign :one
INSERT INTO campaigns (
    store_id, name, redeem_code, offer_category, offer_type, reward_value, min_bill_amount, redemption_limit, target_audience, start_date, end_date, status, banner_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
) RETURNING id, store_id, name, redeem_code, offer_category, offer_type, reward_value, min_bill_amount, redemption_limit, target_audience, start_date, end_date, status, banner_url, created_at, updated_at;

-- name: GetCampaignsByStoreID :many
SELECT 
    id, store_id, name, redeem_code, offer_category, offer_type, reward_value, min_bill_amount, redemption_limit, target_audience, start_date, end_date, status, banner_url, created_at, updated_at
FROM campaigns
WHERE store_id = $1
ORDER BY created_at DESC;


-- name: GetCampaignByID :one
SELECT 
    id, store_id, name, redeem_code, offer_category, offer_type, reward_value, min_bill_amount, redemption_limit, target_audience, start_date, end_date, status, banner_url, created_at, updated_at
FROM campaigns
WHERE id = $1 LIMIT 1;

-- name: UpdateCampaignStatus :exec
UPDATE campaigns
SET status = $2, updated_at = NOW()
WHERE id = $1 AND store_id = $3;
