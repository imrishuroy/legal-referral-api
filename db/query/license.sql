-- name: SaveLicense :one
INSERT INTO licenses (
    user_id,
    name,
    license_number,
    issue_date,
    issue_state
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: UploadLicense :one
UPDATE licenses
SET
    license_url = $1
WHERE
    user_id = $2
RETURNING *;