
ALTER TABLE payment.accounts DROP CONSTRAINT chk_account_owner;
ALTER TABLE payment.accounts ADD CONSTRAINT chk_account_owner CHECK (
    user_id IS NOT NULL OR merchant_id IS NOT NULL OR type = 'platform_commission_pool'
);