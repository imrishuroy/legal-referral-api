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

type commentPostReq struct {
	UserID          string `json:"user_id"`
	SenderID        string `json:"sender_id"`
	PostId          int    `json:"post_id"`
	Content         string `json:"content"`
	ParentCommentId *int32 `json:"parent_comment_id"`
}

func (srv *Server) CommentPost(ctx *gin.Context) {

	var req commentPostReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != req.SenderID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.CommentPostParams{
		UserID:          req.UserID,
		PostID:          int32(req.PostId),
		Content:         req.Content,
		ParentCommentID: req.ParentCommentId,
	}

	postIDStr := strconv.Itoa(req.PostId)

	comment, err := srv.Store.CommentPost(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Prepare notification data
	data := map[string]string{
		"user_id":           req.UserID,
		"sender_id":         req.SenderID,
		"target_id":         postIDStr,
		"target_type":       "comment",
		"notification_type": "like",
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

	ctx.JSON(http.StatusOK, comment)

}

func (srv *Server) ListComments(ctx *gin.Context) {

	postIDStr := ctx.Param("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid post ID"})
		return
	}
	postID32 := int32(postID)

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.ListComments2Params{
		PostID: postID32,
		UserID: authPayload.UID,
	}

	comments, err := srv.Store.ListComments2(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, comments)
}
