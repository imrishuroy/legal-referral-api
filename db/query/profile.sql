-- name: FetchUserProfile :many
SELECT sqlc.embed(users), sqlc.embed(pricing)
FROM users
JOIN pricing ON pricing.user_id = users.user_id
WHERE users.user_id = $1;
