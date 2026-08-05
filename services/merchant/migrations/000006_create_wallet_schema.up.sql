-- 1. Create Enums for Wallet Transaction Types
CREATE TYPE wallet_txn_type AS ENUM (
    'qr_scan',
    'redeem',
    'adjustment'
);

-- 2. Wallets Table
-- Tracks aggregate earnings, available balance, and total redeemed per merchant store.
CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL UNIQUE REFERENCES stores(id) ON DELETE CASCADE,
    
    available_balance NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    total_earnings NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    total_redeemed NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 3. Wallet Transactions Table
-- Immutable ledger of all wallet transactions (earnings and redemptions).
CREATE TABLE wallet_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    
    transaction_id VARCHAR(50) NOT NULL UNIQUE, -- e.g., TXN-001
    title VARCHAR(100) NOT NULL,                -- e.g., 'Payout from Admin', 'QR Payment'
    subtitle VARCHAR(100) NOT NULL,             -- e.g., 'Redeem', 'QrScan'
    amount NUMERIC(15, 2) NOT NULL,             -- Positive for earnings, negative for redemptions
    txn_type wallet_txn_type NOT NULL,          -- 'qr_scan', 'redeem', 'adjustment'
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 4. Performance Indices
CREATE INDEX idx_wallets_store_id ON wallets(store_id);
CREATE INDEX idx_wallet_transactions_store_id_created ON wallet_transactions(store_id, created_at DESC);
CREATE INDEX idx_wallet_transactions_txn_type ON wallet_transactions(store_id, txn_type);
