// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: notification.sql

package db

import (
	"context"
	"time"
)

const createNotification = `-- name: CreateNotification :one
INSERT INTO notifications (
    user_id,
    sender_id,
    target_id,
    target_type,
    notification_type,
    message
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING notification_id, user_id, sender_id, target_id, target_type, notification_type, message, is_read, created_at
`

type CreateNotificationParams struct {
	UserID           string `json:"user_id"`
	SenderID         string `json:"sender_id"`
	TargetID         int32  `json:"target_id"`
	TargetType       string `json:"target_type"`
	NotificationType string `json:"notification_type"`
	Message          string `json:"message"`
}

func (q *Queries) CreateNotification(ctx context.Context, arg CreateNotificationParams) (Notification, error) {
	row := q.db.QueryRow(ctx, createNotification,
		arg.UserID,
		arg.SenderID,
		arg.TargetID,
		arg.TargetType,
		arg.NotificationType,
		arg.Message,
	)
	var i Notification
	err := row.Scan(
		&i.NotificationID,
		&i.UserID,
		&i.SenderID,
		&i.TargetID,
		&i.TargetType,
		&i.NotificationType,
		&i.Message,
		&i.IsRead,
		&i.CreatedAt,
	)
	return i, err
}

const deleteNotificationById = `-- name: DeleteNotificationById :one
DELETE FROM notifications WHERE notification_id = $1 RETURNING notification_id, user_id, sender_id, target_id, target_type, notification_type, message, is_read, created_at
`

func (q *Queries) DeleteNotificationById(ctx context.Context, notificationID int32) (Notification, error) {
	row := q.db.QueryRow(ctx, deleteNotificationById, notificationID)
	var i Notification
	err := row.Scan(
		&i.NotificationID,
		&i.UserID,
		&i.SenderID,
		&i.TargetID,
		&i.TargetType,
		&i.NotificationType,
		&i.Message,
		&i.IsRead,
		&i.CreatedAt,
	)
	return i, err
}

const getNotificationById = `-- name: GetNotificationById :one
SELECT notification_id, user_id, sender_id, target_id, target_type, notification_type, message, is_read, created_at FROM notifications WHERE notification_id = $1
`

func (q *Queries) GetNotificationById(ctx context.Context, notificationID int32) (Notification, error) {
	row := q.db.QueryRow(ctx, getNotificationById, notificationID)
	var i Notification
	err := row.Scan(
		&i.NotificationID,
		&i.UserID,
		&i.SenderID,
		&i.TargetID,
		&i.TargetType,
		&i.NotificationType,
		&i.Message,
		&i.IsRead,
		&i.CreatedAt,
	)
	return i, err
}

const listNotifications = `-- name: ListNotifications :many
SELECT n.notification_id, n.user_id, n.sender_id, n.target_id, n.target_type, n.notification_type, n.message, n.is_read, n.created_at, u.first_name AS sender_first_name, u.last_name AS sender_last_name, u.avatar_url AS sender_avatar_url
FROM notifications n
JOIN users u ON n.sender_id = u.user_id
WHERE n.user_id = $1
ORDER BY n.created_at DESC
LIMIT $2 OFFSET $3
`

type ListNotificationsParams struct {
	UserID string `json:"user_id"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

type ListNotificationsRow struct {
	NotificationID   int32     `json:"notification_id"`
	UserID           string    `json:"user_id"`
	SenderID         string    `json:"sender_id"`
	TargetID         int32     `json:"target_id"`
	TargetType       string    `json:"target_type"`
	NotificationType string    `json:"notification_type"`
	Message          string    `json:"message"`
	IsRead           bool      `json:"is_read"`
	CreatedAt        time.Time `json:"created_at"`
	SenderFirstName  string    `json:"sender_first_name"`
	SenderLastName   string    `json:"sender_last_name"`
	SenderAvatarUrl  *string   `json:"sender_avatar_url"`
}

func (q *Queries) ListNotifications(ctx context.Context, arg ListNotificationsParams) ([]ListNotificationsRow, error) {
	rows, err := q.db.Query(ctx, listNotifications, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListNotificationsRow{}
	for rows.Next() {
		var i ListNotificationsRow
		if err := rows.Scan(
			&i.NotificationID,
			&i.UserID,
			&i.SenderID,
			&i.TargetID,
			&i.TargetType,
			&i.NotificationType,
			&i.Message,
			&i.IsRead,
			&i.CreatedAt,
			&i.SenderFirstName,
			&i.SenderLastName,
			&i.SenderAvatarUrl,
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

const markNotificationAsRead = `-- name: MarkNotificationAsRead :one
UPDATE notifications SET is_read = true WHERE notification_id = $1 RETURNING notification_id, user_id, sender_id, target_id, target_type, notification_type, message, is_read, created_at
`

func (q *Queries) MarkNotificationAsRead(ctx context.Context, notificationID int32) (Notification, error) {
	row := q.db.QueryRow(ctx, markNotificationAsRead, notificationID)
	var i Notification
	err := row.Scan(
		&i.NotificationID,
		&i.UserID,
		&i.SenderID,
		&i.TargetID,
		&i.TargetType,
		&i.NotificationType,
		&i.Message,
		&i.IsRead,
		&i.CreatedAt,
	)
	return i, err
}
