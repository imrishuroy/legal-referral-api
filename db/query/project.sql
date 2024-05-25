-- name: AwardProject :one
INSERT INTO projects (
    referred_user_id,
    referrer_user_id,
    referral_id
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: ListReferrerActiveProjects :many
SELECT
    p.project_id,
    p.status,
    p.created_at,
    p.started_at,
    p.completed_at,
    r.referral_id,
    r.title,
    r.case_description,
    referred_user.user_id AS user_id,
    referred_user.first_name AS first_name,
    referred_user.last_name AS last_name,
    referred_user.avatar_url AS avatar_url,
    referred_user.practice_area AS practice_area
FROM
    projects p
        JOIN users referrer_user ON p.referrer_user_id = referrer_user.user_id
        JOIN users referred_user ON p.referred_user_id = referred_user.user_id
        JOIN referrals r ON p.referral_id = r.referral_id
WHERE
    p.referrer_user_id = @user_id::text
  AND (p.status = 'started' OR p.status = 'accepted' OR p.status = 'complete_initiated')
  ORDER BY p.created_at DESC;

-- name: ListReferredActiveProjects :many
SELECT
    p.project_id,
    p.status,
    p.created_at,
    p.started_at,
    p.completed_at,
    r.referral_id,
    r.title,
    r.case_description,
    referrer_user.user_id AS user_id,
    referrer_user.first_name AS first_name,
    referrer_user.last_name AS last_name,
    referrer_user.avatar_url AS avatar_url,
    referrer_user.practice_area AS practice_area
FROM
    projects p
        JOIN users referrer_user ON p.referrer_user_id = referrer_user.user_id
        JOIN users referred_user ON p.referred_user_id = referred_user.user_id
        JOIN referrals r ON p.referral_id = r.referral_id
WHERE
    p.referred_user_id = @user_id::text
  AND (p.status = 'started' OR p.status = 'accepted')
  ORDER BY p.created_at DESC;

-- name: ListAwardedProjects :many
SELECT
    p.project_id,
    p.status,
    p.created_at,
    p.started_at,
    p.completed_at,
    r.referral_id,
    r.title,
    r.case_description,
    referrer_user.user_id,
    referrer_user.first_name,
    referrer_user.last_name,
    referrer_user.avatar_url AS avatar_url,
    referrer_user.practice_area AS practice_area
FROM
    projects p
        JOIN users referrer_user ON p.referrer_user_id = referrer_user.user_id
        JOIN referrals r ON p.referral_id = r.referral_id
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

-- name: CompleteProject :one
UPDATE projects
SET
    status = 'completed',
    completed_at = current_timestamp
WHERE
    project_id = @project_id::int
  AND referrer_user_id = @user_id::text
RETURNING *;

-- name: InitiateCompleteProject :one
UPDATE projects
SET
    status = 'complete_initiated'
WHERE
    project_id = @project_id::int
  AND referred_user_id = @user_id::text
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
    r.referral_id,
    r.title,
    r.case_description,
    r.preferred_practice_area,
    r.preferred_practice_location,
    referred_user.user_id AS user_id,
    referred_user.first_name AS first_name,
    referred_user.last_name AS last_name,
    referred_user.avatar_url AS avatar_url,
    referred_user.practice_area AS practice_area
FROM
    projects p
        JOIN users referrer_user ON p.referrer_user_id = referrer_user.user_id
        JOIN users referred_user ON p.referred_user_id = referred_user.user_id
        JOIN referrals r ON p.referral_id = r.referral_id
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
    r.referral_id,
    r.title,
    r.case_description,
    r.preferred_practice_area,
    r.preferred_practice_location,
    referrer_user.user_id AS user_id,
    referrer_user.first_name AS first_name,
    referrer_user.last_name AS last_name,
    referrer_user.avatar_url AS avatar_url,
    referrer_user.practice_area AS practice_area
FROM
    projects p
        JOIN users referrer_user ON p.referrer_user_id = referrer_user.user_id
        JOIN users referred_user ON p.referred_user_id = referred_user.user_id
        JOIN referrals r ON p.referral_id = r.referral_id
WHERE
    p.referred_user_id = @user_id::text
  AND (p.status = 'completed')
ORDER BY p.created_at DESC;


