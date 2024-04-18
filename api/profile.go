package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

type fetchUserProfileReq struct {
	UserID string `uri:"user_id" binding:"required"`
}

func (server *Server) fetchUserProfile(ctx *gin.Context) {

	var req fetchUserProfileReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	log.Info().Msgf("user id %s", req.UserID)

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != req.UserID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	profile, err := server.store.FetchUserProfile(ctx, req.UserID)
	if err != nil {
		log.Err(err).Msg("error fetching user profile")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, profile)
}
