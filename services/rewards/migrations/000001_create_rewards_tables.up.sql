CREATE TYPE reward_recipient_type AS ENUM ('user', 'merchant');
CREATE TYPE reward_payout_status AS ENUM ('pending', 'completed', 'failed');

CREATE TABLE reward_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    min_percentage_bps INT NOT NULL CHECK (min_percentage_bps >= 0),
    max_percentage_bps INT NOT NULL CHECK (max_percentage_bps >= min_percentage_bps),
    min_transaction_amount_paise BIGINT NOT NULL DEFAULT 0 CHECK (min_transaction_amount_paise >= 0),
    max_payout_amount_paise BIGINT NOT NULL CHECK (max_payout_amount_paise > 0),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_reward_rules_one_active
    ON reward_rules (active)
    WHERE active = TRUE;

CREATE TABLE reward_payouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL,
    recipient_id UUID NOT NULL,
    recipient_account_id UUID NOT NULL,
    recipient_type reward_recipient_type NOT NULL,
    rule_id UUID REFERENCES reward_rules(id),
    original_amount_paise BIGINT NOT NULL CHECK (original_amount_paise > 0),
    reward_percentage_bps INT NOT NULL CHECK (reward_percentage_bps >= 0),
    reward_amount_paise BIGINT NOT NULL CHECK (reward_amount_paise >= 0),
    status reward_payout_status NOT NULL DEFAULT 'pending',
    payment_transaction_id UUID,
    failure_reason TEXT,
    event_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    UNIQUE (transaction_id, recipient_type)
);

CREATE INDEX idx_reward_payouts_recipient_created
    ON reward_payouts (recipient_id, created_at DESC);

CREATE INDEX idx_reward_payouts_transaction_id
    ON reward_payouts (transaction_id);

CREATE INDEX idx_reward_payouts_status
    ON reward_payouts (status);
