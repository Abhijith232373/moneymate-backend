-- name: InsertTransaction :one
INSERT INTO payment.transactions
    (id, from_account_id, to_account_id, amount, status, idempotency_key, description, completed_at)
VALUES
    ($1, $2, $3, $4, $5::payment.tx_status, $6, $7, $8)
RETURNING id, from_account_id, to_account_id, amount, status::text AS status,
          idempotency_key, COALESCE(description, '') AS description, created_at, completed_at;

-- name: GetTransactionByID :one
SELECT id, from_account_id, to_account_id, amount, status::text AS status,
       idempotency_key, COALESCE(description, '') AS description, created_at, completed_at
FROM payment.transactions
WHERE id = $1;

-- name: GetTransactionByIdempotencyKey :one
SELECT id, from_account_id, to_account_id, amount, status::text AS status,
       idempotency_key, COALESCE(description, '') AS description, created_at, completed_at
FROM payment.transactions
WHERE idempotency_key = $1 AND from_account_id = $2;

-- name: UpdateTransactionStatus :exec
UPDATE payment.transactions
SET status = $2::payment.tx_status,
    completed_at = CASE WHEN $2 = 'completed' THEN NOW() ELSE completed_at END
WHERE id = $1;

-- name: ListTransactionsByAccount :many
SELECT t.id, t.from_account_id, t.to_account_id, t.amount, t.status::text AS status,
       t.idempotency_key, COALESCE(t.description, '') AS description, t.created_at, t.completed_at
FROM payment.transactions t
JOIN payment.journal_entries j ON j.transaction_id = t.id
WHERE j.account_id = $1
ORDER BY t.created_at DESC;