// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: connection.sql

package db

import (
	"context"
	"time"
)

const acceptConnection = `-- name: AcceptConnection :one
UPDATE connection_invitations
SET status = 1
WHERE id = $1 AND status = 0
RETURNING id, sender_id, recipient_id, status, created_at
`

func (q *Queries) AcceptConnection(ctx context.Context, id int32) (ConnectionInvitation, error) {
	row := q.db.QueryRow(ctx, acceptConnection, id)
	var i ConnectionInvitation
	err := row.Scan(
		&i.ID,
		&i.SenderID,
		&i.RecipientID,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const addConnection = `-- name: AddConnection :exec
INSERT INTO connections (sender_id, recipient_id)
    VALUES ($1, $2)
`

type AddConnectionParams struct {
	SenderID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id"`
}

func (q *Queries) AddConnection(ctx context.Context, arg AddConnectionParams) error {
	_, err := q.db.Exec(ctx, addConnection, arg.SenderID, arg.RecipientID)
	return err
}

const listConnectionInvitations = `-- name: ListConnectionInvitations :many
SELECT ci.id, ci.sender_id, ci.recipient_id, ci.status, ci.created_at,
       u.first_name,
       u.last_name,
       u.about,
       u.avatar_url
FROM connection_invitations ci
JOIN users u ON ci.recipient_id = u.user_id
WHERE ci.recipient_id = $1 AND ci.status = 0
ORDER BY ci.created_at DESC
OFFSET $2
LIMIT $3
`

type ListConnectionInvitationsParams struct {
	RecipientID string `json:"recipient_id"`
	Offset      int32  `json:"offset"`
	Limit       int32  `json:"limit"`
}

type ListConnectionInvitationsRow struct {
	ID          int32     `json:"id"`
	SenderID    string    `json:"sender_id"`
	RecipientID string    `json:"recipient_id"`
	Status      int32     `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	About       *string   `json:"about"`
	AvatarUrl   *string   `json:"avatar_url"`
}

func (q *Queries) ListConnectionInvitations(ctx context.Context, arg ListConnectionInvitationsParams) ([]ListConnectionInvitationsRow, error) {
	rows, err := q.db.Query(ctx, listConnectionInvitations, arg.RecipientID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListConnectionInvitationsRow{}
	for rows.Next() {
		var i ListConnectionInvitationsRow
		if err := rows.Scan(
			&i.ID,
			&i.SenderID,
			&i.RecipientID,
			&i.Status,
			&i.CreatedAt,
			&i.FirstName,
			&i.LastName,
			&i.About,
			&i.AvatarUrl,
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

const listConnections = `-- name: ListConnections :many
SELECT ci.id, ci.sender_id, ci.recipient_id, ci.created_at,
       u.first_name,
       u.last_name,
       u.about,
       u.avatar_url
FROM connections ci
JOIN users u ON ci.recipient_id = u.user_id
WHERE sender_id = $3::text OR recipient_id = $3
ORDER BY created_at DESC
OFFSET $1
LIMIT $2
`

type ListConnectionsParams struct {
	Offset int32  `json:"offset"`
	Limit  int32  `json:"limit"`
	UserID string `json:"user_id"`
}

type ListConnectionsRow struct {
	ID          int32     `json:"id"`
	SenderID    string    `json:"sender_id"`
	RecipientID string    `json:"recipient_id"`
	CreatedAt   time.Time `json:"created_at"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	About       *string   `json:"about"`
	AvatarUrl   *string   `json:"avatar_url"`
}

func (q *Queries) ListConnections(ctx context.Context, arg ListConnectionsParams) ([]ListConnectionsRow, error) {
	rows, err := q.db.Query(ctx, listConnections, arg.Offset, arg.Limit, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListConnectionsRow{}
	for rows.Next() {
		var i ListConnectionsRow
		if err := rows.Scan(
			&i.ID,
			&i.SenderID,
			&i.RecipientID,
			&i.CreatedAt,
			&i.FirstName,
			&i.LastName,
			&i.About,
			&i.AvatarUrl,
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

const rejectConnection = `-- name: RejectConnection :exec
UPDATE connection_invitations
    SET status = 3
    WHERE sender_id = $1 AND recipient_id = $2 AND status = 'pending'
`

type RejectConnectionParams struct {
	SenderID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id"`
}

func (q *Queries) RejectConnection(ctx context.Context, arg RejectConnectionParams) error {
	_, err := q.db.Exec(ctx, rejectConnection, arg.SenderID, arg.RecipientID)
	return err
}

const sendConnection = `-- name: SendConnection :one
INSERT INTO connection_invitations (
    sender_id,
    recipient_id
) VALUES ($1, $2)
RETURNING (id)
`

type SendConnectionParams struct {
	SenderID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id"`
}

func (q *Queries) SendConnection(ctx context.Context, arg SendConnectionParams) (int32, error) {
	row := q.db.QueryRow(ctx, sendConnection, arg.SenderID, arg.RecipientID)
	var id int32
	err := row.Scan(&id)
	return id, err
}
