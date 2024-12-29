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
    posts.content,
    posts.media,
    posts.post_type,
    posts.poll_id,
    posts.created_at,
    posts.owner_id,
    post_owner.first_name AS owner_first_name,
    post_owner.last_name AS owner_last_name,
    post_owner.avatar_url AS owner_avatar_url,
    post_owner.practice_area AS owner_practice_area,
    COALESCE(post_stats.likes, 0) AS likes_count,
    COALESCE(post_stats.comments, 0) AS comments_count,
    EXISTS (
        SELECT 1
        FROM likes
        WHERE likes.user_id = posts.owner_id
          AND likes.post_id = posts.post_id
          AND likes.type = 'post'
    ) AS is_liked
FROM posts
         JOIN users post_owner ON posts.owner_id = post_owner.user_id
         LEFT JOIN post_statistics post_stats ON posts.post_id = post_stats.post_id
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

-- name: ListPosts :many
SELECT
    posts.post_id,
    posts.owner_id,
    posts.content,
    posts.media,
    posts.post_type,
    posts.poll_id,
    posts.created_at
FROM posts
WHERE posts.post_id = ANY(sqlc.slice('post_ids'));

-- name: PostsMetaData :many
SELECT
    posts.post_id,
    users.first_name as owner_first_name,
    users.last_name as owner_last_name,
    users.avatar_url as owner_avatar_url,
    users.practice_area as owner_practice_area,
    COALESCE(post_stats.likes, 0) AS likes_count,
    COALESCE(post_stats.comments, 0) AS comments_count,
    EXISTS (
        SELECT 1
        FROM likes
        WHERE likes.user_id = $1 AND likes.post_id = posts.post_id AND likes.type = 'post'
    ) AS is_liked
FROM posts
    LEFT JOIN post_statistics post_stats ON posts.post_id = post_stats.post_id
    JOIN users ON posts.owner_id = users.user_id
WHERE posts.post_id = ANY(sqlc.slice('post_ids'));

