// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"time"
)

type License struct {
	ID            int64  `json:"id"`
	UserID        string `json:"user_id"`
	Name          string `json:"name"`
	LicenseNumber string `json:"license_number"`
	IssueDate     string `json:"issue_date"`
	IssueState    string `json:"issue_state"`
}

type Otp struct {
	SessionID int64     `json:"session_id"`
	Email     string    `json:"email"`
	Channel   string    `json:"channel"`
	CreatedAt time.Time `json:"created_at"`
	Otp       int32     `json:"otp"`
}

type User struct {
	ID               string    `json:"id"`
	Email            string    `json:"email"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Mobile           string    `json:"mobile"`
	IsEmailVerified  bool      `json:"is_email_verified"`
	IsMobileVerified bool      `json:"is_mobile_verified"`
	WizardStep       int32     `json:"wizard_step"`
	WizardCompleted  bool      `json:"wizard_completed"`
	SignUpMethod     int32     `json:"sign_up_method"`
	JoinDate         time.Time `json:"join_date"`
}
