// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: profile.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const fetchUserProfile2 = `-- name: FetchUserProfile2 :one

SELECT
    users.user_id,
    users.first_name,
    users.last_name,
    users.practice_area,
    users.avatar_url,
    users.banner_url,
    users.average_billing_per_client,
    users.case_resolution_rate,
    users.open_to_referral,
    users.about,
    pricing.price_id,
    pricing.service_type,
    pricing.per_hour_price,
    pricing.per_hearing_price,
    pricing.contingency_price,
    pricing.hybrid_price

FROM users
LEFT JOIN pricing ON pricing.user_id = users.user_id
WHERE users.user_id = $1
`

type FetchUserProfile2Row struct {
	UserID                  string         `json:"user_id"`
	FirstName               string         `json:"first_name"`
	LastName                string         `json:"last_name"`
	PracticeArea            *string        `json:"practice_area"`
	AvatarUrl               *string        `json:"avatar_url"`
	BannerUrl               *string        `json:"banner_url"`
	AverageBillingPerClient *int32         `json:"average_billing_per_client"`
	CaseResolutionRate      *int32         `json:"case_resolution_rate"`
	OpenToReferral          bool           `json:"open_to_referral"`
	About                   *string        `json:"about"`
	PriceID                 *int64         `json:"price_id"`
	ServiceType             *string        `json:"service_type"`
	PerHourPrice            pgtype.Numeric `json:"per_hour_price"`
	PerHearingPrice         pgtype.Numeric `json:"per_hearing_price"`
	ContingencyPrice        *string        `json:"contingency_price"`
	HybridPrice             *string        `json:"hybrid_price"`
}

// -- name: FetchUserProfile :one
// SELECT sqlc.embed(users),
// COALESCE(sqlc.embed(pricing), '{}') as pricing
// FROM users
// LEFT JOIN pricing ON pricing.user_id = users.user_id
// WHERE users.user_id = $1;
func (q *Queries) FetchUserProfile2(ctx context.Context, userID string) (FetchUserProfile2Row, error) {
	row := q.db.QueryRow(ctx, fetchUserProfile2, userID)
	var i FetchUserProfile2Row
	err := row.Scan(
		&i.UserID,
		&i.FirstName,
		&i.LastName,
		&i.PracticeArea,
		&i.AvatarUrl,
		&i.BannerUrl,
		&i.AverageBillingPerClient,
		&i.CaseResolutionRate,
		&i.OpenToReferral,
		&i.About,
		&i.PriceID,
		&i.ServiceType,
		&i.PerHourPrice,
		&i.PerHearingPrice,
		&i.ContingencyPrice,
		&i.HybridPrice,
	)
	return i, err
}

const toggleOpenToRefferal = `-- name: ToggleOpenToRefferal :exec
UPDATE users
SET open_to_referral = $2
WHERE user_id = $1
`

type ToggleOpenToRefferalParams struct {
	UserID         string `json:"user_id"`
	OpenToReferral bool   `json:"open_to_referral"`
}

func (q *Queries) ToggleOpenToRefferal(ctx context.Context, arg ToggleOpenToRefferalParams) error {
	_, err := q.db.Exec(ctx, toggleOpenToRefferal, arg.UserID, arg.OpenToReferral)
	return err
}