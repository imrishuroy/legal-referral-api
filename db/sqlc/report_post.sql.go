// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: report_post.sql

package db

import (
	"context"
)

const addReportReason = `-- name: AddReportReason :one
INSERT INTO report_reasons (
    reason
) VALUES (
    $1
) RETURNING reason_id
`

func (q *Queries) AddReportReason(ctx context.Context, reason string) (int32, error) {
	row := q.db.QueryRow(ctx, addReportReason, reason)
	var reason_id int32
	err := row.Scan(&reason_id)
	return reason_id, err
}

const isPostReported = `-- name: IsPostReported :one
SELECT
    CASE WHEN report_id IS NOT NULL THEN true ELSE false END AS is_reported
FROM reported_posts
WHERE post_id = $1 AND reported_by = $2
`

type IsPostReportedParams struct {
	PostID     int32  `json:"post_id"`
	ReportedBy string `json:"reported_by"`
}

func (q *Queries) IsPostReported(ctx context.Context, arg IsPostReportedParams) (bool, error) {
	row := q.db.QueryRow(ctx, isPostReported, arg.PostID, arg.ReportedBy)
	var is_reported bool
	err := row.Scan(&is_reported)
	return is_reported, err
}

const reportPost = `-- name: ReportPost :exec
INSERT INTO reported_posts (
    post_id,
    reported_by,
    reason_id
) VALUES (
    $1, $2, $3
)
`

type ReportPostParams struct {
	PostID     int32  `json:"post_id"`
	ReportedBy string `json:"reported_by"`
	ReasonID   int32  `json:"reason_id"`
}

func (q *Queries) ReportPost(ctx context.Context, arg ReportPostParams) error {
	_, err := q.db.Exec(ctx, reportPost, arg.PostID, arg.ReportedBy, arg.ReasonID)
	return err
}
