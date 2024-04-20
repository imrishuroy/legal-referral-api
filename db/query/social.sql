-- name: AddSocial :one
INSERT INTO socials (
    entity_id,
    entity_type,
    platform,
    link
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

