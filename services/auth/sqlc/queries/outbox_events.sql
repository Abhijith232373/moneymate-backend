-- name: InsertOutboxEvent :exec
INSERT INTO auth.outbox_events (id, topic, payload)
VALUES ($1, $2, $3);

-- name: FetchUnpublishedOutboxEvents :many
SELECT * FROM auth.outbox_events
WHERE published_at IS NULL
ORDER BY created_at ASC
LIMIT $1;

-- name: MarkOutboxEventPublished :exec
UPDATE auth.outbox_events
SET published_at = NOW()
WHERE id = $1;