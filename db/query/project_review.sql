-- name: CreateProjectReview :one
INSERT INTO project_reviews (
    project_id,
    user_id,
    review,
    rating
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetProjectReview :one
SELECT *
FROM project_reviews
WHERE project_id = $1 AND user_id = $2;