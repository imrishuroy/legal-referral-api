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

