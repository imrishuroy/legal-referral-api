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

-- name: ListEducations :many
SELECT * FROM educations WHERE user_id = $1;

-- name: UpdateEducation :one
UPDATE educations SET
    school = $2,
    degree = $3,
    field_of_study = $4,
    start_date = $5,
    end_date = $6,
    current = $7,
    grade = $8,
    achievements = $9,
    skills = $10
WHERE education_id = $1
RETURNING *;

-- name: DeleteEducation :exec
DELETE FROM educations
WHERE education_id = $1;