CREATE TABLE auth.staff (
    id                UUID        PRIMARY KEY,
    email             VARCHAR(255) NOT NULL UNIQUE,
    full_name         TEXT         NOT NULL,
    password_hash     TEXT         NOT NULL,
    status            auth.user_status NOT NULL DEFAULT 'active',
    token_version     BIGINT       NOT NULL DEFAULT 0,
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_staff_email ON auth.staff(email);

CREATE TABLE auth.staff_roles (
    staff_id UUID NOT NULL REFERENCES auth.staff(id) ON DELETE CASCADE,
    role_id  UUID NOT NULL REFERENCES auth.roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (staff_id, role_id)
);
