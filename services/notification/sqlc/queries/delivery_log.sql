-- name: InsertDeliveryLog :exec
INSERT INTO notification.delivery_log
    (inbox_id, device_token_id, provider, provider_message_id, status, error_code)
VALUES ($1, $2, 'fcm', $3, $4, $5);
