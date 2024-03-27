package api

import (
	"crypto/rand"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/rs/zerolog/log"
	"math/big"
	"net/http"
	"time"
)

type sendOTPRequest struct {
	UserId  string `json:"user_id"`
	Channel string `json:"channel"`
}

func (server *Server) sendOTP(ctx *gin.Context) {

	var req sendOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	otp, err := generateOTP(6)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	log.Logger.Info().Msgf("OTP: %d", otp)

	arg := db.StoreOTPParams{
		UserID:  req.UserId,
		Channel: req.Channel,
		Otp:     otp,
	}
	sessionId, err := server.store.StoreOTP(ctx, arg)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"session_id": sessionId})

}

type verifyOTPRequest struct {
	SessionId int64  `json:"session_id"`
	UserId    string `json:"user_id"`
	Channel   string `json:"channel"`
	Otp       int32  `json:"otp"`
}

func (server *Server) verifyOTP(ctx *gin.Context) {
	var req verifyOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	otp, err := server.store.GetOTP(ctx, req.SessionId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	expireTime := otp.CreatedAt.Add(5 * time.Minute)
	if time.Now().After(expireTime) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "OTP expired"})
		return
	}

	if otp.SessionID != req.SessionId || otp.Otp != req.Otp || otp.Channel != req.Channel || otp.UserID != req.UserId {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	err = server.store.DeleteOTP(ctx, req.SessionId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": "otp verified successfully"})

}

func generateOTP(length int) (int32, error) {
	// Define the set of characters allowed in the OTP
	chars := "0123456789"
	otp := make([]byte, length)

	// Generate random indices to select characters from the set
	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return 0, err
		}
		otp[i] = chars[num.Int64()]
	}

	// Convert the OTP bytes to an int32
	var result int32
	for _, digit := range otp {
		result = result*10 + int32(digit-'0')
	}

	return result, nil
}
