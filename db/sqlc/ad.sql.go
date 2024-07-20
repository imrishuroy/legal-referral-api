// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: ad.sql

package db

import (
	"context"
	"time"
)

const createAd = `-- name: CreateAd :one
INSERT INTO ads (
    ad_type,
    title,
    description,
    link,
    media,
    payment_cycle,
    author_id,
    start_date,
    end_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING ad_id, ad_type, title, description, link, media, payment_cycle, author_id, start_date, end_date, created_at
`

type CreateAdParams struct {
	AdType       AdType       `json:"ad_type"`
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	Link         string       `json:"link"`
	Media        []string     `json:"media"`
	PaymentCycle PaymentCycle `json:"payment_cycle"`
	AuthorID     string       `json:"author_id"`
	StartDate    time.Time    `json:"start_date"`
	EndDate      time.Time    `json:"end_date"`
}

func (q *Queries) CreateAd(ctx context.Context, arg CreateAdParams) (Ad, error) {
	row := q.db.QueryRow(ctx, createAd,
		arg.AdType,
		arg.Title,
		arg.Description,
		arg.Link,
		arg.Media,
		arg.PaymentCycle,
		arg.AuthorID,
		arg.StartDate,
		arg.EndDate,
	)
	var i Ad
	err := row.Scan(
		&i.AdID,
		&i.AdType,
		&i.Title,
		&i.Description,
		&i.Link,
		&i.Media,
		&i.PaymentCycle,
		&i.AuthorID,
		&i.StartDate,
		&i.EndDate,
		&i.CreatedAt,
	)
	return i, err
}

const extendAdPeriod = `-- name: ExtendAdPeriod :one
UPDATE ads
SET end_date = $2
WHERE ad_id = $1
RETURNING ad_id, ad_type, title, description, link, media, payment_cycle, author_id, start_date, end_date, created_at
`

type ExtendAdPeriodParams struct {
	AdID    int32     `json:"ad_id"`
	EndDate time.Time `json:"end_date"`
}

func (q *Queries) ExtendAdPeriod(ctx context.Context, arg ExtendAdPeriodParams) (Ad, error) {
	row := q.db.QueryRow(ctx, extendAdPeriod, arg.AdID, arg.EndDate)
	var i Ad
	err := row.Scan(
		&i.AdID,
		&i.AdType,
		&i.Title,
		&i.Description,
		&i.Link,
		&i.Media,
		&i.PaymentCycle,
		&i.AuthorID,
		&i.StartDate,
		&i.EndDate,
		&i.CreatedAt,
	)
	return i, err
}

const getRandomAd = `-- name: GetRandomAd :one
SELECT ad_id, ad_type, title, description, link, media, payment_cycle, author_id, start_date, end_date, created_at
FROM ads
WHERE end_date > NOW()
ORDER BY RANDOM()
LIMIT 1
`

func (q *Queries) GetRandomAd(ctx context.Context) (Ad, error) {
	row := q.db.QueryRow(ctx, getRandomAd)
	var i Ad
	err := row.Scan(
		&i.AdID,
		&i.AdType,
		&i.Title,
		&i.Description,
		&i.Link,
		&i.Media,
		&i.PaymentCycle,
		&i.AuthorID,
		&i.StartDate,
		&i.EndDate,
		&i.CreatedAt,
	)
	return i, err
}

const listExpiredAds = `-- name: ListExpiredAds :many
SELECT ad_id, ad_type, title, description, link, media, payment_cycle, author_id, start_date, end_date, created_at
FROM ads
WHERE end_date < NOW()
ORDER BY start_date DESC
`

func (q *Queries) ListExpiredAds(ctx context.Context) ([]Ad, error) {
	rows, err := q.db.Query(ctx, listExpiredAds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Ad{}
	for rows.Next() {
		var i Ad
		if err := rows.Scan(
			&i.AdID,
			&i.AdType,
			&i.Title,
			&i.Description,
			&i.Link,
			&i.Media,
			&i.PaymentCycle,
			&i.AuthorID,
			&i.StartDate,
			&i.EndDate,
			&i.CreatedAt,
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

const listPlayingAds = `-- name: ListPlayingAds :many
SELECT ad_id, ad_type, title, description, link, media, payment_cycle, author_id, start_date, end_date, created_at
FROM ads
WHERE end_date > NOW()
ORDER BY start_date DESC
`

func (q *Queries) ListPlayingAds(ctx context.Context) ([]Ad, error) {
	rows, err := q.db.Query(ctx, listPlayingAds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Ad{}
	for rows.Next() {
		var i Ad
		if err := rows.Scan(
			&i.AdID,
			&i.AdType,
			&i.Title,
			&i.Description,
			&i.Link,
			&i.Media,
			&i.PaymentCycle,
			&i.AuthorID,
			&i.StartDate,
			&i.EndDate,
			&i.CreatedAt,
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

const listRandomAds = `-- name: ListRandomAds :many
SELECT ad_id, ad_type, title, description, link, media, payment_cycle, author_id, start_date, end_date, created_at
FROM ads
WHERE end_date > NOW()
ORDER BY RANDOM()
LIMIT $1
`

func (q *Queries) ListRandomAds(ctx context.Context, limit int32) ([]Ad, error) {
	rows, err := q.db.Query(ctx, listRandomAds, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Ad{}
	for rows.Next() {
		var i Ad
		if err := rows.Scan(
			&i.AdID,
			&i.AdType,
			&i.Title,
			&i.Description,
			&i.Link,
			&i.Media,
			&i.PaymentCycle,
			&i.AuthorID,
			&i.StartDate,
			&i.EndDate,
			&i.CreatedAt,
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
