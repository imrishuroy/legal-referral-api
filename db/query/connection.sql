-- name: SendConnection :one
INSERT INTO connection_invitations (
    sender_id,
    recipient_id
) VALUES ($1, $2)
RETURNING (id);

-- name: AcceptConnection :one
UPDATE connection_invitations
SET status = 1
WHERE id = $1 AND status = 0
RETURNING *;

-- name: RejectConnection :exec
UPDATE connection_invitations
SET status = 3
WHERE id = $1 AND status = 0
RETURNING *;

-- name: AddConnection :exec
INSERT INTO connections (sender_id, recipient_id)
    VALUES ($1, $2);

-- name: ListConnectionInvitations :many
SELECT ci.*,
       u.first_name,
       u.last_name,
       u.about,
       u.avatar_url
FROM connection_invitations ci
JOIN users u ON ci.sender_id = u.user_id
WHERE ci.recipient_id = $1 AND ci.status = 0
ORDER BY ci.created_at DESC
OFFSET $2
LIMIT $3;

-- name: ListConnections :many
SELECT
    ci.*,
    CASE
        WHEN u1.user_id = sqlc.arg(user_id) THEN u2.first_name
        ELSE u1.first_name
        END as first_name,
    CASE
        WHEN u1.user_id = sqlc.arg(user_id) THEN u2.last_name
        ELSE u1.last_name
        END as last_name,
    CASE
        WHEN u1.user_id = sqlc.arg(user_id) THEN u2.about
        ELSE u1.about
        END as about,
    CASE
        WHEN u1.user_id = sqlc.arg(user_id) THEN u2.avatar_url
        ELSE u1.avatar_url
        END as avatar_url
FROM connections ci
         JOIN users u1 ON ci.sender_id = u1.user_id
         JOIN users u2 ON ci.recipient_id = u2.user_id
WHERE sender_id = sqlc.arg(user_id)::text OR recipient_id = sqlc.arg(user_id)
ORDER BY ci.created_at DESC
OFFSET $1
LIMIT $2;
