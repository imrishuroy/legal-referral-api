package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/rs/zerolog/log"
	"net/http"
)

type singUpRequest struct {
	Email           string `json:"email"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	IsEmailVerified bool   `json:"is_email_verified"`
}

func (server *Server) SignUp(ctx *gin.Context) {
	// TODO: parse the accessToken and verify the user, also
	// check the token email and the request email should be same
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

	// get the email from the token
	email := ExtractEmailFromIDToken(ctx)
	log.Info().Msgf("Email from token: %s", email)
	if email != req.Email {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email in the token and request body should be same"})
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

	if !req.IsEmailVerified {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email is not verified"})
		return
	}

	// create the user
	arg := db.CreateUserParams{
		Email:           req.Email,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		IsEmailVerified: req.IsEmailVerified,
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

//func (server *Server) SignIn(ctx gin.Context) {

//}
