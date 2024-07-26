// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: recommendation.sql

package db

import (
	"context"
)

const cancelRecommendation = `-- name: CancelRecommendation :exec
INSERT INTO canceled_recommendations (
    user_id,
    recommended_user_id
) VALUES (
    $1, $2
)
`

type CancelRecommendationParams struct {
	UserID            string `json:"user_id"`
	RecommendedUserID string `json:"recommended_user_id"`
}

func (q *Queries) CancelRecommendation(ctx context.Context, arg CancelRecommendationParams) error {
	_, err := q.db.Exec(ctx, cancelRecommendation, arg.UserID, arg.RecommendedUserID)
	return err
}

const listRecommendations = `-- name: ListRecommendations :many
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
                u.experience,
                u.practice_location
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
    WHERE (sender_id = $1 AND recipient_id = u.user_id) AND status = 'pending'
)
  AND NOT EXISTS (
    SELECT 1
    FROM canceled_recommendations
    WHERE user_id = $1 AND recommended_user_id = u.user_id
)
OFFSET $2
LIMIT $3
`

type ListRecommendationsParams struct {
	UserID string `json:"user_id"`
	Offset int32  `json:"offset"`
	Limit  int32  `json:"limit"`
}

type ListRecommendationsRow struct {
	UserID           string  `json:"user_id"`
	FirstName        string  `json:"first_name"`
	LastName         string  `json:"last_name"`
	About            *string `json:"about"`
	AvatarUrl        *string `json:"avatar_url"`
	PracticeArea     *string `json:"practice_area"`
	Experience       *string `json:"experience"`
	PracticeLocation *string `json:"practice_location"`
}

func (q *Queries) ListRecommendations(ctx context.Context, arg ListRecommendationsParams) ([]ListRecommendationsRow, error) {
	rows, err := q.db.Query(ctx, listRecommendations, arg.UserID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListRecommendationsRow{}
	for rows.Next() {
		var i ListRecommendationsRow
		if err := rows.Scan(
			&i.UserID,
			&i.FirstName,
			&i.LastName,
			&i.About,
			&i.AvatarUrl,
			&i.PracticeArea,
			&i.Experience,
			&i.PracticeLocation,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listRecommendations2 = `-- name: ListRecommendations2 :many
SELECT
    u.user_id,
    u.first_name,
    u.last_name,
    u.about,
    u.avatar_url,
    u.practice_area,
    u.experience,
    u.practice_location
FROM
    users u
WHERE
    u.user_id != $1
  AND u.user_id NOT IN (
    SELECT recipient_id FROM connection_invitations WHERE sender_id = $1
    UNION
    SELECT sender_id FROM connection_invitations WHERE recipient_id = $1
    UNION
    SELECT recipient_id FROM connections WHERE sender_id = $1
    UNION
    SELECT sender_id FROM connections WHERE recipient_id = $1
    UNION
    SELECT recommended_user_id FROM canceled_recommendations WHERE user_id = $1
)
LIMIT $2 OFFSET $3
`

type ListRecommendations2Params struct {
	UserID string `json:"user_id"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

type ListRecommendations2Row struct {
	UserID           string  `json:"user_id"`
	FirstName        string  `json:"first_name"`
	LastName         string  `json:"last_name"`
	About            *string `json:"about"`
	AvatarUrl        *string `json:"avatar_url"`
	PracticeArea     *string `json:"practice_area"`
	Experience       *string `json:"experience"`
	PracticeLocation *string `json:"practice_location"`
}

func (q *Queries) ListRecommendations2(ctx context.Context, arg ListRecommendations2Params) ([]ListRecommendations2Row, error) {
	rows, err := q.db.Query(ctx, listRecommendations2, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListRecommendations2Row{}
	for rows.Next() {
		var i ListRecommendations2Row
		if err := rows.Scan(
			&i.UserID,
			&i.FirstName,
			&i.LastName,
			&i.About,
			&i.AvatarUrl,
			&i.PracticeArea,
			&i.Experience,
			&i.PracticeLocation,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
