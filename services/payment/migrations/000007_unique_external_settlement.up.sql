CREATE UNIQUE INDEX idx_accounts_external_settlement_singleton
ON payment.accounts ((type))
WHERE type = 'external_settlement';