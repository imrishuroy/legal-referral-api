-- name: SaveExperience :one
INSERT INTO experiences (
    user_id,
    practice_area,
    practice_location,
    experience
) VALUES (
    $1, $2, $3, $4
) RETURNING *;