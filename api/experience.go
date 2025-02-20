package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

type addUpdateExperienceReq struct {
	Title            string      `json:"title" binding:"required"`
	PracticeArea     string      `json:"practice_area" binding:"required"`
	FirmID           int64       `json:"firm_id" binding:"required"`
	PracticeLocation string      `json:"practice_location" binding:"required"`
	StartDate        pgtype.Date `json:"start_date" binding:"required"`
	EndDate          pgtype.Date `json:"end_date" binding:"required"`
	Current          bool        `json:"current"`
	Description      string      `json:"description" binding:"required"`
	Skills           []string    `json:"skills" binding:"required"`
}

type UserExperience struct {
	Experience db.Experience `json:"experience"`
	Firm       db.Firm       `json:"firm"`
}

func (srv *Server) AddExperience(ctx *gin.Context) {
	var req addUpdateExperienceReq
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

	arg := db.AddExperienceParams{
		UserID:           authPayload.UID,
		Title:            req.Title,
		PracticeArea:     req.PracticeArea,
		FirmID:           req.FirmID,
		PracticeLocation: req.PracticeLocation,
		StartDate:        req.StartDate,
		EndDate:          req.EndDate,
		Current:          req.Current,
		Description:      req.Description,
		Skills:           req.Skills,
	}

	// check if end time is greater than start time when end time is provided

	if !req.Current && (req.EndDate.Time.Before(req.StartDate.Time) || req.EndDate.Time.Equal(req.StartDate.Time)) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "End date should be greater than start date"})
		return
	}

	expRes, err := srv.Store.AddExperience(ctx, arg)
	if err != nil {
		log.Error().Err(err).Msg("Error adding experience")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// get the firm details
	firm, err := srv.Store.GetFirm(ctx, expRes.FirmID)
	if err != nil {
		log.Error().Err(err).Msg("Error getting firm details")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	experience := UserExperience{
		Experience: expRes,
		Firm:       firm,
	}

	ctx.JSON(http.StatusOK, experience)
}

type listExperienceResponse struct {
	ExperienceId int64        `json:"experience_id"`
	Title        string       `json:"title"`
	PracticeArea string       `json:"practice_area"`
	Description  string       `json:"description"`
	StartDate    pgtype.Date  `json:"start_date"`
	EndDate      *pgtype.Date `json:"end_date"`
	Current      bool         `json:"current"`
	Skills       []string     `json:"skills"`
	Firm         db.Firm      `json:"firm"`
}

func (srv *Server) ListExperiences(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}

	experiences, err := srv.Store.ListExperiences(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("Error listing experiences")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, experiences)
}

func (srv *Server) UpdateExperience(ctx *gin.Context) {

	var req addUpdateExperienceReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	experienceIDParam := ctx.Param("experience_id")
	experienceID, err := strconv.ParseInt(experienceIDParam, 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Invalid entity id")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid entity id"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		log.Error().Msg("Unauthorized")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.UpdateExperienceParams{
		ExperienceID:     experienceID,
		Title:            req.Title,
		PracticeArea:     req.PracticeArea,
		FirmID:           req.FirmID,
		PracticeLocation: req.PracticeLocation,
		StartDate:        req.StartDate,
		EndDate:          req.EndDate,
		Current:          req.Current,
		Description:      req.Description,
		Skills:           req.Skills,
	}

	// log end date
	log.Info().Msgf("End date: %v", req.EndDate.Time)

	// check if end time is greater than start time when end time is provided
	if !req.Current && (req.EndDate.Time.Before(req.StartDate.Time) || req.EndDate.Time.Equal(req.StartDate.Time)) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "End date should be greater than start date"})
		return
	}

	expRes, err := srv.Store.UpdateExperience(ctx, arg)

	if err != nil {
		log.Error().Err(err).Msg("Error updating experience")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// get the firm details
	firm, err := srv.Store.GetFirm(ctx, expRes.FirmID)
	if err != nil {
		log.Error().Err(err).Msg("Error getting firm details")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	experience := UserExperience{
		Experience: expRes,
		Firm:       firm,
	}

	ctx.JSON(http.StatusOK, experience)

}

func (srv *Server) DeleteExperience(ctx *gin.Context) {
	experienceIDParam := ctx.Param("experience_id")

	experienceID, err := strconv.ParseInt(experienceIDParam, 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Invalid entity id")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid entity id"})
		return
	}

	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		log.Error().Msg("Unauthorized")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}

	err = srv.Store.DeleteExperience(ctx, experienceID)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting experience")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Experience deleted successfully"})
}
