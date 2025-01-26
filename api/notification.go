package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
)

type createNotificationReq struct {
	UserID           string `json:"user_id" binding:"required"`
	SenderID         string `json:"sender_id" binding:"required"`
	TargetID         int32  `json:"target_id" binding:"required"`
	TargetType       string `json:"target_type" binding:"required"`
	NotificationType string `json:"notification_type" binding:"required"`
	Message          string `json:"message" binding:"required"`
}

func (s *Server) createNotification(ctx *gin.Context) {

	var req createNotificationReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != req.SenderID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	args := db.CreateNotificationParams{
		UserID:           req.UserID,
		SenderID:         req.SenderID,
		TargetID:         req.TargetID,
		TargetType:       req.TargetType,
		NotificationType: req.NotificationType,
		Message:          req.Message,
	}

	notification, err := s.store.CreateNotification(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, notification)

}

type listNotificationsReq struct {
	Limit  int32 `form:"limit"`
	Offset int32 `form:"offset"`
}

func (s *Server) listNotifications(ctx *gin.Context) {

	var req listNotificationsReq
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

	arg := db.ListNotificationsParams{
		UserID: userID,
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	notifications, err := s.store.ListNotifications(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, notifications)

}
