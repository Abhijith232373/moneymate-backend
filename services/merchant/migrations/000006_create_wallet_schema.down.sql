-- Drop indexes
DROP INDEX IF EXISTS idx_wallet_transactions_txn_type;
DROP INDEX IF EXISTS idx_wallet_transactions_store_id_created;
DROP INDEX IF EXISTS idx_wallets_store_id;

-- Drop tables
DROP TABLE IF EXISTS wallet_transactions;
DROP TABLE IF EXISTS wallets;

-- Drop types
DROP TYPE IF EXISTS wallet_txn_type;
