CREATE TYPE payment.deposit_status AS ENUM ('created', 'paid', 'failed');

CREATE TABLE payment.deposits (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL,
    account_id          UUID NOT NULL REFERENCES payment.accounts(id),
    razorpay_order_id   TEXT NOT NULL UNIQUE,
    razorpay_payment_id TEXT,
    amount              BIGINT NOT NULL CHECK (amount > 0),
    status              payment.deposit_status NOT NULL DEFAULT 'created',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at        TIMESTAMPTZ
);

CREATE INDEX idx_deposits_user ON payment.deposits(user_id, created_at);