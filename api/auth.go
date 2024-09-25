package api

import (
	"bytes"
	"encoding/json"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/url"
)

type signInReq struct {
	Email             string `json:"email" binding:"required"`
	Password          string `json:"password" binding:"required"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

type signInRes struct {
	Kind         string `json:"kind"`
	LocalId      string `json:"localId"`
	Email        string `json:"email"`
	DisplayName  string `json:"displayName"`
	IdToken      string `json:"idToken"`
	Registered   bool   `json:"registered"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
}

type authResponse struct {
	User         db.User `json:"user"`
	IdToken      string  `json:"id_token"`
	RefreshToken string  `json:"refresh_token"`
	ExpiresIn    string  `json:"expires_in"`
}

func (server *Server) signIn(ctx *gin.Context) {
	var req signInReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	req.ReturnSecureToken = true

	if req.Email == "" || req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Email and Password are required"})
		return
	}

	authURL := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=" + server.config.FirebaseAuthKey

	// Marshal the request and make the API call
	resp, err := makePostRequest(authURL, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to make sign-in request")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer closeResponseBody(resp.Body)

	// Handle error cases based on status code
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("failed to read response body")
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to process response"})
			return
		}
		log.Info().Msg(string(body))
		handleFirebaseError(ctx, body)
		return
	}

	// Decode the response
	var res signInRes
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Error().Err(err).Msg("failed to decode sign-in response")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Retrieve user data from the database
	user, err := server.store.GetUserById(ctx, res.LocalId)
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve user")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Create and return the authentication response
	authResponse := authResponse{
		User:         user,
		IdToken:      res.IdToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    res.ExpiresIn,
	}

	ctx.JSON(http.StatusOK, authResponse)
}

type signUpReq struct {
	//FirstName         string  `json:"first_name" binding:"required"`
	//LastName          string  `json:"last_name" binding:"required"`
	//Email             string  `json:"email" binding:"required"`
	//Mobile            *string `json:"mobile" binding:"required"`
	//AvatarUrl         *string `json:"avatar_url"`
	//Password          string  `json:"password" binding:"required"`
	//ReturnSecureToken bool    `json:"returnSecureToken"`
	Email             string `json:"email" binding:"required"`
	Password          string `json:"password" binding:"required"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

type signUpRes struct {
	Kind         string `json:"kind"`
	IdToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalId      string `json:"localId"`
}

// TODO: convert this to multi-part request and expect avatar as file
func (server *Server) signUp(ctx *gin.Context) {

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error parsing form"})
		return
	}

	email := form.Value["email"]
	password := form.Value["password"]
	firstName := form.Value["first_name"]
	lastName := form.Value["last_name"]
	mobile := form.Value["mobile"]
	avatarUrl := form.Value["avatar_url"]
	avatarFile := form.File["avatar_file"]

	// Check for missing required fields
	requiredFields := []string{"email", "first_name", "mobile"}
	for _, field := range requiredFields {
		if len(form.Value[field]) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing required fields"})
			return
		}
	}

	signUpReq := signUpReq{
		Email:             email[0],
		Password:          password[0],
		ReturnSecureToken: true,
	}

	authURL := "https://identitytoolkit.googleapis.com/v1/accounts:signUp?key=" + server.config.FirebaseAuthKey

	// Marshal the request and make the API call
	resp, err := makePostRequest(authURL, signUpReq)
	if err != nil {
		log.Error().Err(err).Msg("failed to make sign-up request")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer closeResponseBody(resp.Body)

	// Handle error cases based on status code
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("failed to read response body")
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to process response"})
			return
		}
		log.Info().Msg(string(body))
		handleFirebaseError(ctx, body)
		return
	}

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password"})
		return
	}

	// Decode the response
	var res signUpRes
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Error().Err(err).Msg("failed to decode sign-up response")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var userImageUrl string
	if len(avatarUrl) != 0 {
		userImageUrl = avatarUrl[0]
	}

	if len(avatarFile) > 0 {
		userImageFile := avatarFile[0]
		file, err := userImageFile.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error opening file"})
			return
		}
		defer file.Close()
		// create file name with userid and file extension
		fileName := res.LocalId + getFileExtension(userImageFile)

		imageUrl, err := server.uploadFile(file, fileName, userImageFile.Header.Get("Content-Type"), "user-images")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error uploading file"})
			return
		}
		userImageUrl = imageUrl
	}

	arg := db.CreateUserParams{
		UserID:         res.LocalId,
		Email:          email[0],
		Mobile:         &mobile[0],
		FirstName:      firstName[0],
		LastName:       lastName[0],
		SignupMethod:   0,
		EmailVerified:  true,
		MobileVerified: true,
		AvatarUrl:      &userImageUrl,
	}

	user, err := server.store.CreateUser(ctx, arg)

	// Create and return the authentication response
	authResponse := authResponse{
		User:         user,
		IdToken:      res.IdToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    res.ExpiresIn,
	}

	ctx.JSON(http.StatusOK, authResponse)
}

type refreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	GrantType    string `json:"grant_type" binding:"required"`
}

type refreshTokenRes struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	UserId       string `json:"user_id"`
	ProjectId    string `json:"project_id"`
}

func (server *Server) refreshToken(ctx *gin.Context) {

	var req refreshTokenReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	if req.RefreshToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Refresh token is required"})
		return
	}

	authURL := "https://securetoken.googleapis.com/v1/token?key=" + server.config.FirebaseAuthKey

	// Marshal the request and make the API call
	resp, err := makePostRequest(authURL, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to make refresh token request")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer closeResponseBody(resp.Body)

	// Handle error cases based on status code
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("failed to read response body")
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to process response"})
			return
		}
		log.Info().Msg(string(body))
		handleFirebaseError(ctx, body)
		return
	}

	// Decode the response
	var res refreshTokenRes
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Error().Err(err).Msg("failed to decode refresh token response")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)

}

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

// Utility function to make POST request
func makePostRequest(url string, data interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return http.Post(url, "application/json", bytes.NewBuffer(jsonData))
}

// Utility function to safely close response body
func closeResponseBody(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		log.Error().Err(err).Msg("failed to close response body")
	}
}

func handleFirebaseError(ctx *gin.Context, body []byte) {
	var firebaseError struct {
		Error struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	// Attempt to unmarshal the error response
	if err := json.Unmarshal(body, &firebaseError); err != nil {
		log.Error().Err(err).Msg("failed to parse error response")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to process error response"})
		return
	}

	// Handle specific error messages
	switch firebaseError.Error.Message {
	case "EMAIL_EXISTS":
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
	case "INVALID_LOGIN_CREDENTIALS":
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login credentials"})
	default:
		// Generic error handling for any other errors
		ctx.JSON(http.StatusBadRequest, gin.H{"error": firebaseError.Error.Message})
	}
}
