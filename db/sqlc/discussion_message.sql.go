// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: discussion_message.sql

package db

import (
	"context"
	"time"
)

const listDiscussionMessages = `-- name: ListDiscussionMessages :many
SELECT
    m1.message_id, m1.parent_message_id, m1.discussion_id, m1.sender_id, m1.message, m1.sent_at,
    u1.user_id AS sender_id,
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

type ListDiscussionMessagesParams struct {
	DiscussionID int32 `json:"discussion_id"`
	Offset       int32 `json:"offset"`
	Limit        int32 `json:"limit"`
}

type ListDiscussionMessagesRow struct {
	MessageID         int32     `json:"message_id"`
	ParentMessageID   *int32    `json:"parent_message_id"`
	DiscussionID      int32     `json:"discussion_id"`
	SenderID          string    `json:"sender_id"`
	Message           string    `json:"message"`
	SentAt            time.Time `json:"sent_at"`
	SenderID_2        *string   `json:"sender_id_2"`
	SenderFirstName   *string   `json:"sender_first_name"`
	SenderLastName    *string   `json:"sender_last_name"`
	SenderAvatarImage *string   `json:"sender_avatar_image"`
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
			&i.SenderID_2,
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
