-- name: InsertInbox :one
INSERT INTO notification.inbox
    (recipient_type, recipient_id, category, title, body, data, event_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (event_id, recipient_type, recipient_id) DO NOTHING
RETURNING id;

-- name: ListInbox :many
SELECT * FROM notification.inbox
WHERE recipient_type = $1 AND recipient_id = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: GetInbox :one
SELECT * FROM notification.inbox
WHERE id = $1 AND recipient_type = $2 AND recipient_id = $3;

-- name: MarkInboxRead :exec
UPDATE notification.inbox
SET read_at = NOW()
WHERE id = $1 AND recipient_type = $2 AND recipient_id = $3;

-- name: MarkInboxSent :exec
UPDATE notification.inbox
SET status = 'sent', sent_at = NOW()
WHERE id = $1;
