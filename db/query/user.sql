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

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
    first_name = $2,
    last_name = $3,
    mobile = $4,
    is_email_verified = $5,
    is_mobile_verified = $6,
    wizard_step = $7,
    wizard_completed = $8
WHERE
    id = $1
RETURNING *;
