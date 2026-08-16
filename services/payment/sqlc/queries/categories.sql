-- name: CreateCategory :one
INSERT INTO payment.categories (user_id, name)
VALUES ($1, $2)
RETURNING id, user_id, name, created_at, updated_at;

-- name: ListCategoriesByUser :many
SELECT id, user_id, name, created_at, updated_at
FROM payment.categories
WHERE user_id = $1
ORDER BY name;

-- name: GetCategoryByID :one
SELECT id, user_id, name, created_at, updated_at
FROM payment.categories
WHERE id = $1;

-- name: UpdateCategory :one
UPDATE payment.categories
SET name = $2, updated_at = NOW()
WHERE id = $1 AND user_id = $3
RETURNING id, user_id, name, created_at, updated_at;

-- name: DeleteCategory :exec
DELETE FROM payment.categories
WHERE id = $1 AND user_id = $2;


