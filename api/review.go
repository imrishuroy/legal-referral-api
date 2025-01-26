package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
)

type addReviewReq struct {
	ReviewerID string  `json:"reviewer_id" binding:"required"`
	Review     string  `json:"review" binding:"required"`
	Ratting    float64 `json:"ratting" binding:"required"`
}

func (s *Server) addReview(ctx *gin.Context) {
	var req addReviewReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	arg := db.AddReviewParams{
		UserID:     authPayload.UID,
		ReviewerID: req.ReviewerID,
		Review:     req.Review,
		Rating:     req.Ratting,
	}

	review, err := s.Store.AddReview(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, review)
}
