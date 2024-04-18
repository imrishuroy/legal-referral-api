package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"net/http"
)

type addExperienceReq struct {
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

type Experience struct {
	ExperienceID     int64       `json:"experience_id"`
	UserID           string      `json:"user_id"`
	Title            string      `json:"title"`
	PracticeArea     string      `json:"practice_area"`
	Firm             db.Firm     `json:"firm"`
	PracticeLocation string      `json:"practice_location"`
	StartDate        pgtype.Date `json:"start_date"`
	EndDate          pgtype.Date `json:"end_date"`
	Current          bool        `json:"current"`
	Description      string      `json:"description"`
	Skills           []string    `json:"skills"`
}

func (server *Server) addExperience(ctx *gin.Context) {
	var req addExperienceReq
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

	if req.Current && req.EndDate.Time.Before(req.StartDate.Time) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "End date should be greater than start date"})
		return
	}

	expRes, err := server.store.AddExperience(ctx, arg)
	if err != nil {
		log.Error().Err(err).Msg("Error adding experience")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// get the firm details
	firm, err := server.store.GetFirm(ctx, expRes.FirmID)
	if err != nil {
		log.Error().Err(err).Msg("Error getting firm details")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	experience := Experience{
		ExperienceID:     expRes.ExperienceID,
		UserID:           expRes.UserID,
		Title:            expRes.Title,
		PracticeArea:     expRes.PracticeArea,
		Firm:             firm,
		PracticeLocation: expRes.PracticeLocation,
		StartDate:        expRes.StartDate,
		EndDate:          expRes.EndDate,
		Current:          expRes.Current,
		Description:      expRes.Description,
		Skills:           expRes.Skills,
	}

	ctx.JSON(http.StatusOK, experience)
}
