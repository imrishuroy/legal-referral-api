package api

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/rs/zerolog/log"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
	"net/http"
	"time"
)

type sendMobileOTPRequest struct {
	Mobile string `json:"mobile"`
}

func (server *Server) sendMobileOTP(ctx *gin.Context) {
	var req sendMobileOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	params := &openapi.CreateVerificationParams{}
	params.SetTo(req.Mobile)
	params.SetChannel("sms")

	// ctx.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})

	resp, err := server.twilioClient.VerifyV2.CreateVerification(server.config.VerifyServiceSID, params)

	if err != nil {
		fmt.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send verification"})
	} else {
		log.Logger.Info().Msgf("OTP sent successfully to %s", req.Mobile)
		log.Logger.Info().Msgf("SID: %s", *resp.Sid)
		ctx.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
	}
}

type verifyMobileOTPRequest struct {
	UserId string `json:"user_id"`
	Mobile string `json:"mobile"`
	Otp    string `json:"otp"`
}

func (server *Server) verifyMobileOTP(ctx *gin.Context) {
	var req verifyMobileOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := &openapi.CreateVerificationCheckParams{}
	params.SetTo(req.Mobile)
	params.SetCode(req.Otp)

	resp, err := server.twilioClient.VerifyV2.CreateVerificationCheck(server.config.VerifyServiceSID, params)
	log.Err(err).Msg("Error while verifying OTP")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to verify OTP"})
		return
	} else if *resp.Status == "approved" {

		user, err := server.store.GetUserById(ctx, req.UserId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		updateUserArg := db.UpdateUserParams{
			ID:               user.ID,
			FirstName:        user.FirstName,
			LastName:         user.LastName,
			Mobile:           req.Mobile,
			IsEmailVerified:  user.IsEmailVerified,
			IsMobileVerified: true,
			WizardStep:       1,
			WizardCompleted:  user.WizardCompleted,
		}

		_, err = server.store.UpdateUser(ctx, updateUserArg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
		return
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid OTP"})
		return
	}
}

func sendEmailOTP(email string, store db.Store, ctx *gin.Context) (int64, error) {

	otp, err := generateOTP()
	if err != nil {
		return 0, err
	}
	// if len of otp is less than 4 digits, generate a new otp
	for otp < 1000 {
		otp, err = generateOTP()
		if err != nil {
			return 0, err
		}
	}

	storeOtpArg := db.StoreOTPParams{
		Email:   email,
		Channel: "email",
		Otp:     otp,
	}

	sessionId, err := store.StoreOTP(ctx, storeOtpArg)
	if err != nil {
		return 0, err
	}

	return sessionId, nil

}

type sendEmailOTPRequest struct {
	Email   string `json:"email"`
	Channel string `json:"channel"`
}

func (server *Server) sendEmailOTP(ctx *gin.Context) {

	var req sendEmailOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	sessionId, err := sendEmailOTP(req.Email, server.store, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"session_id": sessionId})

}

type verifyEmailOTPRequest struct {
	SessionId int64  `json:"session_id"`
	Email     string `json:"email"`
	Channel   string `json:"channel"`
	Otp       int32  `json:"otp"`
}

func (server *Server) verifyEmailOTP(ctx *gin.Context) {
	var req verifyEmailOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	otp, err := server.store.GetOTP(ctx, req.SessionId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid OTP"})
		return
	}

	expireTime := otp.CreatedAt.Add(5 * time.Minute)
	if time.Now().After(expireTime) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "OTP expired"})
		return
	}

	if otp.SessionID != req.SessionId || otp.Otp != req.Otp || otp.Channel != req.Channel || otp.Email != req.Email {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid OTP"})
		return
	}

	_ = server.store.DeleteOTP(ctx, req.SessionId)

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	updateUserArg := db.UpdateUserParams{
		ID:               user.ID,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		Mobile:           user.Mobile,
		IsEmailVerified:  true,
		IsMobileVerified: user.IsMobileVerified,
		WizardStep:       user.WizardStep,
		WizardCompleted:  user.WizardCompleted,
	}

	_, err = server.store.UpdateUser(ctx, updateUserArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "otp verified successfully"})

}

// Function to generate a random 4-digit OTP
func generateOTP() (int32, error) {
	// Initialize a byte slice to store random bytes
	randBytes := make([]byte, 4)

	// Read random bytes from the crypto/rand package
	_, err := rand.Read(randBytes)
	if err != nil {
		return 0, err
	}

	// Convert random bytes to an integer
	otp := binary.BigEndian.Uint32(randBytes)

	// Ensure the OTP is exactly 4 digits
	otp = otp % 10000

	// Return the OTP
	return int32(otp), nil

}
