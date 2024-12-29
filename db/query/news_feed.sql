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
    posts.post_id,
    posts.content,
    posts.media,
    posts.owner_id,
    posts.post_type,
    posts.poll_id,
    post_owner.first_name AS owner_first_name,
    post_owner.last_name AS owner_last_name,
    post_owner.avatar_url AS owner_avatar_url,
    post_owner.practice_area AS owner_practice_area,
    nf.created_at,
    COALESCE(post_stats.likes, 0) AS likes_count,
    COALESCE(post_stats.comments, 0) AS comments_count,
    EXISTS (
        SELECT 1
        FROM likes
        WHERE likes.user_id = $1 AND likes.post_id = nf.post_id AND likes.type = 'post'
    ) AS is_liked
FROM news_feed nf
         JOIN posts ON nf.post_id = posts.post_id
         JOIN users post_owner ON posts.owner_id = post_owner.user_id
         LEFT JOIN post_statistics post_stats ON nf.post_id = post_stats.post_id
WHERE nf.user_id = $1
ORDER BY nf.created_at DESC
LIMIT $2
OFFSET $3;

-- name: ListNewsFeedItems :many
SELECT *
FROM news_feed
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: ListNewsFeedV3 :many
SELECT * FROM news_feed
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: IgnoreFeed :exec
DELETE FROM news_feed
WHERE feed_id = $1 AND user_id = $2;

