-- name: CreatePIN :exec
INSERT INTO auth.user_pins (
    id,
    user_id,
    pin_hash
)
VALUES ($1, $2,$3);

-- name: GetPINByUserID :one
SELECT * FROM auth.user_pins
WHERE user_id = $1;

-- name: UpdatePINHash :exec
UPDATE auth.user_pins
SET
    pin_hash        = $2,
    failed_attempts = 0,
    locked_until    = NULL,
    updated_at      = NOW()
WHERE user_id = $1;

-- name: IncrementPINFailedAttempts :one
UPDATE auth.user_pins
SET
    failed_attempts = failed_attempts + 1,
    updated_at      = NOW()
WHERE user_id = $1
RETURNING failed_attempts;

-- name: LockPIN :exec
UPDATE auth.user_pins
SET
    locked_until = $2,
    updated_at   = NOW()
WHERE user_id = $1;

-- name: ResetPINFailedAttempts :exec
UPDATE auth.user_pins
SET
    failed_attempts = 0,
    locked_until    = NULL,
    updated_at      = NOW()
WHERE user_id = $1;

-- name: PINExists :one
SELECT EXISTS(
    SELECT 1 FROM auth.user_pins WHERE user_id = $1
) AS exists;