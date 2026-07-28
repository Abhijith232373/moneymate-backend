-- Drop indexes first
DROP INDEX IF EXISTS idx_redemption_requests_store_id_status;
DROP INDEX IF EXISTS idx_reward_transactions_display_id;
DROP INDEX IF EXISTS idx_reward_transactions_store_id_created;
DROP INDEX IF EXISTS idx_reward_balances_store_id;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS redemption_requests;
DROP TABLE IF EXISTS reward_transactions;
DROP TABLE IF EXISTS reward_balances;

-- Drop types
DROP TYPE IF EXISTS redemption_status;
DROP TYPE IF EXISTS reward_transaction_type;
