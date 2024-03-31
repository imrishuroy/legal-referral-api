package api

import (
	"errors"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/imrishuroy/legal-referral/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateUserRequest struct {
	ID               string `json:"id"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Mobile           string `json:"mobile"`
	IsEmailVerified  bool   `json:"is_email_verified"`
	IsMobileVerified bool   `json:"is_mobile_verified"`
	WizardStep       int32  `json:"wizard_step"`
	WizardCompleted  bool   `json:"wizard_completed"`
}

func (server *Server) updateUser(ctx *gin.Context) {
	var req db.UpdateUserParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	arg := db.UpdateUserParams{
		ID:               req.ID,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Mobile:           req.Mobile,
		IsEmailVerified:  req.IsEmailVerified,
		IsMobileVerified: req.IsMobileVerified,
		WizardStep:       req.WizardStep,
		WizardCompleted:  req.WizardCompleted,
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)
}

type createUserReq struct {
	Email           string       `json:"email"`
	FirstName       string       `json:"first_name"`
	LastName        string       `json:"last_name"`
	IsEmailVerified bool         `json:"is_email_verified"`
	SignUpMethod    SignUpMethod `json:"sign_up_method"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Email == "" || req.FirstName == "" || req.LastName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email and Name are required"})
		return
	}

	// search if req email already exists in db
	dbUser, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if !errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	// found the user with req email
	if dbUser.ID != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user with email already exists"})
		return
	}

	// create user

	uuid, err := util.GenerateUUID()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		ID:              uuid,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Email:           req.Email,
		SignUpMethod:    int32(req.SignUpMethod),
		IsEmailVerified: false,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.ForeignKeyViolation || errCode == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}

	ctx.JSON(http.StatusOK, user)
}
