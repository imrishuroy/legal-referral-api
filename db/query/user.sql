-- name: CreateUser :one
INSERT INTO users (
    id,
    first_name,
    last_name,
    mobile_number,
    email,
    bar_licence_no,
    practicing_field,
    experience
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;