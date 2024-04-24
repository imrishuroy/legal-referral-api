-- name: AddFirm :one
INSERT INTO firms (
    name,
    logo_url,
    org_type,
    website,
    location,
    about
) VALUES (
    $1, $2, $3, $4, $5, $6
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