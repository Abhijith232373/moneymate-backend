CREATE TYPE payment.tx_status AS ENUM ('pending', 'completed', 'failed', 'reversed');
CREATE TYPE payment.tx_direction AS ENUM ('debit', 'credit');

CREATE TABLE payment.transactions (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_account_id  UUID NOT NULL REFERENCES payment.accounts(id),
    to_account_id    UUID NOT NULL REFERENCES payment.accounts(id),
    amount           BIGINT NOT NULL CHECK (amount > 0),  -- paise
    status           payment.tx_status NOT NULL DEFAULT 'pending',
    idempotency_key  TEXT NOT NULL,
    description      TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at     TIMESTAMPTZ,
    CONSTRAINT transactions_idempotency_key UNIQUE (idempotency_key, from_account_id),
    CONSTRAINT chk_distinct_accounts CHECK (from_account_id <> to_account_id)
);

CREATE TABLE payment.journal_entries (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES payment.transactions(id),
    account_id     UUID NOT NULL REFERENCES payment.accounts(id),
    amount         BIGINT NOT NULL CHECK (amount > 0),  -- paise
    direction      payment.tx_direction NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_from      ON payment.transactions(from_account_id, created_at);
CREATE INDEX idx_transactions_to        ON payment.transactions(to_account_id, created_at);
CREATE INDEX idx_transactions_idemkey   ON payment.transactions(idempotency_key, from_account_id);
CREATE INDEX idx_journal_tx             ON payment.journal_entries(transaction_id);
CREATE INDEX idx_journal_account        ON payment.journal_entries(account_id, created_at);