CREATE TABLE notification.delivery_log (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inbox_id            UUID NOT NULL REFERENCES notification.inbox(id) ON DELETE CASCADE,
    device_token_id     UUID REFERENCES notification.device_tokens(id),
    provider            TEXT NOT NULL DEFAULT 'fcm',
    provider_message_id TEXT,
    status              TEXT NOT NULL CHECK (status IN ('pending','sent','failed','dropped')),
    error_code          TEXT,
    attempt_count       INT NOT NULL DEFAULT 0,
    last_attempt_at     TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_delivery_log_inbox ON notification.delivery_log(inbox_id);
