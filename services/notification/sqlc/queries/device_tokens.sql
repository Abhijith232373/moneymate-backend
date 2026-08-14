-- name: UpsertDeviceToken :one
INSERT INTO notification.device_tokens
    (recipient_type, recipient_id, device_id, token, platform, app_version)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (recipient_type, recipient_id, device_id)
DO UPDATE SET token = EXCLUDED.token,
              platform = EXCLUDED.platform,
              app_version = EXCLUDED.app_version,
              is_active = TRUE,
              updated_at = NOW()
RETURNING id, recipient_type, recipient_id, device_id, token, platform, app_version, is_active, created_at, updated_at;

-- name: ListActiveTokensByRecipient :many
SELECT id, recipient_type, recipient_id, device_id, token, platform, app_version, is_active, created_at, updated_at
FROM notification.device_tokens
WHERE recipient_type = $1 AND recipient_id = $2 AND is_active = TRUE;

-- name: DeactivateDeviceToken :exec
UPDATE notification.device_tokens
SET is_active = FALSE, updated_at = NOW()
WHERE id = $1;

-- name: DeactivateTokensByDevice :exec
UPDATE notification.device_tokens
SET is_active = FALSE, updated_at = NOW()
WHERE recipient_type = $1 AND recipient_id = $2 AND device_id = $3;
