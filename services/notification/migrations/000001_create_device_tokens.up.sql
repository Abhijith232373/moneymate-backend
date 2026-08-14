CREATE TABLE notification.device_tokens (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipient_type TEXT NOT NULL CHECK (recipient_type IN ('user','merchant')),
    recipient_id   UUID NOT NULL,
    device_id      TEXT NOT NULL,
    token          TEXT NOT NULL UNIQUE,
    platform       TEXT NOT NULL CHECK (platform IN ('ios','android','web')),
    app_version    TEXT,
    is_active      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT device_tokens_one_per_device UNIQUE (recipient_type, recipient_id, device_id)
);
CREATE INDEX idx_device_tokens_recipient ON notification.device_tokens(recipient_type, recipient_id) WHERE is_active;