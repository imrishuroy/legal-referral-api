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
SELECT
    m1.*,
    m2.*
FROM
    messages m1
        LEFT JOIN
    messages m2 ON m1.parent_message_id = m2.message_id
WHERE
    m1.room_id = $1
ORDER BY
    m1.sent_at DESC
OFFSET $2
    LIMIT $3;
