CREATE TABLE auth.user_pins (
    id              UUID        PRIMARY KEY,
    user_id         UUID        NOT NULL UNIQUE REFERENCES auth.users(id) ON DELETE CASCADE,
    pin_hash        TEXT        NOT NULL,              
    failed_attempts INTEGER     NOT NULL DEFAULT 0,   
    locked_until    TIMESTAMPTZ,                      
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_pins_user_id ON auth.user_pins(user_id);