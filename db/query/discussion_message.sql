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
    m2.*
FROM
    discussion_messages m1
        LEFT JOIN
    discussion_messages m2 ON m1.parent_message_id = m2.message_id
WHERE
    m1.discussion_id = $1
ORDER BY
    m1.sent_at DESC
OFFSET $2
LIMIT $3;

-- name: ListDiscussionMessages2 :many
SELECT
    m1.message_id,
    m1.parent_message_id,
    m1.discussion_id,
    m1.sender_id,
    m1.message,
    m1.sent_at,
    u1.avatar_url AS sender_avatar_image,
    u1.first_name AS sender_first_name,
    u1.last_name AS sender_last_name,
    m2.message_id AS reply_message_id,
    m2.parent_message_id AS reply_parent_message_id,
    m2.discussion_id AS reply_discussion_id,
    m2.sender_id AS reply_sender_id,
    m2.message AS reply_message,
    m2.sent_at AS reply_sent_at,
    u2.avatar_url AS reply_sender_avatar_image,
    u2.first_name AS reply_sender_first_name,
    u2.last_name AS reply_sender_last_name
FROM
    discussion_messages m1
        LEFT JOIN
    discussion_messages m2 ON m1.parent_message_id = m2.message_id
        LEFT JOIN
    users u1 ON m1.sender_id = u1.user_id
        LEFT JOIN
    users u2 ON m2.sender_id = u2.user_id
WHERE
    m1.discussion_id = $1
ORDER BY
    m1.sent_at DESC
OFFSET $2
LIMIT $3;

-- name: ListDiscussionMessages3 :many
SELECT
    m1.*,
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