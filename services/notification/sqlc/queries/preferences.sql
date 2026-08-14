-- name: GetPreference :one
SELECT enabled FROM notification.preferences
WHERE recipient_type = $1 AND recipient_id = $2 AND category = $3;

-- name: UpsertPreference :one
INSERT INTO notification.preferences (recipient_type, recipient_id, category, enabled)
VALUES ($1, $2, $3, $4)
ON CONFLICT (recipient_type, recipient_id, category)
DO UPDATE SET enabled = EXCLUDED.enabled, updated_at = NOW()
RETURNING *;

-- name: ListPreferences :many
SELECT * FROM notification.preferences
WHERE recipient_type = $1 AND recipient_id = $2;
