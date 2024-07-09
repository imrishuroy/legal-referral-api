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

func (server *Server) sendMessageToDiscussion(ctx *gin.Context) {
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

	message, err := server.store.SendMessageToDiscussion(ctx, req)
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

func (server *Server) listDiscussionMessages(ctx *gin.Context) {

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
	arg := db.ListDiscussionMessages3Params{
		DiscussionID: int32(discussionID),
		Offset:       (req.Offset - 1) * req.Limit,
		Limit:        req.Limit,
	}

	messages, err := server.store.ListDiscussionMessages3(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messages)

	//var messagesResponse []discussionMessage
	//for _, m := range messages {
	//	var repliedMessage *discussionMessage
	//	if m.ReplyMessageID != nil { // Check if replied message data exists
	//		repliedMessage = &discussionMessage{
	//			MessageID:       *m.ReplyMessageID,
	//			ParentMessageID: m.ReplyParentMessageID,
	//			SenderID:        *m.ReplySenderID,
	//			SenderFirstName: *m.ReplySenderFirstName,
	//			SenderLastName:  *m.ReplySenderLastName,
	//			SenderAvatarUrl: *m.ReplySenderAvatarImage,
	//			Message:         *m.ReplyMessage,
	//			DiscussionID:    *m.ReplyDiscussionID,
	//			SentAt:          m.ReplySentAt.Time,
	//		}
	//	}
	//
	//	messagesResponse = append(messagesResponse, discussionMessage{
	//		MessageID:       m.MessageID,
	//		ParentMessageID: m.ParentMessageID,
	//		SenderID:        m.SenderID,
	//		SenderFirstName: *m.SenderFirstName,
	//		SenderLastName:  *m.SenderLastName,
	//		SenderAvatarUrl: *m.SenderAvatarImage,
	//		Message:         m.Message,
	//		DiscussionID:    m.DiscussionID,
	//		SentAt:          m.SentAt,
	//		RepliedMessage:  repliedMessage,
	//	})
	//}
	//
	//if messagesResponse == nil {
	//	messagesResponse = []discussionMessage{}
	//}
	//ctx.JSON(http.StatusOK, messagesResponse)

}
