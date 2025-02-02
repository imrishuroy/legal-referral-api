package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

func (s *Server) ListChatRooms(ctx *gin.Context) {

	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)

	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	chatRooms, err := s.Store.ListChatRooms(ctx, userID)
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

func (s *Server) CreateChatRoom(ctx *gin.Context) {
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

	// Check if chat room already exists
	getChatRoomArg := db.GetChatRoomParams{
		RoomID:  req.RoomID,
		User1ID: req.User1ID,
	}

	chatRoom, err := s.Store.GetChatRoom(ctx, getChatRoomArg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get chat room")
		if strings.Contains(err.Error(), "no rows in result set") {

			chatRoom, err := s.Store.CreateChatRoom(ctx, arg)
			if err != nil {
				log.Error().Err(err).Msg("Failed to create chat room")
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusOK, chatRoom)
			return

		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, chatRoom)

}
