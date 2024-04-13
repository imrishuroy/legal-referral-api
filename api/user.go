package api

import (
	"errors"
	"firebase.google.com/go/v4/auth"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/rs/zerolog/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignupMethod int32

const (
	Email SignupMethod = iota
	Google
	Microsoft
	LinkedIn
)

func (s SignupMethod) String() string {
	return [...]string{"Email", "Google", "Microsoft", "LinkedIn"}[s]
}
func (s SignupMethod) Int32() int32 {
	return int32(s)
}

type createUserReq struct {
	Email          string       `json:"email"`
	FirstName      string       `json:"first_name"`
	LastName       string       `json:"last_name"`
	Mobile         string       `json:"mobile"`
	ImageUrl       string       `json:"image_url"`
	EmailVerified  bool         `json:"email_verified"`
	MobileVerified bool         `json:"mobile_verified"`
	SignupMethod   SignupMethod `json:"signup_method"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	if req.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Email is required"})
		return
	}

	if req.FirstName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "First Name is required"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// search if req email already exists in db
	dbUser, err := server.store.GetUserById(ctx, authPayload.UID)
	if err != nil {
		if !errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	// found the user with req email
	// Check if dbUser.UserID is not empty, indicating that a user with that email already exists
	if dbUser.UserID != "" {
		// If a user with the provided email already exists, return a bad request response
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "User with email already exists"})
		return
	}

	// create user
	arg := db.CreateUserParams{
		UserID:         authPayload.UID,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		Mobile:         &req.Mobile,
		ImageUrl:       &req.ImageUrl,
		EmailVerified:  req.EmailVerified,
		MobileVerified: req.MobileVerified,
		SignupMethod:   int32(req.SignupMethod),
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

type updateUserImageReq struct {
	UserID   string `uri:"user_id" binding:"required"`
	ImageUrl string `json:"image_url"`
}

func (server *Server) updateUserImage(ctx *gin.Context) {
	var req updateUserImageReq

	// Bind URI parameters
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Err(err).Msg("error binding uri")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid URI parameters"})
		return
	}

	// Bind JSON body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Err(err).Msg("error binding json")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON body"})
		return
	}

	log.Info().Msgf("user id %s", req.UserID)
	log.Info().Msgf("image url %s", req.ImageUrl)

	var profileImageArg = db.UpdateUserImageUrlParams{
		UserID:   req.UserID,
		ImageUrl: &req.ImageUrl,
	}

	_, err := server.store.UpdateUserImageUrl(ctx, profileImageArg)
	if err != nil {
		log.Error().Err(err).Msg("error updating user profile image")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var wizardStepArg = db.UpdateUserWizardStepParams{
		UserID:     req.UserID,
		WizardStep: 1,
	}

	_, err = server.store.UpdateUserWizardStep(ctx, wizardStepArg)
	if err != nil {
		log.Error().Err(err).Msg("message updating user wizard step")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Profile image updated successfully"})
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
