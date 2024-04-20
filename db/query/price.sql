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

-- name: UpdatePrice :one
UPDATE pricing SET
    service_type = $2,
    per_hour_price = $3,
    per_hearing_price = $4,
    contingency_price = $5,
    hybrid_price = $6
WHERE
    price_id = $1
RETURNING *;