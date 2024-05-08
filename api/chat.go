package api

import (
	"crypto/sha1"
	"encoding/hex"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

type listMessagesRequest struct {
	RecipientID string `json:"recipient_id"`
}

func (server *Server) listMessages(ctx *gin.Context) {

	roomID := ctx.Param("room_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}
	messages, err := server.store.ListMessages(ctx, roomID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, messages)
}

func CreateRoomID(userID1, userID2 string) string {
	// Concatenate the user IDs
	concatenated := userID1 + userID2

	// Hash the concatenated string using SHA-1
	hash := sha1.New()
	hash.Write([]byte(concatenated))
	hashed := hash.Sum(nil)

	// Convert the hashed bytes to a hexadecimal string
	roomID := hex.EncodeToString(hashed)

	return roomID
}
