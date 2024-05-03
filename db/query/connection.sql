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

-- name: AddConnection :exec
INSERT INTO connections (sender_id, recipient_id)
    VALUES ($1, $2);

-- name: RejectConnection :exec
UPDATE connection_invitations
    SET status = 3
    WHERE sender_id = $1 AND recipient_id = $2 AND status = 'pending';

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
SELECT ci.*,
       u.first_name,
       u.last_name,
       u.about,
       u.avatar_url
FROM connections ci
JOIN users u ON ci.sender_id = u.user_id
WHERE sender_id = sqlc.arg(user_id)::text OR recipient_id = sqlc.arg(user_id)
ORDER BY created_at DESC
OFFSET $1
LIMIT $2;
