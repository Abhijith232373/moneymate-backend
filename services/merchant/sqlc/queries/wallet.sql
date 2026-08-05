-- name: GetWalletByStoreID :one
SELECT id, store_id, available_balance, total_earnings, total_redeemed, created_at, updated_at
FROM wallets
WHERE store_id = $1 LIMIT 1;

-- name: UpsertWallet :one
INSERT INTO wallets (
    store_id, available_balance, total_earnings, total_redeemed
) VALUES (
    $1, $2, $3, $4
) ON CONFLICT (store_id) DO UPDATE SET
    available_balance = EXCLUDED.available_balance,
    total_earnings = EXCLUDED.total_earnings,
    total_redeemed = EXCLUDED.total_redeemed,
    updated_at = NOW()
RETURNING id, store_id, available_balance, total_earnings, total_redeemed, created_at, updated_at;

-- name: GetWalletTransactions :many
SELECT id, store_id, transaction_id, title, subtitle, amount, txn_type, created_at
FROM wallet_transactions
WHERE store_id = $1
ORDER BY created_at DESC;

-- name: GetWalletTransactionsByType :many
SELECT id, store_id, transaction_id, title, subtitle, amount, txn_type, created_at
FROM wallet_transactions
WHERE store_id = $1 AND txn_type = $2
ORDER BY created_at DESC;

-- name: CreateWalletTransaction :one
INSERT INTO wallet_transactions (
    store_id, transaction_id, title, subtitle, amount, txn_type
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING id, store_id, transaction_id, title, subtitle, amount, txn_type, created_at;
