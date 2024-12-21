-- name: GetAccountInfo :one
SELECT
    users.user_id,
    users.first_name,
    users.last_name,
    users.avatar_url,
    users.practice_area,
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
WHERE users.user_id = $1;


-- name: GetUserRatingInfo :one
SELECT
    AVG(rating) AS average_rating,
    COUNT(*) AS attorneys
FROM reviews
WHERE user_id = $1;

-- name: GetFollowersCount :one
SELECT COUNT(*) AS follower_count
FROM connection_invitations
WHERE recipient_id = $1
  AND status NOT IN ('rejected', 'cancelled');

-- name: GetConnectionsCount :one
SELECT COUNT(*) AS connection_count
FROM connections
WHERE sender_id = $1 OR recipient_id = $1;
