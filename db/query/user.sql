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
    address = $5,
    is_email_verified = $6,
    is_mobile_verified = $7,
    wizard_step = $8,
    wizard_completed = $9
WHERE
    id = $1
RETURNING *;

-- name: GetUserWizardStep :one
SELECT wizard_step
FROM users
WHERE id = $1;

-- name: UpdateUserWizardStep :one
UPDATE users
SET
    wizard_step = $2
WHERE
    id = $1
RETURNING *;

-- name: UpdateUserAboutYou :one
UPDATE users
SET
    first_name = $2,
    last_name = $3,
    address = $4
WHERE
    id = $1
RETURNING *;