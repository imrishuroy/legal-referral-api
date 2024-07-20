-- name: SendMessageToDiscussion :one
INSERT INTO discussion_messages (
    parent_message_id,
    sender_id,
    message,
    discussion_id
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: ListDiscussionMessages :many
SELECT
    m1.*,
    u1.user_id AS sender_id,
    u1.first_name AS sender_first_name,
    u1.last_name AS sender_last_name,
    u1.avatar_url AS sender_avatar_image
FROM
    discussion_messages m1
    LEFT JOIN
    users u1 ON m1.sender_id = u1.user_id
WHERE
    m1.discussion_id = $1
ORDER BY
    m1.sent_at ASC
OFFSET $2
LIMIT $3;