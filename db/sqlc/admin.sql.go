// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: admin.sql

package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const listActiveReferralProjects = `-- name: ListActiveReferralProjects :many
SELECT
    project_id,
    title,
    preferred_practice_area,
    preferred_practice_location,
    case_description,
    referrer_user_id,
    referred_user_id,
    status,
    created_at,
    started_at,
    completed_at
FROM
    projects
WHERE
    status = 'active' AND referrer_user_id = $1
`

func (q *Queries) ListActiveReferralProjects(ctx context.Context, referrerUserID string) ([]Project, error) {
	rows, err := q.db.Query(ctx, listActiveReferralProjects, referrerUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Project{}
	for rows.Next() {
		var i Project
		if err := rows.Scan(
			&i.ProjectID,
			&i.Title,
			&i.PreferredPracticeArea,
			&i.PreferredPracticeLocation,
			&i.CaseDescription,
			&i.ReferrerUserID,
			&i.ReferredUserID,
			&i.Status,
			&i.CreatedAt,
			&i.StartedAt,
			&i.CompletedAt,
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

const listAllReferralProjects = `-- name: ListAllReferralProjects :many
SELECT
    project_id,
    title,
    preferred_practice_area,
    preferred_practice_location,
    case_description,
    referrer_user_id,
    referred_user_id,
    status,
    created_at,
    started_at,
    completed_at
FROM
    projects
WHERE
    referrer_user_id = $1
`

func (q *Queries) ListAllReferralProjects(ctx context.Context, referrerUserID string) ([]Project, error) {
	rows, err := q.db.Query(ctx, listAllReferralProjects, referrerUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Project{}
	for rows.Next() {
		var i Project
		if err := rows.Scan(
			&i.ProjectID,
			&i.Title,
			&i.PreferredPracticeArea,
			&i.PreferredPracticeLocation,
			&i.CaseDescription,
			&i.ReferrerUserID,
			&i.ReferredUserID,
			&i.Status,
			&i.CreatedAt,
			&i.StartedAt,
			&i.CompletedAt,
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

const listAttorneys = `-- name: ListAttorneys :many
WITH conn AS (
    SELECT
        user_id,
        COUNT(*) AS total_connections
    FROM (
             SELECT sender_id AS user_id FROM connections
             UNION ALL
             SELECT recipient_id AS user_id FROM connections
         ) AS all_connections
    GROUP BY user_id
)
SELECT
    u.user_id,
    u.first_name,
    u.last_name,
    u.practice_area,
    u.practice_location,
    p.price_id,
    p.service_type,
    p.per_hour_price,
    p.per_hearing_price,
    p.contingency_price,
    p.hybrid_price,
    COALESCE(conn.total_connections, 0) AS total_connections
FROM
    users u
        LEFT JOIN conn ON u.user_id = conn.user_id
        LEFT JOIN pricing p ON u.user_id = p.user_id
ORDER BY
    u.user_id
LIMIT $1
OFFSET $2
`

type ListAttorneysParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListAttorneysRow struct {
	UserID           string         `json:"user_id"`
	FirstName        string         `json:"first_name"`
	LastName         string         `json:"last_name"`
	PracticeArea     *string        `json:"practice_area"`
	PracticeLocation *string        `json:"practice_location"`
	PriceID          *int64         `json:"price_id"`
	ServiceType      *string        `json:"service_type"`
	PerHourPrice     pgtype.Numeric `json:"per_hour_price"`
	PerHearingPrice  pgtype.Numeric `json:"per_hearing_price"`
	ContingencyPrice *string        `json:"contingency_price"`
	HybridPrice      *string        `json:"hybrid_price"`
	TotalConnections int64          `json:"total_connections"`
}

func (q *Queries) ListAttorneys(ctx context.Context, arg ListAttorneysParams) ([]ListAttorneysRow, error) {
	rows, err := q.db.Query(ctx, listAttorneys, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListAttorneysRow{}
	for rows.Next() {
		var i ListAttorneysRow
		if err := rows.Scan(
			&i.UserID,
			&i.FirstName,
			&i.LastName,
			&i.PracticeArea,
			&i.PracticeLocation,
			&i.PriceID,
			&i.ServiceType,
			&i.PerHourPrice,
			&i.PerHearingPrice,
			&i.ContingencyPrice,
			&i.HybridPrice,
			&i.TotalConnections,
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

const listCompletedReferralProjects = `-- name: ListCompletedReferralProjects :many
SELECT
    project_id,
    title,
    preferred_practice_area,
    preferred_practice_location,
    case_description,
    referrer_user_id,
    referred_user_id,
    status,
    created_at,
    started_at,
    completed_at
FROM
    projects
WHERE
    status = 'completed' AND referrer_user_id = $1
`

func (q *Queries) ListCompletedReferralProjects(ctx context.Context, referrerUserID string) ([]Project, error) {
	rows, err := q.db.Query(ctx, listCompletedReferralProjects, referrerUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Project{}
	for rows.Next() {
		var i Project
		if err := rows.Scan(
			&i.ProjectID,
			&i.Title,
			&i.PreferredPracticeArea,
			&i.PreferredPracticeLocation,
			&i.CaseDescription,
			&i.ReferrerUserID,
			&i.ReferredUserID,
			&i.Status,
			&i.CreatedAt,
			&i.StartedAt,
			&i.CompletedAt,
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

const listLawyers = `-- name: ListLawyers :many

SELECT
    u.user_id,
    u.first_name,
    u.last_name,
    u.avatar_url,
    COUNT(r.referral_user_id) AS referral_count
FROM
    users u
        LEFT JOIN
    referral_users r ON u.user_id = r.referred_user_id
GROUP BY
    u.user_id, u.first_name, u.last_name
ORDER BY
    u.user_id
`

type ListLawyersRow struct {
	UserID        string  `json:"user_id"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	AvatarUrl     *string `json:"avatar_url"`
	ReferralCount int64   `json:"referral_count"`
}

// lawyers
func (q *Queries) ListLawyers(ctx context.Context) ([]ListLawyersRow, error) {
	rows, err := q.db.Query(ctx, listLawyers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListLawyersRow{}
	for rows.Next() {
		var i ListLawyersRow
		if err := rows.Scan(
			&i.UserID,
			&i.FirstName,
			&i.LastName,
			&i.AvatarUrl,
			&i.ReferralCount,
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

const listLicenseUnVerifiedUsers = `-- name: ListLicenseUnVerifiedUsers :many
SELECT
    u.user_id,
    u.first_name,
    u.last_name,
    u.avatar_url,
    u.practice_location,
    u.join_date,
    l.license_id,
    l.license_number,
    l.name AS license_name,
    l.issue_date,
    l.issue_state
FROM
    users u
        LEFT JOIN
    licenses l ON u.user_id = l.user_id
WHERE
     u.license_verified = false
ORDER BY
    u.join_date DESC
LIMIT $1
OFFSET $2
`

type ListLicenseUnVerifiedUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListLicenseUnVerifiedUsersRow struct {
	UserID           string      `json:"user_id"`
	FirstName        string      `json:"first_name"`
	LastName         string      `json:"last_name"`
	AvatarUrl        *string     `json:"avatar_url"`
	PracticeLocation *string     `json:"practice_location"`
	JoinDate         time.Time   `json:"join_date"`
	LicenseID        *int64      `json:"license_id"`
	LicenseNumber    *string     `json:"license_number"`
	LicenseName      *string     `json:"license_name"`
	IssueDate        pgtype.Date `json:"issue_date"`
	IssueState       *string     `json:"issue_state"`
}

// AND u.license_rejected = false
func (q *Queries) ListLicenseUnVerifiedUsers(ctx context.Context, arg ListLicenseUnVerifiedUsersParams) ([]ListLicenseUnVerifiedUsersRow, error) {
	rows, err := q.db.Query(ctx, listLicenseUnVerifiedUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListLicenseUnVerifiedUsersRow{}
	for rows.Next() {
		var i ListLicenseUnVerifiedUsersRow
		if err := rows.Scan(
			&i.UserID,
			&i.FirstName,
			&i.LastName,
			&i.AvatarUrl,
			&i.PracticeLocation,
			&i.JoinDate,
			&i.LicenseID,
			&i.LicenseNumber,
			&i.LicenseName,
			&i.IssueDate,
			&i.IssueState,
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

const listLicenseVerifiedUsers = `-- name: ListLicenseVerifiedUsers :many
SELECT
    u.user_id,
    u.first_name,
    u.last_name,
    u.avatar_url,
    u.practice_location,
    u.join_date,
    l.license_id,
    l.license_number,
    l.name AS license_name,
    l.issue_date,
    l.issue_state
FROM
    users u
        LEFT JOIN
    licenses l ON u.user_id = l.user_id
WHERE
    u.license_verified = true
ORDER BY
    u.join_date DESC
LIMIT $1
OFFSET $2
`

type ListLicenseVerifiedUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListLicenseVerifiedUsersRow struct {
	UserID           string      `json:"user_id"`
	FirstName        string      `json:"first_name"`
	LastName         string      `json:"last_name"`
	AvatarUrl        *string     `json:"avatar_url"`
	PracticeLocation *string     `json:"practice_location"`
	JoinDate         time.Time   `json:"join_date"`
	LicenseID        *int64      `json:"license_id"`
	LicenseNumber    *string     `json:"license_number"`
	LicenseName      *string     `json:"license_name"`
	IssueDate        pgtype.Date `json:"issue_date"`
	IssueState       *string     `json:"issue_state"`
}

func (q *Queries) ListLicenseVerifiedUsers(ctx context.Context, arg ListLicenseVerifiedUsersParams) ([]ListLicenseVerifiedUsersRow, error) {
	rows, err := q.db.Query(ctx, listLicenseVerifiedUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListLicenseVerifiedUsersRow{}
	for rows.Next() {
		var i ListLicenseVerifiedUsersRow
		if err := rows.Scan(
			&i.UserID,
			&i.FirstName,
			&i.LastName,
			&i.AvatarUrl,
			&i.PracticeLocation,
			&i.JoinDate,
			&i.LicenseID,
			&i.LicenseNumber,
			&i.LicenseName,
			&i.IssueDate,
			&i.IssueState,
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
