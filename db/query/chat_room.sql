-- name: ListChatRooms :many
SELECT
    cr.room_id,
    CASE
        WHEN cr.user1_id = $1 THEN u2.user_id
        ELSE u1.user_id
        END AS user_id,
    CASE
        WHEN cr.user1_id = $1 THEN u2.first_name
        ELSE u1.first_name
        END AS first_name,
    CASE
        WHEN cr.user1_id = $1 THEN u2.last_name
        ELSE u1.last_name
        END AS last_name,
    CASE
        WHEN cr.user1_id = $1 THEN u2.avatar_url
        ELSE u1.avatar_url
        END AS avatar_url,
    cr.last_message_at,
    cr.created_at,
    m.message AS last_message,
    m.sent_at AS last_message_sent_at
FROM
    chat_rooms cr
        JOIN
    users u1 ON cr.user1_id = u1.user_id
        JOIN
    users u2 ON cr.user2_id = u2.user_id
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
    user2_id AS user_id,
    (SELECT first_name FROM users WHERE user_id = $3) AS first_name,
    (SELECT last_name FROM users WHERE user_id = $3) AS last_name,
    (SELECT avatar_url FROM users WHERE user_id = $3) AS avatar_url,
    last_message_at,
    created_at;

-- name: GetChatRoom :one
SELECT
    cr.room_id,
    CASE
        WHEN cr.user1_id = $2 THEN u2.user_id
        ELSE u1.user_id
        END AS user_id,
    CASE
        WHEN cr.user1_id = $2 THEN u2.first_name
        ELSE u1.first_name
        END AS first_name,
    CASE
        WHEN cr.user1_id = $2 THEN u2.last_name
        ELSE u1.last_name
        END AS last_name,
    CASE
        WHEN cr.user1_id = $2 THEN u2.avatar_url
        ELSE u1.avatar_url
        END AS avatar_url,
    cr.last_message_at,
    cr.created_at
FROM
    chat_rooms AS cr
        JOIN
    users AS u2 ON cr.user2_id = u2.user_id
        JOIN
    users AS u1 ON cr.user1_id = u1.user_id

WHERE
    cr.room_id = $1 AND (cr.user1_id = $2 OR cr.user2_id = $2);

