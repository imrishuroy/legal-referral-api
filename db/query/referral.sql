-- name: CreateReferral :one
INSERT INTO referrals (
    referred_user_id,
    referrer_user_id,
    title,
    preferred_practice_area,
    preferred_practice_location,
    case_description
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: ListActiveReferrals :many
SELECT * FROM referrals
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

FROM
    referrals r
        JOIN
    users u
    ON
        r.referred_user_id = u.user_id
WHERE
    r.referral_id = $1;

-- name: ListProposals :many
SELECT
    u.user_id AS referrer_user_id,
    u.first_name AS referrer_first_name,
    u.last_name AS referrer_last_name,
    u.practice_area AS referrer_practice_area,
    u.practice_location AS referrer_practice_location,
    u.avatar_url AS referrer_avatar_url,
    r.referral_id,
    r.title,
    r.preferred_practice_area,
    r.preferred_practice_location,
    r.case_description,
    r.created_at,
    r.updated_at
FROM
    referrals r
        JOIN
    users u
    ON
        r.referrer_user_id = u.user_id
WHERE
    r.referred_user_id = $1 AND status = 'active';

-- name: ChangeReferralStatus :one
UPDATE referrals
SET status = $2
WHERE referral_id = $1
RETURNING *;