CREATE TYPE payment.account_type AS ENUM (
    'wallet',
    'pod',
    'merchant_settlement',
    'merchant_payout',
    'platform_commission_pool'
);

CREATE TABLE payment.accounts (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID,
    merchant_id UUID,
    type        payment.account_type NOT NULL,
    currency    TEXT NOT NULL DEFAULT 'INR',
    balance     BIGINT NOT NULL DEFAULT 0,
    version     BIGINT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT wallets_balance_check   CHECK (balance >= 0),
    CONSTRAINT chk_account_owner       CHECK (
        user_id IS NOT NULL OR merchant_id IS NOT NULL OR type = 'platform_commission_pool'
    )
);

CREATE UNIQUE INDEX idx_accounts_wallet_per_user ON payment.accounts(user_id) WHERE type = 'wallet';

CREATE INDEX idx_accounts_user_id     ON payment.accounts(user_id);
CREATE INDEX idx_accounts_merchant_id ON payment.accounts(merchant_id);
