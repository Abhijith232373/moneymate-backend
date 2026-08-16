-- Auth Service
CREATE ROLE auth_user
LOGIN               
PASSWORD 'auth_password';

-- Core Service
CREATE ROLE core_user
LOGIN
PASSWORD 'core_password';

-- Merchant Service
CREATE ROLE merchant_user
LOGIN
PASSWORD 'merchant_password';

-- Rewards Service
CREATE ROLE rewards_user
LOGIN
PASSWORD 'rewards_password';

-- Automation Service
CREATE ROLE automation_user
LOGIN
PASSWORD 'automation_password';

-- Payment Service
CREATE ROLE payment_user
LOGIN
PASSWORD 'payment_password';

-- Support Service
CREATE ROLE support_user
LOGIN
PASSWORD 'support_password';

-- notification Service
CREATE ROLE notification_user
LOGIN
PASSWORD 'notification_password';