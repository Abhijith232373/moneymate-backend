DROP INDEX IF EXISTS idx_reward_payouts_status;
DROP INDEX IF EXISTS idx_reward_payouts_transaction_id;
DROP INDEX IF EXISTS idx_reward_payouts_recipient_created;
DROP TABLE IF EXISTS reward_payouts;

DROP INDEX IF EXISTS idx_reward_rules_one_active;
DROP TABLE IF EXISTS reward_rules;

DROP TYPE IF EXISTS reward_payout_status;
DROP TYPE IF EXISTS reward_recipient_type;
