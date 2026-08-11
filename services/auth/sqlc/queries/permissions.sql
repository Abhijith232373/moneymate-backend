-- name: CreatePermission :one
INSERT INTO auth.permissions (
    id,
    name,
    description
)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetPermissionByID :one
SELECT *
FROM auth.permissions
WHERE id = $1;

-- name: UpdatePermission :one
UPDATE auth.permissions
SET name = $2, description = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: GetPermissionByName :one
SELECT *
FROM auth.permissions
WHERE name = $1;

-- name: ListPermissions :many
SELECT *
FROM auth.permissions
ORDER BY name;

-- name: DeletePermission :exec
DELETE FROM auth.permissions
WHERE id = $1;