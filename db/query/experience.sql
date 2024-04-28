-- name: AddExperience :one
INSERT INTO experiences (
    user_id,
    title,
    practice_area,
    firm_id,
    practice_location,
    start_date,
    end_date,
    current,
    description,
    skills
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- -- name: ListExperiences :many
-- SELECT sqlc.embed(experiences), sqlc.embed(firms)
-- FROM experiences
-- JOIN firms ON experiences.firm_id = firms.firm_id
-- WHERE user_id = $1;

-- name: ListExperiences :many
SELECT sqlc.embed(experiences), sqlc.embed(firms)
FROM experiences
JOIN firms ON experiences.firm_id = firms.firm_id
WHERE user_id = $1
ORDER BY COALESCE(experiences.end_date, CURRENT_DATE) DESC;


--- name: ListExperiences3 :many
-- SELECT
--     experiences.experience_id,
--     experiences.user_id,
--     experiences.title,
--     experiences.practice_area,
--     experiences.firm_id,
--     experiences.practice_location,
--     experiences.start_date,
--     experiences.end_date,
--     experiences.current,
--     experiences.description,
--     experiences.skills,
--     firms.firm_id,
--     firms.name,
--     firms.logo_url,
--     firms.org_type,
--     firms.website,
--     firms.location
-- FROM experiences
-- JOIN firms ON experiences.firm_id = firms.firm_id
-- WHERE user_id = $1;

-- name: UpdateExperience :one
UPDATE experiences
SET
    title = $2,
    practice_area = $3,
    firm_id = $4,
    practice_location = $5,
    start_date = $6,
    end_date = $7,
    current = $8,
    description = $9,
    skills = $10
WHERE experience_id = $1
RETURNING *;

-- name: DeleteExperience :exec
DELETE FROM experiences
WHERE experience_id = $1;
