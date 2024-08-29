-- name: CreateFAQ :one
INSERT INTO faqs (
    question,
    answer
) VALUES (
    $1, $2
) RETURNING *;

-- name: ListFAQs :many
SELECT *
FROM faqs
ORDER BY created_at ASC;