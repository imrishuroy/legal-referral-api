-- name: ListChatRooms :many
SELECT
    cr.room_id,
    cr.user1_id,
    cr.user2_id,
    u2.first_name AS user2_first_name,
    u2.last_name AS user2_last_name,
    u2.avatar_url AS user2_avatar_url,
    cr.last_message_at,
    cr.created_at,
    m.message AS last_message,
    m.sent_at AS last_message_sent_at
FROM
    chat_rooms AS cr
        JOIN
    users AS u1 ON cr.user1_id = u1.user_id
        JOIN
    users AS u2 ON cr.user2_id = u2.user_id
        LEFT JOIN
    messages AS m ON cr.room_id = m.room_id
WHERE
    (cr.user1_id = $1 OR cr.user2_id = $1) AND
    m.message IS NOT NULL AND
    m.sent_at = (
        SELECT MAX(sent_at)
        FROM messages
        WHERE room_id = cr.room_id
    )
ORDER BY
    last_message_sent_at DESC;

-- name: CreateChatRoom :one
INSERT INTO chat_rooms (room_id, user1_id, user2_id)
VALUES ($1, $2, $3)
RETURNING
    room_id,
    user1_id,
    user2_id,
    (SELECT first_name FROM users WHERE user_id = $3) AS user2_first_name,
    (SELECT last_name FROM users WHERE user_id = $3) AS user2_last_name,
    (SELECT avatar_url FROM users WHERE user_id = $3) AS user2_avatar_url,
    last_message_at,
    created_at;

-- name: GetChatRoom :one
SELECT
    cr.room_id,
    cr.user1_id,
    cr.user2_id,
    u2.first_name AS user2_first_name,
    u2.last_name AS user2_last_name,
    u2.avatar_url AS user2_avatar_url,
    cr.last_message_at,
    cr.created_at
FROM
    chat_rooms AS cr
        JOIN
    users AS u2 ON cr.user2_id = u2.user_id
WHERE
    cr.room_id = $1;
