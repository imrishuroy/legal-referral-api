-- name: CreatePost :one
INSERT INTO posts (
    owner_id,
    content,
    media,
    post_type,
    poll_id
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;


-- name: GetPostLikesCount :one
SELECT
    COUNT(*) AS likes_count
FROM likes
WHERE post_id = $1 AND type = 'post';

-- name: GetPostCommentsCount :one
SELECT
    COUNT(*) AS comments_count
FROM comments
WHERE post_id = $1;

-- name: GetPostLikesAndCommentsCount :one
SELECT
    COALESCE(likes_counts.likes_count, 0) AS likes_count,
    COALESCE(comments_counts.comments_count, 0) AS comments_count
FROM posts
            LEFT JOIN (
        SELECT
            post_id,
            COUNT(*) AS likes_count
        FROM likes
        WHERE type = 'post'
        GROUP BY post_id
    ) likes_counts ON posts.post_id = likes_counts.post_id
            LEFT JOIN (
        SELECT
            post_id,
            COUNT(*) AS comments_count
        FROM comments
        GROUP BY post_id
    ) comments_counts ON posts.post_id = comments_counts.post_id
    WHERE posts.post_id = $1;

-- name: GetPosIsLikedByCurrentUser :one
SELECT
    CASE WHEN like_id IS NOT NULL THEN true ELSE false END AS is_liked
FROM likes
WHERE post_id = $1 AND user_id = $2 AND type = 'post';