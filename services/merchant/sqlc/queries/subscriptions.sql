-- name: GetSubscriptionPlans :many
-- GetSubscriptionPlans retrieves the full catalog of active pricing tiers ordered by cost.
SELECT 
    id, plan_code, name, price, billing_cycle, description, max_active_campaigns, is_most_popular, features, is_active, created_at, updated_at
FROM subscription_tiers
WHERE is_active = TRUE
ORDER BY price ASC;

-- name: GetSubscriptionByStoreID :one
-- GetSubscriptionByStoreID fetches the active billing record and renewal lifecycle for a specific merchant store.
SELECT 
    id, store_id, plan_code, status, billing_cycle, current_period_start, current_period_end, auto_renew, created_at, updated_at
FROM merchant_subscriptions
WHERE store_id = $1 LIMIT 1;

-- name: CreateSubscriptionPlan :one
-- CreateSubscriptionPlan inserts or updates a pricing tier catalog entry with JSONB features and promotional limits.
INSERT INTO subscription_tiers (
    plan_code, name, price, billing_cycle, description, max_active_campaigns, is_most_popular, features, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) ON CONFLICT (plan_code) DO UPDATE SET

    name = EXCLUDED.name,
    price = EXCLUDED.price,
    billing_cycle = EXCLUDED.billing_cycle,
    description = EXCLUDED.description,
    max_active_campaigns = EXCLUDED.max_active_campaigns,
    is_most_popular = EXCLUDED.is_most_popular,
    features = EXCLUDED.features,
    is_active = EXCLUDED.is_active,
    updated_at = NOW()
RETURNING id, plan_code, name, price, billing_cycle, description, max_active_campaigns, is_most_popular, features, is_active, created_at, updated_at;

-- name: UpsertMerchantSubscription :one
-- UpsertMerchantSubscription atomically creates or updates a store's billing subscription record.
INSERT INTO merchant_subscriptions (
    id, store_id, plan_code, status, billing_cycle, current_period_start, current_period_end, auto_renew
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) ON CONFLICT (store_id) DO UPDATE SET
    plan_code = EXCLUDED.plan_code,
    status = EXCLUDED.status,
    billing_cycle = EXCLUDED.billing_cycle,
    current_period_start = EXCLUDED.current_period_start,
    current_period_end = EXCLUDED.current_period_end,
    auto_renew = EXCLUDED.auto_renew,
    updated_at = NOW()
RETURNING id, store_id, plan_code, status, billing_cycle, current_period_start, current_period_end, auto_renew, created_at, updated_at;


-- name: UpdateStorePlanEnum :exec
-- UpdateStorePlanEnum synchronizes the core store record's plan enum column with the active subscription tier.
UPDATE stores
SET plan = sqlc.arg('plan')::subscription_plan, updated_at = NOW()
WHERE id = sqlc.arg('id');

