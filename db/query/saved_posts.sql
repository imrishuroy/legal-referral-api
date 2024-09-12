-- name: SavePost :one
INSERT INTO saved_posts (
    user_id,
    post_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: ListSavedPosts :many
SELECT
    saved_posts.saved_post_id,
    sqlc.embed(posts),
    saved_posts.created_at
FROM
    saved_posts
JOIN
    posts ON saved_posts.post_id = posts.id
WHERE
    saved_posts.user_id = $1
ORDER BY
    saved_posts.created_at DESC;