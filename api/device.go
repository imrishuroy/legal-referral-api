package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
)

type saveDeviceReq struct {
	DeviceID    string `json:"device_id" binding:"required"`
	DeviceToken string `json:"device_token" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
}

func (s *Server) SaveDevice(ctx *gin.Context) {
	var req saveDeviceReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	arg := db.SaveDeviceParams{
		DeviceID:    req.DeviceID,
		DeviceToken: req.DeviceToken,
		UserID:      req.UserID,
	}

	err := s.Store.SaveDevice(ctx, arg)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Device saved successfully"})

}
