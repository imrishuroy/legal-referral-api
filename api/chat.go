package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
	"time"
)

type listMessagesRequest struct {
	Offset int32 `form:"offset" binding:"required"`
	Limit  int32 `form:"limit" binding:"required"`
}

type message struct {
	MessageID       int32     `json:"message_id"`
	ParentMessageID *int32    `json:"parent_message_id"`
	SenderID        string    `json:"sender_id"`
	RecipientID     string    `json:"recipient_id"`
	Message         string    `json:"message"`
	HasAttachment   bool      `json:"has_attachment"`
	AttachmentID    *int32    `json:"attachment_id"`
	IsRead          bool      `json:"is_read"`
	RoomID          string    `json:"room_id"`
	SentAt          time.Time `json:"sent_at"`
	RepliedMessage  *message  `json:"replied_message"`
}

func (srv *Server) ListMessages(ctx *gin.Context) {
	var req listMessagesRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	roomID := ctx.Param("room_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}
	arg := db.ListMessagesParams{
		RoomID: roomID,
		Offset: (req.Offset - 1) * req.Limit,
		Limit:  req.Limit,
	}

	messages, err := srv.Store.ListMessages(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var messagesResponse []message
	for _, m := range messages {
		var repliedMessage *message
		if m.MessageID_2 != nil { // Check if replied message data exists
			repliedMessage = &message{
				MessageID:       *m.MessageID_2,
				ParentMessageID: m.ParentMessageID_2,
				SenderID:        *m.SenderID_2,
				RecipientID:     *m.RecipientID_2,
				Message:         *m.Message_2,
				HasAttachment:   *m.HasAttachment_2,
				AttachmentID:    m.AttachmentID_2,
				IsRead:          *m.IsRead_2,
				RoomID:          *m.RoomID_2,
				SentAt:          m.SentAt_2.Time,
			}
		}

		messagesResponse = append(messagesResponse, message{
			MessageID:       m.MessageID,
			ParentMessageID: m.ParentMessageID,
			SenderID:        m.SenderID,
			RecipientID:     m.RecipientID,
			Message:         m.Message,
			HasAttachment:   m.HasAttachment,
			AttachmentID:    m.AttachmentID,
			IsRead:          m.IsRead,
			RoomID:          m.RoomID,
			SentAt:          m.SentAt,
			RepliedMessage:  repliedMessage,
		})
	}

	if messagesResponse == nil {
		messagesResponse = []message{}
	}
	ctx.JSON(http.StatusOK, messagesResponse)
}

//func CreateRoomID(userID1, userID2 string) string {
//	// Concatenate the user IDs
//	concatenated := userID1 + userID2
//
//	// Hash the concatenated string using SHA-1
//	hash := sha1.New()
//	hash.Write([]byte(concatenated))
//	hashed := hash.Sum(nil)
//
//	// Convert the hashed bytes to a hexadecimal string
//	roomID := hex.EncodeToString(hashed)
//
//	return roomID
//}
