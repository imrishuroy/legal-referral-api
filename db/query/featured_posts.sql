-- name: FeaturePost :exec
INSERT INTO featured_posts (
    user_id,
    post_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: UnFeaturePost :exec
DELETE FROM featured_posts
WHERE
    user_id = $1 AND post_id = $2;

-- name: ListFeaturedPosts :many
SELECT
    featured_posts.feature_post_id,
    sqlc.embed(posts),
    featured_posts.created_at
FROM
    featured_posts
JOIN
    posts ON featured_posts.post_id = posts.post_id
ORDER BY
    featured_posts.created_at DESC;

-- name: IsPostFeatured :one
SELECT
    CASE WHEN post_id IS NOT NULL THEN true ELSE false END AS is_featured
FROM featured_posts
WHERE post_id = $1;