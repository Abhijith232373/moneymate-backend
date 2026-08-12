-- 000XXX_create_outbox.up.sql
CREATE TABLE auth.outbox_events (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    topic        TEXT NOT NULL,
    payload      JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ
);

CREATE INDEX idx_outbox_unpublished ON auth.outbox_events(created_at) WHERE published_at IS NULL;