-- name: LikePost :exec
INSERT INTO likes (
    user_id,
    post_id,
    type
) VALUES (
    $1, $2, 'post'
);

-- name: UnlikePost :exec
DELETE FROM likes
WHERE user_id = $1 AND post_id = $2 AND type = 'post';


-- name: ListPostLikes :many
SELECT
    user_id
FROM likes
WHERE post_id = $1 AND type = 'post';

-- name: CheckPostLike :one
SELECT EXISTS(
    SELECT 1
    FROM likes
    WHERE user_id = $1 AND post_id = $2 AND type = 'post'
) as exists;

-- name: ListPostLikedUsers :many
SELECT
    users.user_id,
    first_name,
    last_name,
    avatar_url
FROM likes
    JOIN users ON likes.user_id = users.user_id
    WHERE post_id = $1 AND type = 'post';

-- name: ListPostLikedUsers2 :many
SELECT
    users.user_id,
    first_name,
    last_name,
    avatar_url
FROM likes
         JOIN users ON likes.user_id = users.user_id
WHERE post_id = $1 AND type = 'post'
ORDER BY
    CASE
        WHEN users.user_id = $2 THEN 0
        ELSE 1
        END,
    users.user_id;

-- name: LikeComment :exec
INSERT INTO likes (
    user_id,
    comment_id,
    type
) VALUES (
    $1, $2, 'comment'
);

-- name: UnlikeComment :exec
DELETE FROM likes
WHERE user_id = $1 AND comment_id = $2 AND type = 'comment';

