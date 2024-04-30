-- name: ListRecommendations :many
WITH TargetUserExperiences AS (
    SELECT DISTINCT user_id, practice_area, practice_location, skills
    FROM experiences
    WHERE user_id = $1
),
TargetUserEducations AS (
    SELECT DISTINCT user_id, school, degree, field_of_study, skills
    FROM educations
    WHERE user_id = $1
)
SELECT DISTINCT u.user_id,
                u.first_name,
                u.last_name,
                u.about,
                u.avatar_url,
                u.practice_area,
                u.experience
FROM users u
LEFT JOIN TargetUserExperiences tue ON u.practice_area = tue.practice_area OR u.user_id = tue.user_id
LEFT JOIN TargetUserEducations tue2 ON u.practice_area = tue2.field_of_study OR u.user_id = tue2.user_id
WHERE u.user_id <> $1
OFFSET $2
LIMIT $3;

-- name: ListRecommendations2 :many
WITH TargetUserExperiences AS (
    SELECT DISTINCT user_id, practice_area, practice_location, skills
    FROM experiences
    WHERE user_id = $1
),
     TargetUserEducations AS (
         SELECT DISTINCT user_id, school, degree, field_of_study, skills
         FROM educations
         WHERE user_id = $1
     ),
     UserSkills AS (
         SELECT DISTINCT user_id, unnest(skills) AS skill
         FROM experiences
         WHERE user_id = $1
         UNION
         SELECT DISTINCT user_id, unnest(skills) AS skill
         FROM educations
         WHERE user_id = $1
     )
SELECT DISTINCT u.user_id,
                u.first_name,
                u.last_name,
                u.about,
                u.avatar_url,
                u.practice_area,
                u.experience
FROM users u
         LEFT JOIN TargetUserExperiences tue ON u.practice_area = tue.practice_area OR u.practice_location = tue.practice_location OR u.user_id = tue.user_id
         LEFT JOIN TargetUserEducations tue2 ON u.user_id = tue2.user_id
         LEFT JOIN UserSkills us ON u.user_id = us.user_id
WHERE u.user_id <> $1
  AND (tue.user_id IS NOT NULL OR tue2.user_id IS NOT NULL OR us.skill IS NOT NULL)
  AND NOT EXISTS (
    SELECT 1
    FROM connections
    WHERE (sender_id = $1 AND recipient_id = u.user_id)
       OR (sender_id = u.user_id AND recipient_id = $1)
)
  AND NOT EXISTS (
    SELECT 1
    FROM connection_invitations
    WHERE sender_id = $1 AND recipient_id = u.user_id
)
  AND NOT EXISTS (
    SELECT 1
    FROM canceled_recommendations
    WHERE user_id = $1 AND recommended_user_id = u.user_id
)
OFFSET $2
LIMIT $3;

-- name: CancelRecommendation :exec
INSERT INTO canceled_recommendations (
    user_id,
    recommended_user_id
) VALUES (
    $1, $2
);


