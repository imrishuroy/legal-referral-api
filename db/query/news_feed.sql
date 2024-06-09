-- name: PostToNewsFeed :exec
INSERT INTO news_feed (
    user_id,
    post_id
) VALUES (
    $1, $2
);

-- name: ListNewsFeed :many
SELECT
    u.user_id,
    u.first_name,
    u.last_name,
    u.avatar_url,
    p.post_id,
    p.title,
    p.content,
    p.media,
    p.post_type,
    p.created_at
FROM news_feed nf
    JOIN users u ON nf.user_id = u.user_id
    JOIN posts p ON nf.post_id = p.post_id
WHERE nf.user_id = $1
ORDER BY nf.created_at DESC;

-- name: ListNewsFeed2 :many
SELECT nf.feed_id, sqlc.embed(users), sqlc.embed(posts), nf.created_at
FROM users
         JOIN news_feed nf ON nf.user_id = users.user_id
         JOIN posts ON nf.post_id = posts.post_id
WHERE nf.user_id = $1
ORDER BY nf.created_at DESC;

