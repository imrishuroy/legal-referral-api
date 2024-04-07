package api

import (
	"errors"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/imrishuroy/legal-referral/util"
	"net/http"
)

type signUpCustomRequest struct {
	Email string `json:"email"`
}

func (server *Server) customTokenSignUp(ctx *gin.Context) {
	var req signUpCustomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	uid, err := util.GenerateUUID()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	claims := map[string]interface{}{
		"role":  "user",
		"email": req.Email,
	}
	token, err := server.firebaseAuth.CustomTokenWithClaims(ctx, uid, claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, token)
}

type SignupMethod int

type singUpRequest struct {
	Email         string       `json:"email"`
	FirstName     string       `json:"first_name"`
	LastName      string       `json:"last_name"`
	EmailVerified bool         `json:"email_verified"`
	SignupMethod  SignupMethod `json:"signup_method"`
}

func (server *Server) signUp(ctx *gin.Context) {

	var req singUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	// check all the required fields
	if req.Email == "" || req.FirstName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Email and Name are required"})
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
	if dbUser.UserID != "" {
		// TODO: don't throw error if user not found, instead return the user found, status code should be 200
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "user with email already exists"})
		return
	}

	userId, err := util.GenerateUUID()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// create the user
	arg := db.CreateUserParams{
		UserID:        userId,
		Email:         req.Email,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		SignupMethod:  int32(req.SignupMethod),
		EmailVerified: req.EmailVerified,
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

	err = sendOTP(server, user.Email, "email", server.config.VerifyEmailServiceSID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type signInRequest struct {
	Email string `json:"email"`
}

func (server *Server) signIn(ctx *gin.Context) {
	var req signInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
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
	if dbUser.UserID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user with email does not exists"})
		return
	}

	if !dbUser.EmailVerified {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email is not verified, " +
			"please verify your email to continue"})
		return
	}

	ctx.JSON(http.StatusOK, dbUser)
}
