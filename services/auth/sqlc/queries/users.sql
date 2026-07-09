-- name: CreateUser :one
INSERT INTO auth.users (
    id,
    email,
    phone,
    full_name,
    handle,
    password_hash
)
VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;



-- name: GetUserByID :one
SELECT * FROM auth.users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM auth.users
WHERE email = $1;

-- name: GetUserByPhone :one
SELECT * FROM auth.users
WHERE phone = $1;

-- name: GetUserByHandle :one
SELECT * FROM auth.users
WHERE handle = $1;

-- name: VerifyEmail :exec
UPDATE auth.users
SET
    is_email_verified  = TRUE,
    status = 'active',
    updated_at = NOW()
WHERE id = $1;

-- name: VerifyPhone :exec
UPDATE auth.users
SET
    is_phone_verified = TRUE,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdatePassword :exec
UPDATE auth.users
SET
    password_hash = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserStatus :exec 
UPDATE auth.users
SET
    status = $2,
    updated_at = NOW()
WHERE id = $1;


-- name: IncrementTokenVersion :one
UPDATE auth.users
SET
    token_version = token_version + 1,
    updated_at = NOW()
WHERE id = $1
RETURNING token_version;

-- name: GetTokenVersion :one
SELECT token_version FROM auth.users
WHERE id = $1;

-- name: SoftDeleteUser :exec
UPDATE auth.users
SET
    status = 'deleted',
    updated_at = NOW()
WHERE id = $1;

-- name: HandleExists :one
SELECT EXISTS(
    SELECT 1 FROM auth.users WHERE handle = $1
) AS exists;

-- name: EmailExists :one
SELECT EXISTS(
    SELECT 1 FROM auth.users WHERE email = $1
) AS exists;

-- name: PhoneExists :one
SELECT EXISTS(
    SELECT 1 FROM auth.users WHERE phone = $1
) AS exists;

-- name: ListUsers :many
SELECT * FROM auth.users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT count(*) FROM auth.users;