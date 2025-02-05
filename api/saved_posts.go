package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
	"strconv"
)

type savePostReq struct {
	PostID int32  `json:"post_id" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
}

func (srv *Server) SavePost(ctx *gin.Context) {
	var req savePostReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != req.UserID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.SavePostParams{
		PostID: req.PostID,
		UserID: req.UserID,
	}

	err := srv.Store.SavePost(ctx, arg)
	if err != nil {

		errorCode := db.ErrorCode(err)
		if errorCode == db.UniqueViolation {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Post already saved"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Post saved successfully"})

}

func (srv *Server) UnSavePost(ctx *gin.Context) {

	savedPostParam := ctx.Param("saved_post_id")
	// convert this to int32
	savedPostID, err := strconv.Atoi(savedPostParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	err = srv.Store.UnsavePost(ctx, int32(savedPostID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Post unsaved successfully"})
}

func (srv *Server) ListSavedPosts(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	savedPosts, err := srv.Store.ListSavedPosts(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, savedPosts)
}
