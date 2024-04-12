-- name: AddEducation :one
INSERT INTO educations (
    user_id,
    school,
    degree,
    field_of_study,
    start_date,
    end_date,
    current,
    grade,
    achievements,
    skills
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;