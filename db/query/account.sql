-- name: GetAccountInfo :one
SELECT
    user_id,
    first_name,
    last_name,
    avatar_url,
    practice_area
FROM users
WHERE user_id = $1;

-- name: GetUserRatingInfo :one
SELECT
    AVG(rating) AS average_rating,
    COUNT(*) AS attorneys
FROM reviews
WHERE user_id = $1;


SELECT
    user_id,
    AVG(rating) AS average_rating
FROM
    reviews
WHERE
    user_id = $1
GROUP BY
    user_id;


-- name: GetFollowersCount :one
SELECT COUNT(*) AS follower_count
FROM connection_invitations
WHERE recipient_id = $1
  AND status NOT IN ('rejected', 'cancelled');

-- name: GetConnectionsCount :one
SELECT COUNT(*) AS connection_count
FROM connections
WHERE sender_id = $1 OR recipient_id = $1;
