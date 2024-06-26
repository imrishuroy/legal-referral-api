package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
	"strconv"
)

func (server *Server) createProjectReview(ctx *gin.Context) {
	var req *db.CreateProjectReviewParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != req.UserID {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.CreateProjectReviewParams{
		ProjectID: req.ProjectID,
		UserID:    req.UserID,
		Review:    req.Review,
		Rating:    req.Rating,
	}

	review, err := server.store.CreateProjectReview(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, review)
}

func (server *Server) getProjectReview(ctx *gin.Context) {
	projectIdParam := ctx.Param("project_id")
	if projectIdParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "project_id is required"})
		return
	}
	// convert project_id to int
	projectID, err := strconv.Atoi(projectIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "project_id must be a number"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.GetProjectReviewParams{
		ProjectID: int32(projectID),
		UserID:    authPayload.UID,
	}

	review, err := server.store.GetProjectReview(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, review)
}
