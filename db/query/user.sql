-- name: CreateUser :one
INSERT INTO users (
    id,
    email,
    first_name,
    last_name,
    sign_up_method,
    is_email_verified
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;