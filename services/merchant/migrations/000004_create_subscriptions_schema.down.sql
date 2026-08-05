-- Drop indexes first
DROP INDEX IF EXISTS idx_merchant_subscriptions_status_end;
DROP INDEX IF EXISTS idx_merchant_subscriptions_store_id;
DROP INDEX IF EXISTS idx_subscription_tiers_code_active;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS merchant_subscriptions;
DROP TABLE IF EXISTS subscription_tiers;


-- Drop types
DROP TYPE IF EXISTS billing_cycle_type;
DROP TYPE IF EXISTS subscription_billing_status;
