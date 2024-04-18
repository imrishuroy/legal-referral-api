-- name: AddPrice :one
INSERT INTO pricing (
    user_id,
    service_type,
    per_hour_price,
    per_hearing_price,
    contingency_price,
    hybrid_price
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;