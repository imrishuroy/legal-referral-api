package api

import (
	"errors"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
	"strconv"
)

type reportPostReq struct {
	PostID     int32  `json:"post_id" binding:"required"`
	ReportedBy string `json:"reported_by" binding:"required"`
	Reason     string `json:"reason" binding:"required"`
}

func (s *Server) reportPost(ctx *gin.Context) {

	var req reportPostReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != req.ReportedBy {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	reasonID, err := s.store.AddReportReason(ctx, req.Reason)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.ReportPostParams{
		PostID:     req.PostID,
		ReportedBy: req.ReportedBy,
		ReasonID:   reasonID,
	}

	err = s.store.ReportPost(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Post reported successfully"})
}

type getReportedPostsReq struct {
}

func (s *Server) isPostReported(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	postIDStr := ctx.Param("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid post ID"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.IsPostReportedParams{
		ReportedBy: userID,
		PostID:     int32(postID),
	}

	isReported, err := s.store.IsPostReported(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusOK, false)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, isReported)
}
