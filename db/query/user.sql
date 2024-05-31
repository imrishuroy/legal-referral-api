-- name: CreateUser :one
INSERT INTO users (
    user_id,
    email,
    mobile,
    first_name,
    last_name,
    signup_method,
    email_verified,
    mobile_verified,
    avatar_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetUserById :one
SELECT * FROM users
WHERE user_id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserWizardStep :one
SELECT wizard_step
FROM users
WHERE user_id = $1;

-- name: UpdateUserWizardStep :one
UPDATE users
SET
    wizard_step = $2
WHERE
    user_id = $1
RETURNING *;

-- name: UpdateUserAvatarUrl :one
UPDATE users
SET
    avatar_url = $2
WHERE
    user_id = $1
RETURNING *;

-- name: MarkWizardCompleted :one
UPDATE users
SET
    wizard_completed = $2
WHERE
    user_id = $1
RETURNING *;

-- name: UpdateEmailVerificationStatus :one
UPDATE users
SET
    email_verified = $2
WHERE
    user_id = $1
RETURNING *;

-- name: UpdateMobileVerificationStatus :one
UPDATE users
SET
    mobile = $2,
    mobile_verified = $3
WHERE
    user_id = $1
RETURNING *;

-- name: SaveAboutYou :one
UPDATE users
SET
    address = $2,
    practice_area = $3,
    practice_location = $4,
    experience = $5,
    wizard_completed = $6
WHERE
    user_id = $1
RETURNING *;

-- name: UpdateUserInfo :one
UPDATE users
SET
    first_name = $2,
    last_name = $3,
    average_billing_per_client = $4,
    case_resolution_rate = $5,
    about = $6
WHERE
    user_id = $1
RETURNING *;

-- name: UpdateUserBannerImage :exec
UPDATE users
SET
    banner_url = $2
WHERE
    user_id = $1;

-- name: ListConnectedUsers :many
SELECT
    u.user_id,
    u.avatar_url,
    u.first_name,
    u.last_name
FROM
    connections c
        JOIN
    users u
    ON
        (c.recipient_id = u.user_id OR c.sender_id = u.user_id)
WHERE
    (c.sender_id = @user_id::text OR c.recipient_id = @user_id::text)
  AND u.user_id != @user_id::text
ORDER BY
    c.created_at
LIMIT $1
OFFSET $2;

-- name: ListUsers :many
SELECT
    user_id,
    first_name,
    last_name,
    avatar_url,
    join_date
FROM
    users
WHERE
    user_id != $1
ORDER BY
    join_date DESC
LIMIT $2
OFFSET $3;
