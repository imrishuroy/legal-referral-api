-- name: CreatePoll :one
INSERT INTO polls (
    owner_id,
    title,
    options,
    end_time
) VALUES (
    $1, $2, $3, $4
) RETURNING *;