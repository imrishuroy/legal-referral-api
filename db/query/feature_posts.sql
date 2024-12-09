-- name: FeaturePost :exec
INSERT INTO feature_posts (
    user_id,
    post_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: UnFeaturePost :exec
DELETE FROM feature_posts
WHERE
    user_id = $1 AND post_id = $2;

-- name: ListFeaturePosts :many
SELECT
    feature_posts.feature_post_id,
    sqlc.embed(posts),
    feature_posts.created_at
FROM
    feature_posts
JOIN
    posts ON feature_posts.post_id = posts.post_id
ORDER BY
    feature_posts.created_at DESC;