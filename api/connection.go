package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
	"strconv"
)

type ConnectionStatus int32

const (
	pending = iota
	accepted
	rejected
)

func (s ConnectionStatus) String() string {
	return [...]string{"pending", "accepted", "rejected"}[s]
}
func (s ConnectionStatus) Int32() int32 {
	return int32(s)
}

type sendConnectionReq struct {
	SenderID    string `json:"sender_id" binding:"required"`
	RecipientID string `json:"recipient_id" binding:"required"`
}

type sendConnectionRes struct {
	ID      int32  `json:"id"`
	Message string `json:"message"`
}

func (server *Server) sendConnection(ctx *gin.Context) {
	var req sendConnectionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != req.SenderID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.SendConnectionParams{
		SenderID:    req.SenderID,
		RecipientID: req.RecipientID,
	}

	connID, err := server.store.SendConnection(ctx, arg)
	if err != nil {
		errorCode := db.ErrorCode(err)
		if errorCode == db.UniqueViolation {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Connection request already sent"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := sendConnectionRes{
		ID:      connID,
		Message: "Connection request sent",
	}

	ctx.JSON(http.StatusOK, res)
}

func (server *Server) acceptConnection(ctx *gin.Context) {
	connIDParams := ctx.Param("id")
	connID, err := strconv.ParseInt(connIDParams, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid connection ID"})
		return
	}

	arg := db.AcceptConnectionTxParams{
		ID: int32(connID),
	}

	conn, err := server.store.AcceptConnectionTx(ctx, arg)
	if err != nil {
		errorCode := db.ErrorCode(err)
		if errorCode == db.UniqueViolation {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Connection request already accepted"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, conn)
}

type rejectConnectionReq struct {
	SenderID    string `json:"sender_id" binding:"required"`
	RecipientID string `json:"recipient_id" binding:"required"`
}

func (server *Server) rejectConnection(ctx *gin.Context) {
	var req rejectConnectionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	arg := db.RejectConnectionParams{
		SenderID:    req.SenderID,
		RecipientID: req.RecipientID,
	}

	err := server.store.RejectConnection(ctx, arg)
	if err != nil {
		errorCode := db.ErrorCode(err)
		if errorCode == db.UniqueViolation {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Connection request already rejected"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Connection request rejected"})
}

type listConnectionInvitationRes struct {
	ID         int32            `json:"id"`
	Recipient  db.User          `json:"recipient"`
	Status     ConnectionStatus `json:"status"`
	CreateTime string           `json:"create_time"`
}

type listConnectionInvitationsReq struct {
	Limit  int32 `form:"limit" binding:"required"`
	Offset int32 `form:"offset" binding:"required"`
}

func (server *Server) listConnectionInvitations(ctx *gin.Context) {

	var req listConnectionInvitationsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.ListConnectionInvitationsParams{
		RecipientID: userID,
		Limit:       req.Limit,
		Offset:      (req.Offset - 1) * req.Limit,
	}

	connections, err := server.store.ListConnectionInvitations(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, connections)
}

type listConnectionsReq struct {
	Limit  int32 `form:"limit" binding:"required"`
	Offset int32 `form:"offset" binding:"required"`
}

func (server *Server) listConnections(ctx *gin.Context) {

	var req listConnectionsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.ListConnectionsParams{
		UserID: userID,
		Limit:  req.Limit,
		Offset: (req.Offset - 1) * req.Limit,
	}

	connections, err := server.store.ListConnections(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, connections)
}
