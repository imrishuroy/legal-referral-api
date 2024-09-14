-- name: SavePost :exec
INSERT INTO saved_posts (
    user_id,
    post_id
) VALUES (
    $1, $2
);

-- name: UnsavePost :exec
DELETE FROM saved_posts
WHERE
    saved_post_id = $1;

-- name: ListSavedPosts :many
SELECT
    saved_posts.saved_post_id,
    sqlc.embed(posts),
    sqlc.embed(users),
    saved_posts.created_at
FROM
    saved_posts
JOIN
    posts ON saved_posts.post_id = posts.post_id
JOIN
    users ON posts.owner_id = users.user_id
WHERE
    saved_posts.user_id = $1
ORDER BY
    saved_posts.created_at DESC;