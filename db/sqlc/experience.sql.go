// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: experience.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addExperience = `-- name: AddExperience :one
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
) RETURNING experience_id, user_id, title, practice_area, firm_id, practice_location, start_date, end_date, current, description, skills
`

type AddExperienceParams struct {
	UserID           string      `json:"user_id"`
	Title            string      `json:"title"`
	PracticeArea     string      `json:"practice_area"`
	FirmID           int64       `json:"firm_id"`
	PracticeLocation string      `json:"practice_location"`
	StartDate        pgtype.Date `json:"start_date"`
	EndDate          pgtype.Date `json:"end_date"`
	Current          bool        `json:"current"`
	Description      string      `json:"description"`
	Skills           []string    `json:"skills"`
}

func (q *Queries) AddExperience(ctx context.Context, arg AddExperienceParams) (Experience, error) {
	row := q.db.QueryRow(ctx, addExperience,
		arg.UserID,
		arg.Title,
		arg.PracticeArea,
		arg.FirmID,
		arg.PracticeLocation,
		arg.StartDate,
		arg.EndDate,
		arg.Current,
		arg.Description,
		arg.Skills,
	)
	var i Experience
	err := row.Scan(
		&i.ExperienceID,
		&i.UserID,
		&i.Title,
		&i.PracticeArea,
		&i.FirmID,
		&i.PracticeLocation,
		&i.StartDate,
		&i.EndDate,
		&i.Current,
		&i.Description,
		&i.Skills,
	)
	return i, err
}

const deleteExperience = `-- name: DeleteExperience :exec
DELETE FROM experiences
WHERE experience_id = $1
`

func (q *Queries) DeleteExperience(ctx context.Context, experienceID int64) error {
	_, err := q.db.Exec(ctx, deleteExperience, experienceID)
	return err
}

const listExperiences = `-- name: ListExperiences :many
SELECT experiences.experience_id, experiences.user_id, experiences.title, experiences.practice_area, experiences.firm_id, experiences.practice_location, experiences.start_date, experiences.end_date, experiences.current, experiences.description, experiences.skills, firms.firm_id, firms.name, firms.logo_url, firms.org_type, firms.website, firms.location, firms.about
FROM experiences
JOIN firms ON experiences.firm_id = firms.firm_id
WHERE user_id = $1
ORDER BY COALESCE(experiences.end_date, CURRENT_DATE) DESC
`

type ListExperiencesRow struct {
	Experience Experience `json:"experience"`
	Firm       Firm       `json:"firm"`
}

func (q *Queries) ListExperiences(ctx context.Context, userID string) ([]ListExperiencesRow, error) {
	rows, err := q.db.Query(ctx, listExperiences, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListExperiencesRow{}
	for rows.Next() {
		var i ListExperiencesRow
		if err := rows.Scan(
			&i.Experience.ExperienceID,
			&i.Experience.UserID,
			&i.Experience.Title,
			&i.Experience.PracticeArea,
			&i.Experience.FirmID,
			&i.Experience.PracticeLocation,
			&i.Experience.StartDate,
			&i.Experience.EndDate,
			&i.Experience.Current,
			&i.Experience.Description,
			&i.Experience.Skills,
			&i.Firm.FirmID,
			&i.Firm.Name,
			&i.Firm.LogoUrl,
			&i.Firm.OrgType,
			&i.Firm.Website,
			&i.Firm.Location,
			&i.Firm.About,
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

const updateExperience = `-- name: UpdateExperience :one
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
RETURNING experience_id, user_id, title, practice_area, firm_id, practice_location, start_date, end_date, current, description, skills
`

type UpdateExperienceParams struct {
	ExperienceID     int64       `json:"experience_id"`
	Title            string      `json:"title"`
	PracticeArea     string      `json:"practice_area"`
	FirmID           int64       `json:"firm_id"`
	PracticeLocation string      `json:"practice_location"`
	StartDate        pgtype.Date `json:"start_date"`
	EndDate          pgtype.Date `json:"end_date"`
	Current          bool        `json:"current"`
	Description      string      `json:"description"`
	Skills           []string    `json:"skills"`
}

func (q *Queries) UpdateExperience(ctx context.Context, arg UpdateExperienceParams) (Experience, error) {
	row := q.db.QueryRow(ctx, updateExperience,
		arg.ExperienceID,
		arg.Title,
		arg.PracticeArea,
		arg.FirmID,
		arg.PracticeLocation,
		arg.StartDate,
		arg.EndDate,
		arg.Current,
		arg.Description,
		arg.Skills,
	)
	var i Experience
	err := row.Scan(
		&i.ExperienceID,
		&i.UserID,
		&i.Title,
		&i.PracticeArea,
		&i.FirmID,
		&i.PracticeLocation,
		&i.StartDate,
		&i.EndDate,
		&i.Current,
		&i.Description,
		&i.Skills,
	)
	return i, err
}
