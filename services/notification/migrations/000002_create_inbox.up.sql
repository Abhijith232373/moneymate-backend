CREATE TYPE notification.category AS ENUM
    ('bill_due','debt','transfer','merchant','campaign','offer','promo','system');

CREATE TABLE notification.inbox (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipient_type TEXT NOT NULL CHECK (recipient_type IN ('user','merchant')),
    recipient_id   UUID NOT NULL,
    category       notification.category NOT NULL,
    title          TEXT NOT NULL,
    body           TEXT NOT NULL,
    data           JSONB NOT NULL DEFAULT '{}',
    event_id       TEXT NOT NULL,
    status         TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','sent','failed')),
    read_at        TIMESTAMPTZ,
    sent_at        TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT inbox_dedup UNIQUE (event_id, recipient_type, recipient_id)
);
CREATE INDEX idx_inbox_recipient ON notification.inbox(recipient_type, recipient_id, created_at DESC);
