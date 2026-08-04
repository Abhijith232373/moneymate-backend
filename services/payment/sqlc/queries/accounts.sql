-- name: CreateAccount :one
INSERT INTO payment.accounts (user_id, merchant_id, type, currency)
VALUES ($1, $2, $3::payment.account_type, $4)
RETURNING id, user_id, merchant_id, type::text AS type, currency,
          balance, version, created_at, updated_at;

-- name: GetAccountByID :one
SELECT id, user_id, merchant_id, type::text AS type, currency,
       balance, version, created_at, updated_at
FROM payment.accounts
WHERE id = $1;

-- name: GetWalletByUserID :one
SELECT id, user_id, merchant_id, type::text AS type, currency,
       balance, version, created_at, updated_at
FROM payment.accounts
WHERE user_id = $1 AND type = 'wallet';

-- name: ListAccountsByUser :many
SELECT id, user_id, merchant_id, type::text AS type, currency,
       balance, version, created_at, updated_at
FROM payment.accounts
WHERE user_id = $1
ORDER BY created_at;

-- name: GetAccountByIDForUpdate :one
SELECT id, user_id, merchant_id, type::text AS type, currency,
       balance, version, created_at, updated_at
FROM payment.accounts
WHERE id = $1
FOR UPDATE;

-- name: AddBalance :exec
UPDATE payment.accounts
SET balance = balance + $2, version = version + 1, updated_at = NOW()
WHERE id = $1;