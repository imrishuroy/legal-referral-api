package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
)

type listNewsFeedReq struct {
	Limit  int32 `form:"limit" binding:"required"`
	Offset int32 `form:"offset" binding:"required"`
}

func (server *Server) listNewsFeed(ctx *gin.Context) {

	var req listNewsFeedReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.ListNewsFeed2Params{
		UserID: userID,
		Limit:  req.Limit,
		Offset: (req.Offset - 1) * req.Limit,
	}

	feed, err := server.store.ListNewsFeed2(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, feed)

}
