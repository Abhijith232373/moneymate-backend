-- 1. Create Enums for Subscription Billing Statuses and Billing Cycles
CREATE TYPE subscription_billing_status AS ENUM (
    'active',
    'past_due',
    'canceled',
    'trialing'
);

CREATE TYPE billing_cycle_type AS ENUM (
    'monthly',
    'annual'
);

-- 2. Subscription Tiers Catalog Table
-- Stores dynamic pricing, features, and limits for each tier (Essential, Growth, Enterprise).
-- Built for millions of merchants so plan configurations can evolve without code deployments.
CREATE TABLE subscription_tiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_code VARCHAR(50) NOT NULL UNIQUE, -- 'essential', 'growth', 'enterprise'
    name VARCHAR(100) NOT NULL,            -- 'Essential', 'Growth', 'Enterprise'
    price NUMERIC(10, 2) NOT NULL,         -- 0.00, 29.00, 99.00
    billing_cycle billing_cycle_type NOT NULL DEFAULT 'monthly',
    description TEXT NOT NULL,
    max_active_campaigns INT NOT NULL,     -- 1 for Essential, 5 for Growth, -1 (unlimited) for Enterprise
    is_most_popular BOOLEAN NOT NULL DEFAULT FALSE,
    features JSONB NOT NULL,               -- Array of feature strings or UI check/cross items
    
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 3. Merchant Subscriptions Table
-- Tracks the active subscription lifecycle, billing dates, and tier for every store.
CREATE TABLE merchant_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    store_id UUID NOT NULL UNIQUE REFERENCES stores(id) ON DELETE CASCADE,
    plan_code VARCHAR(50) NOT NULL REFERENCES subscription_tiers(plan_code) ON UPDATE CASCADE,
    
    status subscription_billing_status NOT NULL DEFAULT 'active',
    billing_cycle billing_cycle_type NOT NULL DEFAULT 'monthly',
    current_period_start TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    current_period_end TIMESTAMPTZ NOT NULL DEFAULT (NOW() + INTERVAL '30 days'),
    auto_renew BOOLEAN NOT NULL DEFAULT TRUE,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 5. Performance Indices for high-volume billing queries and plan lookups
CREATE INDEX idx_subscription_tiers_code_active ON subscription_tiers(plan_code) WHERE is_active = TRUE;

CREATE INDEX idx_merchant_subscriptions_store_id ON merchant_subscriptions(store_id);
CREATE INDEX idx_merchant_subscriptions_status_end ON merchant_subscriptions(status, current_period_end);
