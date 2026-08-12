ALTER TABLE payment.accounts ADD COLUMN handle TEXT;
CREATE UNIQUE INDEX idx_accounts_handle ON payment.accounts(handle) WHERE handle IS NOT NULL;