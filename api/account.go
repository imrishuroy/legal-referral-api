package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ratingInfo struct {
	AverageRating float64 `json:"average_rating"`
	Attorneys     int64   `json:"attorneys"`
}

type accountInfo struct {
	UserID           string      `json:"user_id"`
	FirstName        string      `json:"first_name"`
	LastName         string      `json:"last_name"`
	AvatarUrl        *string     `json:"avatar_url"`
	PracticeArea     *string     `json:"practice_area"`
	RatingInfo       *ratingInfo `json:"rating_info"`
	FollowersCount   int64       `json:"followers_count"`
	ConnectionsCount int64       `json:"connections_count"`
}

func (s *Server) getAccountInfo(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	acInfo, err := s.store.GetAccountInfo(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	avgRating, ok := acInfo.AverageRating.(float64)
	if !ok {
		avgRating = 0.0
	}

	attorneys, ok := acInfo.Attorneys.(int64)
	if !ok {
		attorneys = 0
	}

	followers, ok := acInfo.FollowersCount.(int64)
	if !ok {
		followers = 0
	}

	connections, ok := acInfo.ConnectionsCount.(int64)
	if !ok {
		connections = 0
	}

	account := accountInfo{
		UserID:       acInfo.UserID,
		FirstName:    acInfo.FirstName,
		LastName:     acInfo.LastName,
		AvatarUrl:    acInfo.AvatarUrl,
		PracticeArea: acInfo.PracticeArea,
		RatingInfo: &ratingInfo{
			AverageRating: avgRating,
			Attorneys:     attorneys,
		},
		FollowersCount:   followers,
		ConnectionsCount: connections,
	}

	ctx.JSON(http.StatusOK, account)

}
