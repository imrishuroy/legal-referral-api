package api

import (
	"firebase.google.com/go/v4/auth"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SignupMethod int32

const (
	Email SignupMethod = iota
	Google
	LinkedIn
	Facebook
)

func (s SignupMethod) String() string {
	return [...]string{"Email", "Google", "Microsoft", "LinkedIn"}[s]
}
func (s SignupMethod) Int32() int32 {
	return int32(s)
}

func (server *Server) createUser(ctx *gin.Context) {

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error parsing form"})
		return
	}

	email := form.Value["email"]
	firstName := form.Value["first_name"]
	lastName := form.Value["last_name"]
	mobile := form.Value["mobile"]
	signupMethod := form.Value["signup_method"]
	imageUrl := form.Value["image_url"]

	// Check for missing required fields
	requiredFields := []string{"email", "first_name", "mobile", "signup_method"}
	for _, field := range requiredFields {
		if len(form.Value[field]) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing required fields"})
			return
		}
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var userImageUrl string
	if len(imageUrl) != 0 {
		userImageUrl = imageUrl[0]
	}

	if len(form.File["user_image"]) > 0 {
		userImageFile := form.File["user_image"][0]
		file, err := userImageFile.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error opening file"})
			return
		}
		defer file.Close()
		// create file name with userid and file extension
		fileName := authPayload.UID + getFileExtension(userImageFile)

		url, err := server.uploadfile(file, fileName, userImageFile.Header.Get("Content-Type"), "user-images")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error uploading file"})
			return
		}
		userImageUrl = url
	}

	// convert signup method to int32
	userSignUpMethod, err := strconv.ParseInt(signupMethod[0], 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid signup method"})
		return
	}

	// create user
	arg := db.CreateUserParams{
		UserID:         authPayload.UID,
		FirstName:      firstName[0],
		LastName:       lastName[0],
		Email:          email[0],
		Mobile:         &mobile[0],
		AvatarUrl:      &userImageUrl,
		EmailVerified:  true,
		MobileVerified: int32(userSignUpMethod) == 0,
		SignupMethod:   int32(userSignUpMethod),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		errCode := db.ErrorCode(err)
		switch errCode {
		case db.ForeignKeyViolation, db.UniqueViolation:
			ctx.JSON(http.StatusForbidden, gin.H{"message": "User already exists"})
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type getUserByIdReq struct {
	UserID string `uri:"user_id" binding:"required"`
}

func (server *Server) getUserById(ctx *gin.Context) {
	var req getUserByIdReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != req.UserID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	user, _ := server.store.GetUserById(ctx, req.UserID)

	// if the user not found returning nil, not error
	if user.UserID == "" {
		ctx.JSON(http.StatusOK, nil)
		return
	}
	ctx.JSON(http.StatusOK, user)
	return
}

type updateUserRequest struct {
	ID              string `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Mobile          string `json:"mobile"`
	Address         string `json:"address"`
	EmailVerified   bool   `json:"email_verified"`
	MobileVerified  bool   `json:"mobile_verified"`
	WizardStep      int32  `json:"wizard_step"`
	WizardCompleted bool   `json:"wizard_completed"`
}

type getUserWizardStepReq struct {
	UserID string `uri:"user_id" binding:"required"`
}

func (server *Server) getUserWizardStep(ctx *gin.Context) {
	var req getUserWizardStepReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Err(err).Msg("error binding uri")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	log.Info().Msgf("user id %s", req.UserID)

	step, err := server.store.GetUserWizardStep(ctx, req.UserID)
	if err != nil {
		log.Error().Err(err).Msg("message getting user wizard step")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, step)
}

type markWizardCompletedReq struct {
	UserID string `uri:"user_id" binding:"required"`
}

func (server *Server) markWizardCompleted(ctx *gin.Context) {
	var req markWizardCompletedReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	markWizardCompArg := db.MarkWizardCompletedParams{
		UserID:          req.UserID,
		WizardCompleted: true,
	}

	_, err := server.store.MarkWizardCompleted(ctx, markWizardCompArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Wizard marked as completed"})
}

type saveAboutYouReq struct {
	Address          string `json:"address"`
	PracticeArea     string `json:"practice_area"`
	PracticeLocation string `json:"practice_location"`
	Experience       string `json:"experience"`
}

func (server *Server) saveAboutYou(ctx *gin.Context) {

	var req saveAboutYouReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.SaveAboutYouParams{
		UserID:           authPayload.UID,
		Address:          &req.Address,
		PracticeArea:     &req.PracticeArea,
		PracticeLocation: &req.PracticeLocation,
		Experience:       &req.Experience,
		WizardCompleted:  true,
	}

	_, err := server.store.SaveAboutYou(ctx, arg)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Error updating user about you")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "About you saved successfully"})
}

type updateUserInfo struct {
	FirstName               string `json:"first_name" binding:"required"`
	LastName                string `json:"last_name" binding:"required"`
	AverageBillingPerClient int32  `json:"average_billing_per_client" binding:"required"`
	CaseResolutionRate      int32  `json:"case_resolution_rate" binding:"required"`
	About                   string `json:"about" binding:"required"`
}

func (server *Server) updateUserInfo(ctx *gin.Context) {
	var req updateUserInfo
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	arg := db.UpdateUserInfoParams{
		UserID:                  authPayload.UID,
		FirstName:               req.FirstName,
		LastName:                req.LastName,
		AverageBillingPerClient: &req.AverageBillingPerClient,
		CaseResolutionRate:      &req.CaseResolutionRate,
		About:                   &req.About,
	}

	user, err := server.store.UpdateUserInfo(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
