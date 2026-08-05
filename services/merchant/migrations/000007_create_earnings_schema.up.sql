-- 1. Earnings Stats Table
-- Tracks total lifetime scans and total earned from milestones
CREATE TABLE earnings_stats (
    store_id UUID PRIMARY KEY REFERENCES stores(id) ON DELETE CASCADE,
    total_scans BIGINT NOT NULL DEFAULT 0,
    total_earned NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Earnings Payout Requests Table
-- Tracks requested milestone payouts
CREATE TABLE earnings_payout_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    milestone_scans INT NOT NULL,
    reward_amount NUMERIC(15, 2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'requested', -- 'requested', 'paid'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_earnings_payout_requests_store_id ON earnings_payout_requests(store_id);
-- Ensure a store can only request a specific milestone once
CREATE UNIQUE INDEX idx_unique_milestone_payout ON earnings_payout_requests(store_id, milestone_scans);
