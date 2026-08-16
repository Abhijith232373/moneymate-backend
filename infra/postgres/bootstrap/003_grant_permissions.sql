-- 1. Revoke default public access from schemas (Security Best Practice)
REVOKE ALL ON SCHEMA auth FROM PUBLIC;
REVOKE ALL ON SCHEMA core FROM PUBLIC;
REVOKE ALL ON SCHEMA merchant FROM PUBLIC;
REVOKE ALL ON SCHEMA rewards FROM PUBLIC;
REVOKE ALL ON SCHEMA automation FROM PUBLIC;
REVOKE ALL ON SCHEMA payment FROM PUBLIC;
REVOKE ALL ON SCHEMA support FROM PUBLIC;

-- 2. Grant USAGE (the ability to see and access the schema) to the correct user
GRANT USAGE ON SCHEMA auth TO auth_user;
GRANT USAGE ON SCHEMA core TO core_user;
GRANT USAGE ON SCHEMA merchant TO merchant_user;
GRANT USAGE ON SCHEMA rewards TO rewards_user;
GRANT USAGE ON SCHEMA automation TO automation_user;
GRANT USAGE ON SCHEMA payment TO payment_user;
GRANT USAGE ON SCHEMA support TO support_user;

-- 3. Grant full CRUD rights on any EXISTING tables/sequences
-- (Even though they are empty now, this is good for idempotency)
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA auth TO auth_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA auth TO auth_user;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA core TO core_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA core TO core_user;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA merchant TO merchant_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA merchant TO merchant_user;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA rewards TO rewards_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA rewards TO rewards_user;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA automation TO automation_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA automation TO automation_user;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA payment TO payment_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA payment TO payment_user;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA support TO support_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA support TO support_user;





GRANT CREATE ON SCHEMA auth TO auth_user;
GRANT CREATE ON SCHEMA core TO core_user;
GRANT CREATE ON SCHEMA merchant TO merchant_user;
GRANT CREATE ON SCHEMA rewards TO rewards_user;
GRANT CREATE ON SCHEMA automation TO automation_user;
GRANT CREATE ON SCHEMA payment TO payment_user;
GRANT CREATE ON SCHEMA support TO support_user;


-- 4. Grant CREATE on the database itself
-- (needed so each service's migration tool can run "CREATE SCHEMA IF NOT EXISTS ..."
--  even when the schema already exists — Postgres checks this permission regardless)
GRANT CREATE ON DATABASE moneymate TO auth_user;
GRANT CREATE ON DATABASE moneymate TO core_user;
GRANT CREATE ON DATABASE moneymate TO merchant_user;
GRANT CREATE ON DATABASE moneymate TO rewards_user;
GRANT CREATE ON DATABASE moneymate TO automation_user;
GRANT CREATE ON DATABASE moneymate TO payment_user;
GRANT CREATE ON DATABASE moneymate TO support_user;



REVOKE ALL ON SCHEMA notification FROM PUBLIC;
GRANT USAGE ON SCHEMA notification TO notification_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA notification TO notification_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA notification TO notification_user;
GRANT CREATE ON SCHEMA notification TO notification_user;
GRANT CREATE ON DATABASE moneymate TO notification_user;    