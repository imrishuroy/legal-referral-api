-- name: AddFirm :one
INSERT INTO firms (
    name,
    owner_user_id,
    logo_url,
    org_type,
    website,
    location,
    about
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)RETURNING *;

-- name: ListFirms :many
SELECT * FROM firms
WHERE @query::text = '' OR name ILIKE '%' || @query || '%'
ORDER BY firm_id
LIMIT $1
OFFSET $2;

-- name: GetFirm :one
SELECT * FROM firms
WHERE firm_id = $1;

-- name: ListFirmsByOwner :many
SELECT * FROM firms
WHERE owner_user_id = $1;