-- name: CreateReferral :one
INSERT INTO projects (
    referrer_user_id,
    title,
    preferred_practice_area,
    preferred_practice_location,
    case_description,
    status
) VALUES (
    $1, $2, $3, $4, $5, 'active'
) RETURNING *;

-- name: AddReferredUserToProject :one
INSERT INTO referral_users (
    project_id,
    referred_user_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: ListActiveReferrals :many
SELECT * FROM projects
WHERE referrer_user_id = @user_id::text AND (status = 'active' OR status = 'awarded')
ORDER BY created_at DESC;

-- name: ListReferredUsers :many
SELECT
    u.user_id,
    u.first_name,
    u.last_name,
    u.avatar_url,
    u.practice_area,
    u.practice_location,
    u.average_billing_per_client
FROM referral_users ru
    JOIN users u ON ru.referred_user_id = u.user_id
WHERE ru.project_id = $1;

-- name: ListActiveProposals :many
SELECT
    referrer.user_id AS user_id,
    referrer.first_name AS first_name,
    referrer.last_name AS last_name,
    referrer.practice_area AS practice_area,
    referrer.practice_location AS practice_location,
    referrer.avatar_url AS avatar_url,
    p.project_id,
    p.title,
    p.preferred_practice_area,
    p.preferred_practice_location,
    p.case_description,
    p.status,
    p.created_at

FROM
    projects p
        INNER JOIN
    referral_users ru ON p.project_id = ru.project_id
        INNER JOIN
    users referrer ON p.referrer_user_id = referrer.user_id
WHERE
    ru.referred_user_id = @user_id::text AND p.status = 'active'
ORDER BY
    p.created_at DESC;

-- name: GetProjectStatus :one
SELECT status FROM projects WHERE project_id = $1;

-- name: AwardProject :one
UPDATE projects
SET
    referred_user_id = $2,
    status = 'awarded'
WHERE project_id = $1
RETURNING *;

-- name: ListAwardedProjects :many
SELECT
    p.project_id,
    p.status,
    p.created_at,
    p.started_at,
    p.completed_at,
    p.title,
    p.case_description,
    referrer_user.user_id AS user_id,
    referrer_user.first_name AS first_name,
    referrer_user.last_name AS last_name,
    referrer_user.avatar_url AS avatar_url,
    referrer_user.practice_area AS practice_area
FROM
    projects p
        JOIN users referrer_user ON p.referrer_user_id = referrer_user.user_id
WHERE
    p.referred_user_id = @user_id::text
  AND p.status = 'awarded'
ORDER BY p.created_at DESC;

-- name: AcceptProject :one
UPDATE projects
SET
    status = 'accepted'
WHERE
    project_id = @project_id::int
  AND referred_user_id = @user_id::text
  AND status = 'awarded'
RETURNING *;
--
-- name: RejectProject :one
UPDATE projects
SET
    status = 'rejected'
WHERE
    project_id = @project_id::int
  AND referred_user_id = @user_id::text
RETURNING *;

-- name: StartProject :one
UPDATE projects
SET
    status = 'started',
    started_at = current_timestamp
WHERE
    project_id = @project_id::int
  AND referred_user_id = @user_id::text
RETURNING *;

-- name: ListReferrerActiveProjects :many
SELECT
    p.project_id,
    p.status,
    p.created_at,
    p.started_at,
    p.completed_at,
    p.title,
    p.case_description,
    referred_user.user_id AS user_id,
    referred_user.first_name AS first_name,
    referred_user.last_name AS last_name,
    referred_user.avatar_url AS avatar_url,
    referred_user.practice_area AS practice_area
FROM
    projects p
        JOIN users referred_user ON p.referred_user_id = referred_user.user_id
WHERE
    p.referrer_user_id = @user_id::text
  AND (p.status = 'started' OR p.status = 'accepted' OR p.status = 'complete_initiated')
  ORDER BY p.created_at DESC;
--
-- name: ListReferredActiveProjects :many
SELECT
    p.project_id,
    p.status,
    p.created_at,
    p.started_at,
    p.completed_at,
    p.title,
    p.case_description,
    referrer_user.user_id AS user_id,
    referrer_user.first_name AS first_name,
    referrer_user.last_name AS last_name,
    referrer_user.avatar_url AS avatar_url,
    referrer_user.practice_area AS practice_area
FROM
    projects p
        JOIN users referrer_user ON p.referrer_user_id = referrer_user.user_id
WHERE
    p.referred_user_id = @user_id::text
  AND (p.status = 'started' OR p.status = 'accepted' OR p.status = 'complete_initiated')
  ORDER BY p.created_at DESC;

-- name: InitiateCompleteProject :one
UPDATE projects
SET
    status = 'complete_initiated'
WHERE
    project_id = @project_id::int
  AND referred_user_id = @user_id::text
RETURNING *;

-- name: CompleteProject :one
UPDATE projects
SET
    status = 'completed',
    completed_at = current_timestamp
WHERE
    project_id = @project_id::int
  AND referrer_user_id = @user_id::text
RETURNING *;

-- name: CancelCompleteProjectInitiation :one
UPDATE projects
SET
    status = 'started'
WHERE
    project_id = @project_id::int
  AND referrer_user_id = @user_id::text
RETURNING *;

-- name: ListReferrerCompletedProjects :many
SELECT
    p.project_id,
    p.status,
    p.created_at,
    p.started_at,
    p.completed_at,
    p.title,
    p.case_description,
    p.preferred_practice_area,
    p.preferred_practice_location,
    referred_user.user_id AS user_id,
    referred_user.first_name AS first_name,
    referred_user.last_name AS last_name,
    referred_user.avatar_url AS avatar_url,
    referred_user.practice_area AS practice_area
FROM
    projects p
        JOIN users referred_user ON p.referred_user_id = referred_user.user_id
WHERE
    p.referrer_user_id = @user_id::text
  AND (p.status = 'completed')
ORDER BY p.created_at DESC;

-- name: ListReferredCompletedProjects :many
SELECT
    p.project_id,
    p.status,
    p.created_at,
    p.started_at,
    p.completed_at,
    p.title,
    p.case_description,
    p.preferred_practice_area,
    p.preferred_practice_location,
    referrer_user.user_id AS user_id,
    referrer_user.first_name AS first_name,
    referrer_user.last_name AS last_name,
    referrer_user.avatar_url AS avatar_url,
    referrer_user.practice_area AS practice_area
FROM
    projects p
        JOIN users referrer_user ON p.referrer_user_id = referrer_user.user_id
WHERE
    p.referred_user_id = @user_id::text
  AND (p.status = 'completed')
ORDER BY p.created_at DESC;

