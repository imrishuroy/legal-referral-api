// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: saved_posts.sql

package db

import (
	"context"
	"time"
)

const listSavedPosts = `-- name: ListSavedPosts :many
SELECT
    saved_posts.saved_post_id,
    posts.post_id, posts.owner_id, posts.content, posts.media, posts.post_type, posts.poll_id, posts.created_at,
    saved_posts.created_at
FROM
    saved_posts
JOIN
    posts ON saved_posts.post_id = posts.id
WHERE
    saved_posts.user_id = $1
ORDER BY
    saved_posts.created_at DESC
`

type ListSavedPostsRow struct {
	SavedPostID int32     `json:"saved_post_id"`
	Post        Post      `json:"post"`
	CreatedAt   time.Time `json:"created_at"`
}

func (q *Queries) ListSavedPosts(ctx context.Context, userID string) ([]ListSavedPostsRow, error) {
	rows, err := q.db.Query(ctx, listSavedPosts, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListSavedPostsRow{}
	for rows.Next() {
		var i ListSavedPostsRow
		if err := rows.Scan(
			&i.SavedPostID,
			&i.Post.PostID,
			&i.Post.OwnerID,
			&i.Post.Content,
			&i.Post.Media,
			&i.Post.PostType,
			&i.Post.PollID,
			&i.Post.CreatedAt,
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

const savePost = `-- name: SavePost :one
INSERT INTO saved_posts (
    user_id,
    post_id
) VALUES (
    $1, $2
) RETURNING saved_post_id, post_id, user_id, created_at
`

type SavePostParams struct {
	UserID string `json:"user_id"`
	PostID int32  `json:"post_id"`
}

func (q *Queries) SavePost(ctx context.Context, arg SavePostParams) (SavedPost, error) {
	row := q.db.QueryRow(ctx, savePost, arg.UserID, arg.PostID)
	var i SavedPost
	err := row.Scan(
		&i.SavedPostID,
		&i.PostID,
		&i.UserID,
		&i.CreatedAt,
	)
	return i, err
}
