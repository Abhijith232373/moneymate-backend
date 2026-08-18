
DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'auth_user') THEN
    CREATE ROLE auth_user LOGIN PASSWORD 'auth_password';
  END IF;

  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'core_user') THEN
    CREATE ROLE core_user LOGIN PASSWORD 'core_password';
  END IF;

  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'merchant_user') THEN
    CREATE ROLE merchant_user LOGIN PASSWORD 'merchant_password';
  END IF;

  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'rewards_user') THEN
    CREATE ROLE rewards_user LOGIN PASSWORD 'rewards_password';
  END IF;

  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'payment_user') THEN
    CREATE ROLE payment_user LOGIN PASSWORD 'payment_password';
  END IF;

  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'support_user') THEN
    CREATE ROLE support_user LOGIN PASSWORD 'support_password';
  END IF;

  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'notification_user') THEN
    CREATE ROLE notification_user LOGIN PASSWORD 'notification_password';
  END IF;
END $$;