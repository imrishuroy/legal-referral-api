-- name: CreatePost :one
INSERT INTO posts (
    owner_id,
    title,
    content,
    media,
    post_type,
    poll_id
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

