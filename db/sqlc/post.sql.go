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

const getPosIsLikedByCurrentUser = `-- name: GetPosIsLikedByCurrentUser :one
SELECT
    CASE WHEN like_id IS NOT NULL THEN true ELSE false END AS is_liked
FROM likes
WHERE post_id = $1 AND user_id = $2 AND type = 'post'
`

type GetPosIsLikedByCurrentUserParams struct {
	PostID *int32 `json:"post_id"`
	UserID string `json:"user_id"`
}

func (q *Queries) GetPosIsLikedByCurrentUser(ctx context.Context, arg GetPosIsLikedByCurrentUserParams) (bool, error) {
	row := q.db.QueryRow(ctx, getPosIsLikedByCurrentUser, arg.PostID, arg.UserID)
	var is_liked bool
	err := row.Scan(&is_liked)
	return is_liked, err
}

const getPostCommentsCount = `-- name: GetPostCommentsCount :one
SELECT
    COUNT(*) AS comments_count
FROM comments
WHERE post_id = $1
`

func (q *Queries) GetPostCommentsCount(ctx context.Context, postID int32) (int64, error) {
	row := q.db.QueryRow(ctx, getPostCommentsCount, postID)
	var comments_count int64
	err := row.Scan(&comments_count)
	return comments_count, err
}

const getPostLikesAndCommentsCount = `-- name: GetPostLikesAndCommentsCount :one
SELECT
    COALESCE(likes_counts.likes_count, 0) AS likes_count,
    COALESCE(comments_counts.comments_count, 0) AS comments_count
FROM posts
            LEFT JOIN (
        SELECT
            post_id,
            COUNT(*) AS likes_count
        FROM likes
        WHERE type = 'post'
        GROUP BY post_id
    ) likes_counts ON posts.post_id = likes_counts.post_id
            LEFT JOIN (
        SELECT
            post_id,
            COUNT(*) AS comments_count
        FROM comments
        GROUP BY post_id
    ) comments_counts ON posts.post_id = comments_counts.post_id
    WHERE posts.post_id = $1
`

type GetPostLikesAndCommentsCountRow struct {
	LikesCount    int64 `json:"likes_count"`
	CommentsCount int64 `json:"comments_count"`
}

func (q *Queries) GetPostLikesAndCommentsCount(ctx context.Context, postID int32) (GetPostLikesAndCommentsCountRow, error) {
	row := q.db.QueryRow(ctx, getPostLikesAndCommentsCount, postID)
	var i GetPostLikesAndCommentsCountRow
	err := row.Scan(&i.LikesCount, &i.CommentsCount)
	return i, err
}

const getPostLikesCount = `-- name: GetPostLikesCount :one
SELECT
    COUNT(*) AS likes_count
FROM likes
WHERE post_id = $1 AND type = 'post'
`

func (q *Queries) GetPostLikesCount(ctx context.Context, postID *int32) (int64, error) {
	row := q.db.QueryRow(ctx, getPostLikesCount, postID)
	var likes_count int64
	err := row.Scan(&likes_count)
	return likes_count, err
}
