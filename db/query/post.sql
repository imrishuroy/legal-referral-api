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

-- name: DeletePost :exec
DELETE FROM posts
WHERE post_id = $1 AND owner_id = $2;

-- name: GetPost :one
SELECT
    posts.post_id,
    posts.owner_id,
    posts.content,
    posts.media,
    posts.post_type,
    posts.poll_id,
    posts.created_at
FROM posts
WHERE posts.post_id = $1;

-- name: SearchPosts :many
SELECT
    posts.post_id,
    posts.owner_id,
    users.first_name as owner_first_name,
    users.last_name as owner_last_name,
    users.avatar_url as owner_avatar_url,
    users.practice_area as owner_practice_area,
    posts.content,
    posts.media,
    posts.post_type,
    posts.poll_id,
    posts.created_at
FROM posts
JOIN users ON posts.owner_id = users.user_id
WHERE posts.content ILIKE '%' || @SearchQuery::text || '%'
ORDER BY posts.created_at DESC
LIMIT $1
OFFSET $2;

-- name: GetPostV2 :one
SELECT
    posts.post_id,
    posts.owner_id,
    users.first_name as owner_first_name,
    users.last_name as owner_last_name,
    users.avatar_url as owner_avatar_url,
    users.practice_area as owner_practice_area,
    posts.content,
    posts.media,
    posts.post_type,
    posts.poll_id,
    posts.created_at
FROM posts
JOIN users ON posts.owner_id = users.user_id
WHERE posts.post_id = $1;
