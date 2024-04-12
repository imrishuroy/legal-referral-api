-- name: AddSocial :one
INSERT INTO socials (
    user_id,
    platform_name,
    link_url
) VALUES (
    $1, $2, $3
) RETURNING *;