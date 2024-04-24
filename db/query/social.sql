-- name: AddSocial :one
INSERT INTO socials (
    entity_id,
    entity_type,
    platform,
    link
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: UpdateSocial :one
UPDATE socials
SET
    platform = $2,
    link = $3
WHERE social_id = $1
RETURNING *;

-- name: ListSocials :many
SELECT * FROM socials
WHERE entity_id = $1 AND entity_type = $2;

-- name: DeleteSocial :exec
DELETE FROM socials
WHERE social_id = $1;
