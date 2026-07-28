-- name: CreateCampaign :one
INSERT INTO campaigns (
    store_id, name, offer_type, reward_value, min_bill_amount, target_audience, start_date, end_date, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING id, store_id, name, offer_type, reward_value, min_bill_amount, target_audience, start_date, end_date, is_active, created_at, updated_at;

-- name: GetCampaignsByStoreID :many
SELECT 
    id, store_id, name, offer_type, reward_value, min_bill_amount, target_audience, start_date, end_date, is_active, created_at, updated_at
FROM campaigns
WHERE store_id = $1
ORDER BY created_at DESC;

-- name: GetCampaignsByOwnerID :many
SELECT 
    c.id, c.store_id, c.name, c.offer_type, c.reward_value, c.min_bill_amount, c.target_audience, c.start_date, c.end_date, c.is_active, c.created_at, c.updated_at
FROM campaigns c
JOIN stores s ON c.store_id = s.id
WHERE s.owner_id = $1
ORDER BY c.created_at DESC;

-- name: GetCampaignByID :one
SELECT 
    id, store_id, name, offer_type, reward_value, min_bill_amount, target_audience, start_date, end_date, is_active, created_at, updated_at
FROM campaigns
WHERE id = $1 LIMIT 1;

-- name: UpdateCampaignStatus :exec
UPDATE campaigns
SET is_active = $2, updated_at = NOW()
WHERE id = $1 AND store_id = $3;
