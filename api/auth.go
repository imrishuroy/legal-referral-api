package api

import (
	"encoding/json"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/url"
)

type resetPasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (server *Server) resetPassword(ctx *gin.Context) {
	var req resetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	if req.Email == "" || req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Email and Password are required"})
		return
	}

	user, err := server.firebaseAuth.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// delete the user
	err = server.firebaseAuth.DeleteUser(ctx, user.UID)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete user")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// create the user
	userArg := &auth.UserToCreate{}
	userArg.Email(req.Email)
	userArg.EmailVerified(false)
	userArg.Password(req.Password)
	userArg.UID(user.UID)

	_, err = server.firebaseAuth.CreateUser(ctx, userArg)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
	return
}

type linkedinLoginRequest struct {
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

type linkedinLoginResponse struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

func (server *Server) linkedinLogin(ctx *gin.Context) {

	var req linkedinLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	token, err := validateLinkedinToken(req.AccessToken, server.config.LinkedinClientID, server.config.LinkedinClientSecret)
	if err != nil {
		return
	}

	if !token.Active {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid access token"})
		return
	}

	userRecord, _ := server.firebaseAuth.GetUserByEmail(ctx, req.Email)
	if userRecord != nil {
		userID := userRecord.UserInfo.UID
		token, err := server.firebaseAuth.CustomToken(ctx, userID)
		if err != nil {
			log.Error().Err(err).Msg("failed to create custom token")
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, linkedinLoginResponse{UserID: userID, Token: token})
		return
	} else {
		userID := uuid.New().String()

		user := &auth.UserToCreate{}
		user.Email(req.Email)
		user.EmailVerified(true)
		user.UID(userID)

		createUser, err := server.firebaseAuth.CreateUser(ctx, user)
		if err != nil {
			return
		}

		token, err := server.firebaseAuth.CustomToken(ctx, createUser.UID)
		if err != nil {
			log.Error().Err(err).Msg("failed to create custom token")
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, linkedinLoginResponse{UserID: createUser.UID, Token: token})
	}
}

// TokenInfo represents the structure of the response from the introspection endpoint
type TokenInfo struct {
	Active       bool   `json:"active"`
	ClientID     string `json:"client_id"`
	AuthorizedAt int64  `json:"authorized_at"`
	CreatedAt    int64  `json:"created_at"`
	Status       string `json:"status"`
	ExpiresAt    int64  `json:"expires_at"`
	Scope        string `json:"scope"`
	AuthType     string `json:"auth_type"`
}

// IntrospectToken sends a POST request to the introspection endpoint to validate the token
func validateLinkedinToken(token string, clientID string, clientSecret string) (TokenInfo, error) {

	requestBody := url.Values{}
	requestBody.Set("token", token)
	requestBody.Set("client_id", clientID)
	requestBody.Set("client_secret", clientSecret)

	// Send POST request to the introspection endpoint
	response, err := http.PostForm("https://www.linkedin.com/oauth/v2/introspectToken", requestBody)
	if err != nil {
		return TokenInfo{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	// Decode the response body into a TokenInfo struct
	var tokenInfo TokenInfo
	err = json.NewDecoder(response.Body).Decode(&tokenInfo)
	if err != nil {
		return TokenInfo{}, err
	}

	return tokenInfo, nil
}
