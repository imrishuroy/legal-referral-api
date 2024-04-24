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
    pricing.hybrid_price

FROM users
LEFT JOIN pricing ON pricing.user_id = users.user_id
WHERE users.user_id = $1;

-- name: ToggleOpenToRefferal :exec
UPDATE users
SET open_to_referral = $2
WHERE user_id = $1;