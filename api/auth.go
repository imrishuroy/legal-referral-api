package api

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
)

// SignUpMethod represents different methods of signing up.
type SignUpMethod int

const (
	Email SignUpMethod = iota

	Google

	Apple
)

const IDLength = 32

type singUpRequest struct {
	Email           string       `json:"email"`
	FirstName       string       `json:"first_name"`
	LastName        string       `json:"last_name"`
	IsEmailVerified bool         `json:"is_email_verified"`
	SignUpMethod    SignUpMethod `json:"sign_up_method"`
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

	var userId string
	var email string

	if req.SignUpMethod == Email {
		// get the email from the token
		email = ExtractEmailFromIDToken(ctx)
		userId = generateUserID()
	} else {
		// check if this can give you email ( firebase )
		authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
		email = authPayload.Claims["email"].(string)
		userId = authPayload.UID
	}
	if email != req.Email {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email in the token and request body should be same"})
		return
	}

	// search if req email already exists in db
	dbUser, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if !errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
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
		ID:              userId,
		Email:           req.Email,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		SignUpMethod:    int32(req.SignUpMethod),
		IsEmailVerified: req.IsEmailVerified,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.ForeignKeyViolation || errCode == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)
}

type signInRequest struct {
	Email string `json:"email"`
}

func (server *Server) SignIn(ctx *gin.Context) {
	var req signInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	provider := ctx.GetHeader(provider)
	if provider == "" {
		err := errors.New("provider header is not provided")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse(err))
		ctx.Abort()
		return
	}
	signUpMethod, err := getSignUpMethod(provider)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid provider"})
		ctx.Abort()
		return
	}

	var email string

	switch signUpMethod {
	case Email:
		email = ExtractEmailFromIDToken(ctx)

	case Google, Apple:
		authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
		email = authPayload.Claims["email"].(string)

	default:
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid provider"})
	}

	if email != req.Email {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	// search if req email already exists in db
	dbUser, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if !errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
	}

	// found the user with req email
	if dbUser.ID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user with email does not exists"})
		return
	}

	ctx.JSON(http.StatusOK, dbUser)

}

func generateUserID() string {
	// Create a byte slice to store random bytes
	idBytes := make([]byte, IDLength)

	// Read random bytes from crypto/rand
	_, err := rand.Read(idBytes)
	if err != nil {
		panic(err) // Handle error
	}

	// Encode random bytes to base64 string
	id := base64.URLEncoding.EncodeToString(idBytes)

	// Trim any trailing "=" characters
	id = id[:IDLength]

	return id
}
