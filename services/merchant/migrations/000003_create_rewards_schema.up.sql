-- 1. Create Enums for Reward Transaction Types and Redemption Statuses
CREATE TYPE reward_transaction_type AS ENUM (
    'earning',
    'redemption',
    'adjustment',
    'bonus'
);

CREATE TYPE redemption_status AS ENUM (
    'processing',
    'completed',
    'failed',
    'rejected'
);

-- 2. Reward Balances Table (Tracks aggregate earnings, balance, and quick stats per merchant store)
-- Optimized for high-concurrency read/write operations for millions of active merchants.
CREATE TABLE reward_balances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL UNIQUE REFERENCES stores(id) ON DELETE CASCADE,
    
    -- Financial and statistical metrics displayed on the Rewards Center UI
    available_balance NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    total_scans BIGINT NOT NULL DEFAULT 0,
    premium_points BIGINT NOT NULL DEFAULT 0,
    weekly_growth_percentage NUMERIC(5, 2) NOT NULL DEFAULT 0.00,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 3. Reward Transactions Table (Immutable ledger of all QR scan earnings and redemptions)
CREATE TABLE reward_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    
    -- Transaction presentation details matching UI cards
    campaign_name VARCHAR(255) NOT NULL,
    display_id VARCHAR(50) NOT NULL, -- e.g., #QR-8820
    status VARCHAR(50) NOT NULL DEFAULT 'Settled', -- Settled, Pending, Redeemed
    amount NUMERIC(15, 2) NOT NULL,
    transaction_type reward_transaction_type NOT NULL DEFAULT 'earning',
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 4. Redemption Requests Table (Audit trail for bank transfer withdrawals)
CREATE TABLE redemption_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    
    amount NUMERIC(15, 2) NOT NULL,
    bank_transfer_authorized BOOLEAN NOT NULL DEFAULT TRUE,
    status redemption_status NOT NULL DEFAULT 'processing',
    reference_id VARCHAR(100) UNIQUE, -- Bank payout reference number
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 5. Performance Indices for high-volume querying and filtering
CREATE INDEX idx_reward_balances_store_id ON reward_balances(store_id);
CREATE INDEX idx_reward_transactions_store_id_created ON reward_transactions(store_id, created_at DESC);
CREATE INDEX idx_reward_transactions_display_id ON reward_transactions(display_id);
CREATE INDEX idx_redemption_requests_store_id_status ON redemption_requests(store_id, status);
