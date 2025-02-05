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

type addUpdateEducationReq struct {
	School       string      `json:"school" binding:"required"`
	Degree       string      `json:"degree" binding:"required"`
	FieldOfStudy string      `json:"field_of_study" binding:"required"`
	StartDate    pgtype.Date `json:"start_date" binding:"required"`
	EndDate      pgtype.Date `json:"end_date" binding:"required"`
	Current      bool        `json:"current"`
	Grade        string      `json:"grade" binding:"required"`
	Achievements string      `json:"achievements" binding:"required"`
	Skills       []string    `json:"skills" binding:"required"`
}

func (srv *Server) AddEducation(ctx *gin.Context) {
	var req addUpdateEducationReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.AddEducationParams{
		UserID:       authPayload.UID,
		School:       req.School,
		Degree:       req.Degree,
		FieldOfStudy: req.FieldOfStudy,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Current:      req.Current,
		Grade:        req.Grade,
		Achievements: req.Achievements,
		Skills:       req.Skills,
	}

	if !req.Current && (req.EndDate.Time.Before(req.StartDate.Time) || req.EndDate.Time.Equal(req.StartDate.Time)) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "End date should be greater than start date"})
		return
	}

	education, err := srv.Store.AddEducation(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, education)
}

func (srv *Server) ListEducations(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user id"})
		return
	}

	educations, err := srv.Store.ListEducations(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, educations)
}

func (srv *Server) UpdateEducation(ctx *gin.Context) {

	var req addUpdateEducationReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	educationIDParam := ctx.Param("education_id")
	educationID, err := strconv.ParseInt(educationIDParam, 10, 64)
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

	arg := db.UpdateEducationParams{
		EducationID:  educationID,
		School:       req.School,
		Degree:       req.Degree,
		FieldOfStudy: req.FieldOfStudy,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Current:      req.Current,
		Grade:        req.Grade,
		Achievements: req.Achievements,
		Skills:       req.Skills,
	}

	if !req.Current && (req.EndDate.Time.Before(req.StartDate.Time) || req.EndDate.Time.Equal(req.StartDate.Time)) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "End date should be greater than start date"})
		return
	}

	education, err := srv.Store.UpdateEducation(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, education)

}

func (srv *Server) DeleteEducation(ctx *gin.Context) {
	educationIDParam := ctx.Param("education_id")
	educationID, err := strconv.ParseInt(educationIDParam, 10, 64)
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

	err = srv.Store.DeleteEducation(ctx, educationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Education deleted successfully"})
}
