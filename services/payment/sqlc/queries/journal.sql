-- name: InsertJournalEntry :exec
INSERT INTO payment.journal_entries (id, transaction_id, account_id, amount, direction)
VALUES ($1, $2, $3, $4, $5::payment.tx_direction);

-- name: GetEntriesByTransactionID :many
SELECT id, transaction_id, account_id, amount, direction::text AS direction, created_at
FROM payment.journal_entries
WHERE transaction_id = $1
ORDER BY created_at;

-- name: ListEntriesByAccount :many
SELECT id, transaction_id, account_id, amount, direction::text AS direction, created_at
FROM payment.journal_entries
WHERE account_id = $1
ORDER BY created_at DESC;