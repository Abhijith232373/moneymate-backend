-- name: CreateDeposit :one
INSERT INTO payment.deposits (
    id,
    user_id,
    account_id,
    razorpay_order_id,
    amount
)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetDepositByOrderID :one
SELECT * FROM payment.deposits
WHERE razorpay_order_id = $1;

-- name: MarkDepositPaid :exec
UPDATE payment.deposits
SET status = 'paid', razorpay_payment_id = $2, completed_at = $3
WHERE id = $1;

-- name: MarkDepositFailed :exec
UPDATE payment.deposits
SET status = 'failed', razorpay_payment_id = $2
WHERE id = $1;

-- name: ListDeposits :many
SELECT * FROM payment.deposits
WHERE
    (sqlc.narg('status')::payment.deposit_status IS NULL OR status = sqlc.narg('status'))
    AND (sqlc.narg('user_id')::uuid IS NULL OR user_id = sqlc.narg('user_id'))
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountDeposits :one
SELECT COUNT(*) FROM payment.deposits
WHERE
    (sqlc.narg('status')::payment.deposit_status IS NULL OR status = sqlc.narg('status'))
    AND (sqlc.narg('user_id')::uuid IS NULL OR user_id = sqlc.narg('user_id'));

-- name: ListWithdrawals :many
SELECT * FROM payment.transactions
WHERE to_account_id = @settlement_account_id::uuid
    AND (sqlc.narg('from_account_id')::uuid IS NULL OR from_account_id = sqlc.narg('from_account_id'))
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountWithdrawals :one
SELECT COUNT(*) FROM payment.transactions
WHERE to_account_id = @settlement_account_id::uuid
    AND (sqlc.narg('from_account_id')::uuid IS NULL OR from_account_id = sqlc.narg('from_account_id'));

-- name: MarkDepositPaidIfCreated :execrows
UPDATE payment.deposits
SET status = 'paid', razorpay_payment_id = $2, completed_at = $3
WHERE id = $1 AND status = 'created';