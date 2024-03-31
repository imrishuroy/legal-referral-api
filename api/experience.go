package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/rs/zerolog/log"
	"net/http"
)

type saveExperienceRequest struct {
	UserId           string `json:"userId"`
	PracticeArea     string `json:"practice_area"`
	PracticeLocation string `json:"practice_location"`
	Experience       int32  `json:"experience"`
}

func (server *Server) saveExperience(gin *gin.Context) {
	var req saveExperienceRequest
	if err := gin.ShouldBindJSON(&req); err != nil {
		gin.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.SaveExperienceParams{
		UserID:           req.UserId,
		PracticeArea:     req.PracticeArea,
		PracticeLocation: req.PracticeLocation,
		Experience:       req.Experience,
	}

	experience, err := server.store.SaveExperience(gin, arg)
	if err != nil {
		gin.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	gin.JSON(http.StatusOK, experience)

}

type saveAboutYouRequest struct {
	UserId           string `json:"user_id"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Address          string `json:"address"`
	PracticeArea     string `json:"practice_area"`
	PracticeLocation string `json:"practice_location"`
	Experience       int32  `json:"experience"`
}

func (server *Server) saveAboutYou(ctx *gin.Context) {

	var req saveAboutYouRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// updating the address in user table
	aboutArg := db.UpdateUserAboutYouParams{
		ID:        req.UserId,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Address:   req.Address,
	}

	_, err := server.store.UpdateUserAboutYou(ctx, aboutArg)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Error updating user about you")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}

	expArg := db.SaveExperienceParams{
		UserID:           req.UserId,
		PracticeArea:     req.PracticeArea,
		PracticeLocation: req.PracticeLocation,
		Experience:       req.Experience,
	}

	_, err = server.store.SaveExperience(ctx, expArg)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Error saving experience")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// update the wizard step
	wizardStepArg := db.UpdateUserWizardStepParams{
		ID:         req.UserId,
		WizardStep: 3,
	}

	_, err = server.store.UpdateUserWizardStep(ctx, wizardStepArg)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Error updating wizard step")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "About you saved successfully"})

}
