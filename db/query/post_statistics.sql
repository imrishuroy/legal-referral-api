-- name: IncrementLikes :exec
INSERT INTO post_statistics (post_id, likes, updated_at)
VALUES ($1, 1, CURRENT_TIMESTAMP)
ON CONFLICT (post_id)
    DO UPDATE SET likes = post_statistics.likes + 1,
    updated_at = CURRENT_TIMESTAMP;

-- name: IncrementComments :exec
INSERT INTO post_statistics (post_id, comments, updated_at)
VALUES ($1, 1, CURRENT_TIMESTAMP)
ON CONFLICT (post_id)
    DO UPDATE SET comments = post_statistics.comments + 1,
    updated_at = CURRENT_TIMESTAMP;

-- name: GetPostStats :one
SELECT * FROM post_statistics WHERE post_id = $1;

-- name: DecrementLikes :exec
UPDATE post_statistics SET likes = likes - 1 WHERE post_id = $1;