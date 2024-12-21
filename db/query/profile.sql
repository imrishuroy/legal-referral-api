-- name: UpdateUserAvatar :exec
UPDATE users
SET avatar_url = $2
WHERE user_id = $1;

-- name: FetchUserProfile :one

SELECT
    users.user_id,
    users.first_name,
    users.last_name,
    users.practice_area,
    users.avatar_url,
    users.banner_url,
    users.average_billing_per_client,
    users.case_resolution_rate,
    users.open_to_referral,
    users.about,
    pricing.price_id,
    pricing.service_type,
    pricing.per_hour_price,
    pricing.per_hearing_price,
    pricing.contingency_price,
    pricing.hybrid_price,
    COALESCE((SELECT AVG(rating) FROM reviews WHERE user_id = users.user_id), 0.0) AS average_rating,
    COALESCE((SELECT COUNT(*) FROM reviews WHERE user_id = users.user_id), 0) AS attorneys,
    COALESCE((SELECT COUNT(*)
              FROM connection_invitations
              WHERE recipient_id = users.user_id
                AND status NOT IN ('rejected', 'cancelled')), 0) AS followers_count,
    COALESCE((SELECT COUNT(*)
              FROM connections
              WHERE sender_id = users.user_id OR recipient_id = users.user_id), 0) AS connections_count
FROM users
         LEFT JOIN pricing ON users.user_id = pricing.user_id
WHERE users.user_id = $1;




-- name: ToggleOpenToRefferal :exec
UPDATE users
SET open_to_referral = $2
WHERE user_id = $1;