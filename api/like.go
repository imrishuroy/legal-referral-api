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

func (server *Server) likePost(ctx *gin.Context) {
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

	log.Printf("Post ID: %d", postID32)

	arg := db.LikePostParams{
		UserID: authPayload.UID,
		PostID: &postID32,
	}

	err = server.store.LikePost(ctx, arg)
	if err != nil {
		log.Error().Err(err).Msg("Error liking post")

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// create a data map
	data := map[string]string{
		"user_id": authPayload.UID,
		"post_id": postIDStr,
	}

	//Convert the map to a JSON string
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling data")
	}

	jsonString := string(jsonData)

	//server.publishToKafka("likes", authPayload.UID, string(postID32))
	server.publishToKafka("likes", authPayload.UID, jsonString)
}

func (server *Server) unlikePost(ctx *gin.Context) {
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

	err = server.store.UnlikePost(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
}

func (server *Server) listPostLikedUsers(ctx *gin.Context) {
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

	users, err := server.store.ListPostLikedUsers2(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (server *Server) likeComment(ctx *gin.Context) {
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

	err = server.store.LikeComment(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
}

func (server *Server) unlikeComment(ctx *gin.Context) {
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

	err = server.store.UnlikeComment(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
}
