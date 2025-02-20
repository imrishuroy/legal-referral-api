// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: activity.sql

package db

import (
	"context"
	"time"
)

const getUserFollowersCount = `-- name: GetUserFollowersCount :one
SELECT
COALESCE((SELECT COUNT(*)
          FROM connection_invitations
          WHERE recipient_id = users.user_id
            AND status NOT IN ('rejected', 'cancelled')), 0) AS followers_count
FROM users
WHERE users.user_id = $1
`

func (q *Queries) GetUserFollowersCount(ctx context.Context, userID string) (interface{}, error) {
	row := q.db.QueryRow(ctx, getUserFollowersCount, userID)
	var followers_count interface{}
	err := row.Scan(&followers_count)
	return followers_count, err
}

const listUserComments = `-- name: ListUserComments :many
SELECT
       comments.comment_id,
       comments.post_id,
       comments.content,
       comments.created_at,
       comments.parent_comment_id,
       users.user_id AS author_user_id,
       users.first_name AS author_first_name,
       users.last_name AS author_last_name,
       users.avatar_url AS author_avatar_url
FROM comments
         JOIN users ON comments.user_id = users.user_id
WHERE comments.user_id = $1
ORDER BY comments.created_at DESC
LIMIT $2
OFFSET $3
`

type ListUserCommentsParams struct {
	UserID string `json:"user_id"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

type ListUserCommentsRow struct {
	CommentID       int32     `json:"comment_id"`
	PostID          int32     `json:"post_id"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"created_at"`
	ParentCommentID *int32    `json:"parent_comment_id"`
	AuthorUserID    string    `json:"author_user_id"`
	AuthorFirstName string    `json:"author_first_name"`
	AuthorLastName  string    `json:"author_last_name"`
	AuthorAvatarUrl *string   `json:"author_avatar_url"`
}

func (q *Queries) ListUserComments(ctx context.Context, arg ListUserCommentsParams) ([]ListUserCommentsRow, error) {
	rows, err := q.db.Query(ctx, listUserComments, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListUserCommentsRow{}
	for rows.Next() {
		var i ListUserCommentsRow
		if err := rows.Scan(
			&i.CommentID,
			&i.PostID,
			&i.Content,
			&i.CreatedAt,
			&i.ParentCommentID,
			&i.AuthorUserID,
			&i.AuthorFirstName,
			&i.AuthorLastName,
			&i.AuthorAvatarUrl,
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

const listUserPosts = `-- name: ListUserPosts :many
SELECT posts.post_id, posts.owner_id, posts.content, posts.media, posts.post_type, posts.poll_id, posts.created_at,
       post_owner.first_name AS owner_first_name,
       post_owner.last_name AS owner_last_name,
       post_owner.avatar_url AS owner_avatar_url,
       post_owner.practice_area AS owner_practice_area,
       COALESCE(post_stats.likes, 0) AS likes_count,
       COALESCE(post_stats.comments, 0) AS comments_count
FROM posts
         JOIN users post_owner ON posts.owner_id = post_owner.user_id
         LEFT JOIN post_statistics post_stats ON posts.post_id = post_stats.post_id
WHERE posts.owner_id = $1
ORDER BY posts.created_at DESC
LIMIT $2
OFFSET $3
`

type ListUserPostsParams struct {
	OwnerID string `json:"owner_id"`
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
}

type ListUserPostsRow struct {
	PostID            int32     `json:"post_id"`
	OwnerID           string    `json:"owner_id"`
	Content           *string   `json:"content"`
	Media             []string  `json:"media"`
	PostType          PostType  `json:"post_type"`
	PollID            *int32    `json:"poll_id"`
	CreatedAt         time.Time `json:"created_at"`
	OwnerFirstName    string    `json:"owner_first_name"`
	OwnerLastName     string    `json:"owner_last_name"`
	OwnerAvatarUrl    *string   `json:"owner_avatar_url"`
	OwnerPracticeArea *string   `json:"owner_practice_area"`
	LikesCount        int64     `json:"likes_count"`
	CommentsCount     int64     `json:"comments_count"`
}

func (q *Queries) ListUserPosts(ctx context.Context, arg ListUserPostsParams) ([]ListUserPostsRow, error) {
	rows, err := q.db.Query(ctx, listUserPosts, arg.OwnerID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListUserPostsRow{}
	for rows.Next() {
		var i ListUserPostsRow
		if err := rows.Scan(
			&i.PostID,
			&i.OwnerID,
			&i.Content,
			&i.Media,
			&i.PostType,
			&i.PollID,
			&i.CreatedAt,
			&i.OwnerFirstName,
			&i.OwnerLastName,
			&i.OwnerAvatarUrl,
			&i.OwnerPracticeArea,
			&i.LikesCount,
			&i.CommentsCount,
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
