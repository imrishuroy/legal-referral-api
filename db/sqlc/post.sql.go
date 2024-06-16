// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: post.sql

package db

import (
	"context"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts (
    owner_id,
    content,
    media,
    post_type,
    poll_id
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING post_id, owner_id, content, media, post_type, poll_id, created_at
`

type CreatePostParams struct {
	OwnerID  string   `json:"owner_id"`
	Content  *string  `json:"content"`
	Media    []string `json:"media"`
	PostType PostType `json:"post_type"`
	PollID   *int32   `json:"poll_id"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRow(ctx, createPost,
		arg.OwnerID,
		arg.Content,
		arg.Media,
		arg.PostType,
		arg.PollID,
	)
	var i Post
	err := row.Scan(
		&i.PostID,
		&i.OwnerID,
		&i.Content,
		&i.Media,
		&i.PostType,
		&i.PollID,
		&i.CreatedAt,
	)
	return i, err
}
