-- name: CreateMessage :one
INSERT INTO messages (
    parent_message_id,
    sender_id,
    recipient_id,
    message,
    has_attachment,
    attachment_id,
    is_read,
    room_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: ListMessages :many
SELECT * FROM messages
WHERE room_id = $1;
