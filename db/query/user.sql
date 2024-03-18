-- name: CreateUser :one
INSERT INTO users (
    email,
    first_name,
    last_name,
    is_email_verified
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;