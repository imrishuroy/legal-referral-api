package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
	"strconv"
	"time"
)

type discussionMessage struct {
	MessageID       int32              `json:"message_id"`
	ParentMessageID *int32             `json:"parent_message_id"`
	SenderID        string             `json:"sender_id"`
	SenderAvatarUrl string             `json:"sender_avatar_url"`
	SenderFirstName string             `json:"sender_first_name"`
	SenderLastName  string             `json:"sender_last_name"`
	Message         string             `json:"message"`
	DiscussionID    int32              `json:"discussion_id"`
	SentAt          time.Time          `json:"sent_at"`
	RepliedMessage  *discussionMessage `json:"replied_message"`
}

func (srv *Server) SendMessageToDiscussion(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}

	var req db.SendMessageToDiscussionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	message, err := srv.Store.SendMessageToDiscussion(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, message)
}

type listDiscussionMessages struct {
	Offset int32 `form:"offset" binding:"required"`
	Limit  int32 `form:"limit" binding:"required"`
}

func (srv *Server) ListDiscussionMessages(ctx *gin.Context) {

	var req listDiscussionMessages

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	discussionIDStr := ctx.Param("discussion_id")
	discussionID, err := strconv.Atoi(discussionIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}
	arg := db.ListDiscussionMessagesParams{
		DiscussionID: int32(discussionID),
		Offset:       (req.Offset - 1) * req.Limit,
		Limit:        req.Limit,
	}

	messages, err := srv.Store.ListDiscussionMessages(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messages)

}
