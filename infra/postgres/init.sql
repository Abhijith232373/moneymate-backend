

SELECT 'CREATE DATABASE moneymate'
WHERE NOT EXISTS (
    SELECT FROM pg_database WHERE datname = 'moneymate'
)\gexec

\c moneymate;

-- ── Schemas ─────────────────────────────────────────────────────
CREATE SCHEMA IF NOT EXISTS auth;
CREATE SCHEMA IF NOT EXISTS core;
CREATE SCHEMA IF NOT EXISTS merchant;
CREATE SCHEMA IF NOT EXISTS rewards;
CREATE SCHEMA IF NOT EXISTS automation;

-- ── Auth Role ───────────────────────────────────────────────────
DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'auth_user') THEN
    CREATE ROLE auth_user LOGIN PASSWORD 'auth_password_change_in_prod';
  END IF;
END
$$;

GRANT CONNECT ON DATABASE moneymate TO auth_user;
GRANT USAGE ON SCHEMA auth TO auth_user;
GRANT ALL ON ALL TABLES IN SCHEMA auth TO auth_user;
GRANT ALL ON ALL SEQUENCES IN SCHEMA auth TO auth_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA auth
    GRANT ALL ON TABLES TO auth_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA auth
    GRANT ALL ON SEQUENCES TO auth_user;

-- ── Core Role ───────────────────────────────────────────────────
DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'core_user') THEN
    CREATE ROLE core_user LOGIN PASSWORD 'core_password_change_in_prod';
  END IF;
END
$$;

GRANT CONNECT ON DATABASE moneymate TO core_user;
GRANT USAGE ON SCHEMA core TO core_user;
GRANT ALL ON ALL TABLES IN SCHEMA core TO core_user;
GRANT ALL ON ALL SEQUENCES IN SCHEMA core TO core_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA core
    GRANT ALL ON TABLES TO core_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA core
    GRANT ALL ON SEQUENCES TO core_user;

-- ── Merchant Role ───────────────────────────────────────────────
DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'merchant_user') THEN
    CREATE ROLE merchant_user LOGIN PASSWORD 'merchant_password_change_in_prod';
  END IF;
END
$$;

GRANT CONNECT ON DATABASE moneymate TO merchant_user;
GRANT USAGE ON SCHEMA merchant TO merchant_user;
GRANT ALL ON ALL TABLES IN SCHEMA merchant TO merchant_user;
GRANT ALL ON ALL SEQUENCES IN SCHEMA merchant TO merchant_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA merchant
    GRANT ALL ON TABLES TO merchant_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA merchant
    GRANT ALL ON SEQUENCES TO merchant_user;

