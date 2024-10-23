-- name: CreateNotification :one
INSERT INTO notifications (
    user_id,
    sender_id,
    target_id,
    target_type,
    notification_type,
    message
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: ListNotifications :many
SELECT n.*, u.first_name AS sender_first_name, u.last_name AS sender_last_name, u.avatar_url AS sender_avatar_url
FROM notifications n
JOIN users u ON n.sender_id = u.user_id
WHERE n.user_id = $1
ORDER BY n.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetNotificationById :one
SELECT * FROM notifications WHERE notification_id = $1;

-- name: MarkNotificationAsRead :one
UPDATE notifications SET is_read = true WHERE notification_id = $1 RETURNING *;

-- name: DeleteNotificationById :one
DELETE FROM notifications WHERE notification_id = $1 RETURNING *;
