-- name: PostToNewsFeed :exec
INSERT INTO news_feed (
    user_id,
    post_id
) VALUES (
    $1, $2
);

-- name: ListNewsFeed :many
SELECT
    nf.feed_id,
    sqlc.embed(users),
    sqlc.embed(posts),
    nf.created_at,
    COALESCE(likes_counts.likes_count, 0) AS likes_count,
    COALESCE(comments_counts.comments_count, 0) AS comments_count,
    CASE WHEN user_likes.like_id IS NOT NULL THEN true ELSE false END AS is_liked
FROM users
         JOIN news_feed nf ON nf.user_id = users.user_id
         JOIN posts ON nf.post_id = posts.post_id
         LEFT JOIN (
    SELECT
        post_id,
        COUNT(*) AS likes_count
    FROM likes
    WHERE type = 'post'
    GROUP BY post_id
) likes_counts ON nf.post_id = likes_counts.post_id
         LEFT JOIN (
    SELECT
        post_id,
        COUNT(*) AS comments_count
    FROM comments
    GROUP BY post_id
) comments_counts ON nf.post_id = comments_counts.post_id
LEFT JOIN (
    SELECT
        like_id,
        post_id
    FROM likes
    WHERE likes.user_id = $1 AND type = 'post'
) user_likes ON nf.post_id = user_likes.post_id
WHERE nf.user_id = $1
ORDER BY nf.created_at DESC
LIMIT $2
OFFSET $3;

-- name: ListNewsFeed2 :many
SELECT
    nf.feed_id,
    sqlc.embed(post_owner), -- Embed the post owner data
    sqlc.embed(posts),
    nf.created_at,
    COALESCE(likes_counts.likes_count, 0) AS likes_count,
    COALESCE(comments_counts.comments_count, 0) AS comments_count,
    CASE WHEN user_likes.like_id IS NOT NULL THEN true ELSE false END AS is_liked
FROM news_feed nf
         JOIN posts ON nf.post_id = posts.post_id
         JOIN users post_owner ON posts.owner_id = post_owner.user_id -- Join with post owner
         LEFT JOIN (
    SELECT
        post_id,
        COUNT(*) AS likes_count
    FROM likes
    WHERE type = 'post'
    GROUP BY post_id
) likes_counts ON nf.post_id = likes_counts.post_id
         LEFT JOIN (
    SELECT
        post_id,
        COUNT(*) AS comments_count
    FROM comments
    GROUP BY post_id
) comments_counts ON nf.post_id = comments_counts.post_id
         LEFT JOIN (
    SELECT
        like_id,
        post_id
    FROM likes
    WHERE likes.user_id = $1 AND type = 'post'
) user_likes ON nf.post_id = user_likes.post_id
WHERE nf.user_id = $1
ORDER BY nf.created_at DESC
LIMIT $2
OFFSET $3;






