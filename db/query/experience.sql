-- name: AddExperience :one
INSERT INTO experiences (
    user_id,
    title,
    practice_area,
    company_name,
    practice_location,
    start_date,
    end_date,
    current,
    description,
    skills
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;