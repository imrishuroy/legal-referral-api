package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
)

type listRecommendationsReq struct {
	Offset int32 `form:"offset" binding:"required"`
	Limit  int32 `form:"limit" binding:"required"`
}

func (s *Server) ListRecommendations(ctx *gin.Context) {
	var req listRecommendationsReq

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	arg := db.ListRecommendations2Params{
		UserID: userID,
		Offset: (req.Offset - 1) * req.Limit,
		Limit:  req.Limit,
	}

	recommendations, err := s.Store.ListRecommendations2(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, recommendations)
}

type cancelRecommendationReq struct {
	RecommendedUserID string `json:"recommended_user_id" binding:"required"`
}

func (s *Server) CancelRecommendation(ctx *gin.Context) {
	var req cancelRecommendationReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.CancelRecommendationParams{
		UserID:            authPayload.UID,
		RecommendedUserID: req.RecommendedUserID,
	}

	if err := s.Store.CancelRecommendation(ctx, arg); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Recommendation canceled"})
}
