DROP INDEX IF EXISTS payment.idx_accounts_handle;
ALTER TABLE payment.accounts DROP COLUMN IF EXISTS handle;