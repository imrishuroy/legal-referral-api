package api

import (
	"encoding/json"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

type likePostReq struct {
	PostUserID    string `json:"post_user_id" binding:"required"`
	CurrentUserID string `json:"current_user_id" binding:"required"`
}

func (srv *Server) LikePost(ctx *gin.Context) {
	var req likePostReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	postIDStr := ctx.Param("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid post ID"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != req.CurrentUserID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	postID32 := int32(postID)
	arg := db.LikePostParams{
		UserID: authPayload.UID,
		PostID: &postID32,
	}

	alreadyLiked := false
	err = srv.Store.LikePost(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) != db.UniqueViolation {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		} else {
			alreadyLiked = true
		}
	}

	if !alreadyLiked {

		// Prepare notification data
		data := map[string]string{
			"user_id":           req.PostUserID,
			"sender_id":         req.CurrentUserID,
			"target_id":         postIDStr,
			"target_type":       "post",
			"notification_type": "like",
			"already_liked":     strconv.FormatBool(alreadyLiked),
		}

		// Convert the map to a JSON string
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Error().Err(err).Msg("Error marshalling data")
		}

		// Launch a goroutine to publish to Kafka
		go func() {
			jsonString := string(jsonData)
			srv.publishToKafka("likes", authPayload.UID, jsonString)
		}()
	}
}

func (srv *Server) UnlikePost(ctx *gin.Context) {
	postIDStr := ctx.Param("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid post ID"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	postID32 := int32(postID)

	arg := db.UnlikePostParams{
		UserID: authPayload.UID,
		PostID: &postID32,
	}

	err = srv.Store.UnlikePost(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) != db.UniqueViolation {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	// Decrement likes
	err = srv.Store.DecrementLikes(ctx, postID32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

}

func (srv *Server) ListPostLikedUsers(ctx *gin.Context) {
	postIDStr := ctx.Param("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid post ID"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	postID32 := int32(postID)

	arg := db.ListPostLikedUsers2Params{
		PostID: &postID32,
		UserID: authPayload.UID,
	}

	users, err := srv.Store.ListPostLikedUsers2(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (srv *Server) LikeComment(ctx *gin.Context) {
	commentIDStr := ctx.Param("comment_id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid comment ID"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	commentID32 := int32(commentID)

	arg := db.LikeCommentParams{
		UserID:    authPayload.UID,
		CommentID: &commentID32,
	}

	err = srv.Store.LikeComment(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
}

func (srv *Server) UnlikeComment(ctx *gin.Context) {
	commentIDStr := ctx.Param("comment_id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid comment ID"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	commentID32 := int32(commentID)

	arg := db.UnlikeCommentParams{
		UserID:    authPayload.UID,
		CommentID: &commentID32,
	}

	err = srv.Store.UnlikeComment(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
}
