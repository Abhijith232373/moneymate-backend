-- Revoke default public access
REVOKE ALL ON SCHEMA auth FROM PUBLIC;
REVOKE ALL ON SCHEMA core FROM PUBLIC;
REVOKE ALL ON SCHEMA merchant FROM PUBLIC;
REVOKE ALL ON SCHEMA rewards FROM PUBLIC;
REVOKE ALL ON SCHEMA payment FROM PUBLIC;
REVOKE ALL ON SCHEMA support FROM PUBLIC;
REVOKE ALL ON SCHEMA notification FROM PUBLIC;

-- Grant USAGE
GRANT USAGE ON SCHEMA auth TO auth_user;
GRANT USAGE ON SCHEMA core TO core_user;
GRANT USAGE ON SCHEMA merchant TO merchant_user;
GRANT USAGE ON SCHEMA rewards TO rewards_user;
GRANT USAGE ON SCHEMA payment TO payment_user;
GRANT USAGE ON SCHEMA support TO support_user;
GRANT USAGE ON SCHEMA notification TO notification_user;

-- Grant CRUD on existing tables/sequences
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA auth TO auth_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA auth TO auth_user;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA core TO core_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA core TO core_user;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA merchant TO merchant_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA merchant TO merchant_user;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA rewards TO rewards_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA rewards TO rewards_user;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA payment TO payment_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA payment TO payment_user;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA support TO support_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA support TO support_user;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA notification TO notification_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA notification TO notification_user;

-- Grant CREATE on schemas (for migration tools)
GRANT CREATE ON SCHEMA auth TO auth_user;
GRANT CREATE ON SCHEMA core TO core_user;
GRANT CREATE ON SCHEMA merchant TO merchant_user;
GRANT CREATE ON SCHEMA rewards TO rewards_user;
GRANT CREATE ON SCHEMA payment TO payment_user;
GRANT CREATE ON SCHEMA support TO support_user;
GRANT CREATE ON SCHEMA notification TO notification_user;

-- Grant CREATE on database
GRANT CREATE ON DATABASE moneymate TO auth_user;
GRANT CREATE ON DATABASE moneymate TO core_user;
GRANT CREATE ON DATABASE moneymate TO merchant_user;
GRANT CREATE ON DATABASE moneymate TO rewards_user;
GRANT CREATE ON DATABASE moneymate TO payment_user;
GRANT CREATE ON DATABASE moneymate TO support_user;
GRANT CREATE ON DATABASE moneymate TO notification_user;