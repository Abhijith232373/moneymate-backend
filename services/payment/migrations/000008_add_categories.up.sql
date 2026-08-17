-- 0000XX_add_categories.up.sql
CREATE TABLE payment.categories (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL,
    name       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT categories_user_name_unique UNIQUE (user_id, name)
);
CREATE INDEX idx_categories_user ON payment.categories(user_id);

ALTER TABLE payment.transactions
    ADD COLUMN category_id UUID REFERENCES payment.categories(id) ON DELETE SET NULL;
CREATE INDEX idx_transactions_category ON payment.transactions(category_id);
