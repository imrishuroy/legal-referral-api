// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: discussion_message.sql

package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const listDiscussionMessages = `-- name: ListDiscussionMessages :many
SELECT
    m1.message_id, m1.parent_message_id, m1.discussion_id, m1.sender_id, m1.message, m1.sent_at,
    m2.message_id, m2.parent_message_id, m2.discussion_id, m2.sender_id, m2.message, m2.sent_at
FROM
    discussion_messages m1
        LEFT JOIN
    discussion_messages m2 ON m1.parent_message_id = m2.message_id
WHERE
    m1.discussion_id = $1
ORDER BY
    m1.sent_at DESC
OFFSET $2
LIMIT $3
`

type ListDiscussionMessagesParams struct {
	DiscussionID int32 `json:"discussion_id"`
	Offset       int32 `json:"offset"`
	Limit        int32 `json:"limit"`
}

type ListDiscussionMessagesRow struct {
	MessageID         int32              `json:"message_id"`
	ParentMessageID   *int32             `json:"parent_message_id"`
	DiscussionID      int32              `json:"discussion_id"`
	SenderID          string             `json:"sender_id"`
	Message           string             `json:"message"`
	SentAt            time.Time          `json:"sent_at"`
	MessageID_2       *int32             `json:"message_id_2"`
	ParentMessageID_2 *int32             `json:"parent_message_id_2"`
	DiscussionID_2    *int32             `json:"discussion_id_2"`
	SenderID_2        *string            `json:"sender_id_2"`
	Message_2         *string            `json:"message_2"`
	SentAt_2          pgtype.Timestamptz `json:"sent_at_2"`
}

func (q *Queries) ListDiscussionMessages(ctx context.Context, arg ListDiscussionMessagesParams) ([]ListDiscussionMessagesRow, error) {
	rows, err := q.db.Query(ctx, listDiscussionMessages, arg.DiscussionID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListDiscussionMessagesRow{}
	for rows.Next() {
		var i ListDiscussionMessagesRow
		if err := rows.Scan(
			&i.MessageID,
			&i.ParentMessageID,
			&i.DiscussionID,
			&i.SenderID,
			&i.Message,
			&i.SentAt,
			&i.MessageID_2,
			&i.ParentMessageID_2,
			&i.DiscussionID_2,
			&i.SenderID_2,
			&i.Message_2,
			&i.SentAt_2,
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

const listDiscussionMessages2 = `-- name: ListDiscussionMessages2 :many
SELECT
    m1.message_id,
    m1.parent_message_id,
    m1.discussion_id,
    m1.sender_id,
    m1.message,
    m1.sent_at,
    u1.avatar_url AS sender_avatar_image,
    u1.first_name AS sender_first_name,
    u1.last_name AS sender_last_name,
    m2.message_id AS reply_message_id,
    m2.parent_message_id AS reply_parent_message_id,
    m2.discussion_id AS reply_discussion_id,
    m2.sender_id AS reply_sender_id,
    m2.message AS reply_message,
    m2.sent_at AS reply_sent_at,
    u2.avatar_url AS reply_sender_avatar_image,
    u2.first_name AS reply_sender_first_name,
    u2.last_name AS reply_sender_last_name
FROM
    discussion_messages m1
        LEFT JOIN
    discussion_messages m2 ON m1.parent_message_id = m2.message_id
        LEFT JOIN
    users u1 ON m1.sender_id = u1.user_id
        LEFT JOIN
    users u2 ON m2.sender_id = u2.user_id
WHERE
    m1.discussion_id = $1
ORDER BY
    m1.sent_at DESC
OFFSET $2
LIMIT $3
`

type ListDiscussionMessages2Params struct {
	DiscussionID int32 `json:"discussion_id"`
	Offset       int32 `json:"offset"`
	Limit        int32 `json:"limit"`
}

type ListDiscussionMessages2Row struct {
	MessageID              int32              `json:"message_id"`
	ParentMessageID        *int32             `json:"parent_message_id"`
	DiscussionID           int32              `json:"discussion_id"`
	SenderID               string             `json:"sender_id"`
	Message                string             `json:"message"`
	SentAt                 time.Time          `json:"sent_at"`
	SenderAvatarImage      *string            `json:"sender_avatar_image"`
	SenderFirstName        *string            `json:"sender_first_name"`
	SenderLastName         *string            `json:"sender_last_name"`
	ReplyMessageID         *int32             `json:"reply_message_id"`
	ReplyParentMessageID   *int32             `json:"reply_parent_message_id"`
	ReplyDiscussionID      *int32             `json:"reply_discussion_id"`
	ReplySenderID          *string            `json:"reply_sender_id"`
	ReplyMessage           *string            `json:"reply_message"`
	ReplySentAt            pgtype.Timestamptz `json:"reply_sent_at"`
	ReplySenderAvatarImage *string            `json:"reply_sender_avatar_image"`
	ReplySenderFirstName   *string            `json:"reply_sender_first_name"`
	ReplySenderLastName    *string            `json:"reply_sender_last_name"`
}

func (q *Queries) ListDiscussionMessages2(ctx context.Context, arg ListDiscussionMessages2Params) ([]ListDiscussionMessages2Row, error) {
	rows, err := q.db.Query(ctx, listDiscussionMessages2, arg.DiscussionID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListDiscussionMessages2Row{}
	for rows.Next() {
		var i ListDiscussionMessages2Row
		if err := rows.Scan(
			&i.MessageID,
			&i.ParentMessageID,
			&i.DiscussionID,
			&i.SenderID,
			&i.Message,
			&i.SentAt,
			&i.SenderAvatarImage,
			&i.SenderFirstName,
			&i.SenderLastName,
			&i.ReplyMessageID,
			&i.ReplyParentMessageID,
			&i.ReplyDiscussionID,
			&i.ReplySenderID,
			&i.ReplyMessage,
			&i.ReplySentAt,
			&i.ReplySenderAvatarImage,
			&i.ReplySenderFirstName,
			&i.ReplySenderLastName,
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

const listDiscussionMessages3 = `-- name: ListDiscussionMessages3 :many
SELECT
    m1.message_id, m1.parent_message_id, m1.discussion_id, m1.sender_id, m1.message, m1.sent_at,
    u1.first_name AS sender_first_name,
    u1.last_name AS sender_last_name,
    u1.avatar_url AS sender_avatar_image
FROM
    discussion_messages m1
    LEFT JOIN
    users u1 ON m1.sender_id = u1.user_id
WHERE
    m1.discussion_id = $1
ORDER BY
    m1.sent_at ASC
OFFSET $2
LIMIT $3
`

type ListDiscussionMessages3Params struct {
	DiscussionID int32 `json:"discussion_id"`
	Offset       int32 `json:"offset"`
	Limit        int32 `json:"limit"`
}

type ListDiscussionMessages3Row struct {
	MessageID         int32     `json:"message_id"`
	ParentMessageID   *int32    `json:"parent_message_id"`
	DiscussionID      int32     `json:"discussion_id"`
	SenderID          string    `json:"sender_id"`
	Message           string    `json:"message"`
	SentAt            time.Time `json:"sent_at"`
	SenderFirstName   *string   `json:"sender_first_name"`
	SenderLastName    *string   `json:"sender_last_name"`
	SenderAvatarImage *string   `json:"sender_avatar_image"`
}

func (q *Queries) ListDiscussionMessages3(ctx context.Context, arg ListDiscussionMessages3Params) ([]ListDiscussionMessages3Row, error) {
	rows, err := q.db.Query(ctx, listDiscussionMessages3, arg.DiscussionID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListDiscussionMessages3Row{}
	for rows.Next() {
		var i ListDiscussionMessages3Row
		if err := rows.Scan(
			&i.MessageID,
			&i.ParentMessageID,
			&i.DiscussionID,
			&i.SenderID,
			&i.Message,
			&i.SentAt,
			&i.SenderFirstName,
			&i.SenderLastName,
			&i.SenderAvatarImage,
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

const sendMessageToDiscussion = `-- name: SendMessageToDiscussion :one
INSERT INTO discussion_messages (
    parent_message_id,
    sender_id,
    message,
    discussion_id
) VALUES (
    $1, $2, $3, $4
) RETURNING message_id, parent_message_id, discussion_id, sender_id, message, sent_at
`

type SendMessageToDiscussionParams struct {
	ParentMessageID *int32 `json:"parent_message_id"`
	SenderID        string `json:"sender_id"`
	Message         string `json:"message"`
	DiscussionID    int32  `json:"discussion_id"`
}

func (q *Queries) SendMessageToDiscussion(ctx context.Context, arg SendMessageToDiscussionParams) (DiscussionMessage, error) {
	row := q.db.QueryRow(ctx, sendMessageToDiscussion,
		arg.ParentMessageID,
		arg.SenderID,
		arg.Message,
		arg.DiscussionID,
	)
	var i DiscussionMessage
	err := row.Scan(
		&i.MessageID,
		&i.ParentMessageID,
		&i.DiscussionID,
		&i.SenderID,
		&i.Message,
		&i.SentAt,
	)
	return i, err
}
