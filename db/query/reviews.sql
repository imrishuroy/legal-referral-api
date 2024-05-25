-- name: AddReview :one
INSERT INTO reviews (
    user_id,
    reviewer_id,
    review,
    rating
) VALUES (
   $1, $2, $3, $4
) RETURNING *;
