-- 1. Campaigns / Offers Table
-- Separated from core merchant onboarding schema to keep domains modular and independently evolvable.
CREATE TABLE campaigns (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id            UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    
    -- Campaign Setup
    name                VARCHAR(255) NOT NULL,
    redeem_code         VARCHAR(50) UNIQUE,
    offer_category      VARCHAR(100),
    offer_type          VARCHAR(100) NOT NULL, 
    reward_value        NUMERIC(10, 2) NOT NULL DEFAULT 0,
    min_bill_amount     NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    redemption_limit    INT NOT NULL DEFAULT 0,
    target_audience     VARCHAR(100) NOT NULL,
    banner_url          TEXT,
    
    start_date          TIMESTAMPTZ NOT NULL,
    end_date            TIMESTAMPTZ NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'active',
    
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Performance Index for active campaigns query optimization
CREATE INDEX idx_campaigns_store_id_active ON campaigns(store_id) WHERE status = 'active';
