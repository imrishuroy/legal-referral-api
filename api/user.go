package api

import (
	"errors"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createUserReq struct {
	UserId         string       `json:"user_id"`
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

	// search if req email already exists in db
	dbUser, err := server.store.GetUserById(ctx, req.UserId)
	if err != nil {
		if !errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	// found the user with req email
	if dbUser.UserID != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "user with email already exists"})
		return
	}

	// create user
	arg := db.CreateUserParams{
		UserID:         req.UserId,
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
	user, err := server.store.GetUserById(ctx, req.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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

func (server *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	var mobileTxt pgtype.Text
	mobileTxt.String = req.Mobile

	var addressTxt pgtype.Text
	addressTxt.String = req.Address

	arg := db.UpdateUserParams{
		UserID:          req.ID,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Mobile:          &req.Mobile,
		Address:         &req.Address,
		EmailVerified:   req.EmailVerified,
		MobileVerified:  req.MobileVerified,
		WizardStep:      req.WizardStep,
		WizardCompleted: req.WizardCompleted,
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)
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
