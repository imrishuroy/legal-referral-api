// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    id,
    first_name,
    last_name,
    mobile_number,
    email,
    bar_licence_no,
    practicing_field,
    experience
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING id, first_name, last_name, mobile_number, email, bar_licence_no, practicing_field, experience, join_date
`

type CreateUserParams struct {
	ID              string `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	MobileNumber    string `json:"mobile_number"`
	Email           string `json:"email"`
	BarLicenceNo    string `json:"bar_licence_no"`
	PracticingField string `json:"practicing_field"`
	Experience      int32  `json:"experience"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.ID,
		arg.FirstName,
		arg.LastName,
		arg.MobileNumber,
		arg.Email,
		arg.BarLicenceNo,
		arg.PracticingField,
		arg.Experience,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.MobileNumber,
		&i.Email,
		&i.BarLicenceNo,
		&i.PracticingField,
		&i.Experience,
		&i.JoinDate,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, first_name, last_name, mobile_number, email, bar_licence_no, practicing_field, experience, join_date FROM users
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.MobileNumber,
		&i.Email,
		&i.BarLicenceNo,
		&i.PracticingField,
		&i.Experience,
		&i.JoinDate,
	)
	return i, err
}
