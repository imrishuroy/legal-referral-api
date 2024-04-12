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
	CompanyName      string      `json:"company_name" binding:"required"`
	PracticeLocation string      `json:"practice_location" binding:"required"`
	StartDate        pgtype.Date `json:"start_date" binding:"required"`
	EndDate          pgtype.Date `json:"end_date" binding:"required"`
	Current          bool        `json:"current"`
	Description      string      `json:"description" binding:"required"`
	Skills           []string    `json:"skills" binding:"required"`
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
		CompanyName:      req.CompanyName,
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

	experience, err := server.store.AddExperience(ctx, arg)
	if err != nil {
		log.Error().Err(err).Msg("Error adding experience")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, experience)
}
