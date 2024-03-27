package api

import (
	"errors"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
)

type SignUpMethod int

type singUpRequest struct {
	Email           string       `json:"email"`
	FirstName       string       `json:"first_name"`
	LastName        string       `json:"last_name"`
	IsEmailVerified bool         `json:"is_email_verified"`
	SignUpMethod    SignUpMethod `json:"sign_up_method"`
}

func (server *Server) signUp(ctx *gin.Context) {

	var req singUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// check all the required fields
	if req.Email == "" || req.FirstName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email and Name are required"})
		return
	}

	//var userId string
	//var email string

	//authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	//email = authPayload.Claims["email"].(string)
	//userId = authPayload.UID

	//if email != req.Email {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email in the token and request body should be same"})
	//	return
	//}

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

	// create the user
	//arg := db.CreateUserParams{
	//	//ID:              userId,
	//	ID:              GenerateRandomAlphaNumeric(16),
	//	Email:           req.Email,
	//	FirstName:       req.FirstName,
	//	LastName:        req.LastName,
	//	SignUpMethod:    int32(req.SignUpMethod),
	//	IsEmailVerified: req.IsEmailVerified,
	//}
	//
	//user, err := server.store.CreateUser(ctx, arg)
	//if err != nil {
	//	errCode := db.ErrorCode(err)
	//	if errCode == db.ForeignKeyViolation || errCode == db.UniqueViolation {
	//		ctx.JSON(http.StatusForbidden, errorResponse(err))
	//		return
	//	}
	//	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	//	return
	//}
	ctx.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

type signInRequest struct {
	Email string `json:"email"`
}

func (server *Server) signIn(ctx *gin.Context) {
	var req signInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	email := authPayload.Claims["email"].(string)

	if email != req.Email {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
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
	if dbUser.ID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user with email does not exists"})
		return
	}

	if !dbUser.IsEmailVerified {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email is not verified, " +
			"please verify your email to continue"})
		return
	}

	ctx.JSON(http.StatusOK, dbUser)
}
