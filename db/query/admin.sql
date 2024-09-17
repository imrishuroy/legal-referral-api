-- name: ListLicenseVerifiedUsers :many
SELECT
    u.user_id,
    u.first_name,
    u.last_name,
    u.avatar_url,
    u.practice_location,
    u.join_date,
    l.license_id,
    l.license_number,
    l.name AS license_name,
    l.issue_date,
    l.issue_state,
    l.license_url
FROM
    users u
        LEFT JOIN
    licenses l ON u.user_id = l.user_id
WHERE
    u.license_verified = true
ORDER BY
    u.join_date DESC
LIMIT $1
OFFSET $2;


-- name: ListLicenseUnVerifiedUsers :many
SELECT
    u.user_id,
    u.first_name,
    u.last_name,
    u.avatar_url,
    u.practice_location,
    u.join_date,
    l.license_id,
    l.license_number,
    l.name AS license_name,
    l.issue_date,
    l.issue_state,
    l.license_url
FROM
    users u
        LEFT JOIN
    licenses l ON u.user_id = l.user_id
WHERE
     u.license_verified = false
--     AND u.license_rejected = false
ORDER BY
    u.join_date DESC
LIMIT $1
OFFSET $2;



-- name: ListAttorneys :many
WITH conn AS (
    SELECT
        user_id,
        COUNT(*) AS total_connections
    FROM (
             SELECT sender_id AS user_id FROM connections
             UNION ALL
             SELECT recipient_id AS user_id FROM connections
         ) AS all_connections
    GROUP BY user_id
)
SELECT
    u.user_id,
    u.first_name,
    u.last_name,
    u.practice_area,
    u.practice_location,
    p.price_id,
    p.service_type,
    p.per_hour_price,
    p.per_hearing_price,
    p.contingency_price,
    p.hybrid_price,
    COALESCE(conn.total_connections, 0) AS total_connections
FROM
    users u
        LEFT JOIN conn ON u.user_id = conn.user_id
        LEFT JOIN pricing p ON u.user_id = p.user_id
ORDER BY
    u.user_id
LIMIT $1
OFFSET $2;

-- lawyers

-- name: ListLawyers :many
SELECT
    u.user_id,
    u.first_name,
    u.last_name,
    u.avatar_url,
    COUNT(r.referral_user_id) AS referral_count
FROM
    users u
        LEFT JOIN
    referral_users r ON u.user_id = r.referred_user_id
GROUP BY
    u.user_id, u.first_name, u.last_name
ORDER BY
    u.user_id;

-- name: ListAllReferralProjects :many
SELECT
    project_id,
    title,
    preferred_practice_area,
    preferred_practice_location,
    case_description,
    referrer_user_id,
    referred_user_id,
    status,
    created_at,
    started_at,
    completed_at
FROM
    projects
WHERE
    referrer_user_id = $1;

-- name: ListCompletedReferralProjects :many
SELECT
    project_id,
    title,
    preferred_practice_area,
    preferred_practice_location,
    case_description,
    referrer_user_id,
    referred_user_id,
    status,
    created_at,
    started_at,
    completed_at
FROM
    projects
WHERE
    status = 'completed' AND referrer_user_id = $1;

-- name: ListActiveReferralProjects :many
SELECT
    project_id,
    title,
    preferred_practice_area,
    preferred_practice_location,
    case_description,
    referrer_user_id,
    referred_user_id,
    status,
    created_at,
    started_at,
    completed_at
FROM
    projects
WHERE
    status = 'active' AND referrer_user_id = $1;