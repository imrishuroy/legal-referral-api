package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
	"strconv"
)

type saveFeaturePostReq struct {
	PostID int32  `json:"post_id"`
	UserID string `json:"user_id"`
}

func (server *Server) saveFeaturePost(ctx *gin.Context) {
	var req saveFeaturePostReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != req.UserID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.SaveFeaturePostParams{
		PostID: req.PostID,
		UserID: req.UserID,
	}

	err := server.store.SaveFeaturePost(ctx, arg)
	if err != nil {

		errorCode := db.ErrorCode(err)
		if errorCode == db.UniqueViolation {
			ctx.JSON(400, gin.H{"message": "Post already saved"})
			return
		}
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, gin.H{"message": "success"})
}

func (server *Server) unSaveFeaturePost(ctx *gin.Context) {
	postIdStr := ctx.Param("post_id")

	postID, err := strconv.Atoi(postIdStr)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.UnSaveFeaturePostParams{
		PostID: int32(postID),
		UserID: authPayload.UID,
	}

	err = server.store.UnSaveFeaturePost(ctx, arg)

	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, gin.H{"message": "success"})
}

func (server *Server) listFeaturePosts(ctx *gin.Context) {

	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	posts, err := server.store.ListFeaturePosts(ctx)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, posts)
}
