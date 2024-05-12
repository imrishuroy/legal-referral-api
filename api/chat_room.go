package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

func (server *Server) listChatRooms(ctx *gin.Context) {

	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)

	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	chatRooms, err := server.store.ListChatRooms(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": errorResponse(err)})
		return
	}

	ctx.JSON(http.StatusOK, chatRooms)
}

type createChatRoomReq struct {
	RoomID  string `json:"room_id" binding:"required"`
	User1ID string `json:"user1_id" binding:"required"`
	User2ID string `json:"user2_id" binding:"required"`
}

func (server *Server) createChatRoom(ctx *gin.Context) {
	var req createChatRoomReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != req.User1ID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.CreateChatRoomParams{
		RoomID:  req.RoomID,
		User1ID: req.User1ID,
		User2ID: req.User2ID,
	}

	chatRoom, err := server.store.CreateChatRoom(ctx, arg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create chat room")
		// {"level":"error","error":"ERROR: duplicate key value violates unique constraint \"chat_rooms_pkey\" (SQLSTATE 23505)","time":"2024-05-12T14:32:31+05:30","message":"Failed to create chat room"}
		//handle this error
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {

			chatRoom, err := server.store.GetChatRoom(ctx, req.RoomID)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusOK, chatRoom)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, chatRoom)
}
