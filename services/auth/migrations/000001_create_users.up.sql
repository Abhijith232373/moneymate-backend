CREATE TYPE auth.user_status AS ENUM (
    'pending',
    'active',
    'suspended',
    'deleted'
);

CREATE TABLE auth.users (
    id                UUID        PRIMARY KEY ,
    email             VARCHAR(255) NOT NULL UNIQUE,
    phone             VARCHAR(20)  UNIQUE,      
    full_name         TEXT         NOT NULL,
    handle            TEXT         NOT NULL UNIQUE,      -- @username for P2P
    password_hash     TEXT,      
    status            auth.user_status NOT NULL DEFAULT 'pending',
    token_version     BIGINT       NOT NULL DEFAULT 0,

    is_email_verified BOOLEAN      NOT NULL DEFAULT FALSE,
    is_phone_verified BOOLEAN      NOT NULL DEFAULT FALSE,

    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email  ON auth.users(email);
CREATE INDEX idx_users_handle ON auth.users(handle);
CREATE INDEX idx_users_phone  ON auth.users(phone) WHERE phone IS NOT NULL;