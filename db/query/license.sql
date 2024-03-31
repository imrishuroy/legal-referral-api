-- name: SaveLicense :one
INSERT INTO license (
    user_id,
    name,
    license_number,
    issue_date,
    issue_state
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;