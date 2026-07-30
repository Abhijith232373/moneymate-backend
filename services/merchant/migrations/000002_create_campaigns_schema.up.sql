-- 1. Campaigns / Offers Table
-- Separated from core merchant onboarding schema to keep domains modular and independently evolvable.
CREATE TABLE campaigns (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id            UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    
    -- Campaign Setup
    name                VARCHAR(255) NOT NULL,
    offer_type          VARCHAR(100) NOT NULL, 
    reward_value        NUMERIC(10, 2) NOT NULL,
    min_bill_amount     NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    target_audience     VARCHAR(100) NOT NULL,
    banner_url          TEXT,
    
    start_date          DATE NOT NULL,
    end_date            DATE NOT NULL,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Performance Index for active campaigns query optimization
CREATE INDEX idx_campaigns_store_id_active ON campaigns(store_id) WHERE is_active = TRUE;
