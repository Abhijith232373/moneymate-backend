CREATE TABLE notification.preferences (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipient_type TEXT NOT NULL CHECK (recipient_type IN ('user','merchant')),
    recipient_id   UUID NOT NULL,
    category       notification.category NOT NULL,
    enabled        BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT preferences_unique UNIQUE (recipient_type, recipient_id, category)
);
