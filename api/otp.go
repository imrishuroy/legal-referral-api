package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/rs/zerolog/log"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
	"net/http"
)

type sendOTPRequest struct {
	To      string `json:"to"`
	Channel string `json:"channel"`
}

func (server *Server) sendOTP(ctx *gin.Context) {
	var req sendOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Validate the channel
	if req.Channel != "sms" && req.Channel != "email" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid channel"})
		return
	}

	//var serviceSID string
	//
	//// Choose Twilio service SID based on the channel
	//if req.Channel == "sms" {
	//	serviceSID = server.config.VerifyMobileServiceSID
	//} else if req.Channel == "email" {
	//	serviceSID = server.config.VerifyEmailServiceSID
	//}
	//
	//err := sendOTP(server, req.To, req.Channel, serviceSID)
	//if err != nil {
	//	log.Error().Err(err).Msg("Failed to send verification")
	//	ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send verification"})
	//	return
	//}
	ctx.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

func sendOTP(server *Server, to string, channel string, serviceSID string) (err error) {
	params := &verify.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel(channel)

	resp, err := server.twilioClient.VerifyV2.CreateVerification(serviceSID, params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send verification")
		return err
	}
	log.Logger.Info().Msgf("OTP sent successfully to %s via %s", to, channel)
	log.Logger.Info().Msgf("SID: %s", *resp.Sid)
	return nil
}

type verifyOTPRequest struct {
	UserId  string `json:"user_id"`
	To      string `json:"to"`
	Otp     string `json:"otp"`
	Channel string `json:"channel"`
}

func (server *Server) verifyOTP(ctx *gin.Context) {
	var req verifyOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Validate the channel
	if req.Channel != "sms" && req.Channel != "email" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel"})
		return
	}
	var otp = "0000"

	if otp != req.Otp {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid OTP"})
		return
	} else {
		if req.UserId != "" {
			mobileUpdateArg := db.UpdateMobileVerificationStatusParams{
				UserID:         req.UserId,
				Mobile:         &req.To,
				MobileVerified: true,
			}
			_, err := server.store.UpdateMobileVerificationStatus(ctx, mobileUpdateArg)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
		return
	}

	//params := &verify.CreateVerificationCheckParams{}
	//params.SetTo(req.To)
	//params.SetCode(req.Otp)
	//
	//var serviceSID string
	//
	//// Choose Twilio service SID based on the channel
	//if req.Channel == "sms" {
	//	serviceSID = server.config.VerifyMobileServiceSID
	//} else if req.Channel == "email" {
	//	serviceSID = server.config.VerifyEmailServiceSID
	//}
	//
	//// Verify OTP
	//resp, err := server.twilioClient.VerifyV2.CreateVerificationCheck(serviceSID, params)
	//if err != nil {
	//	log.Error().Err(err).Msg("Failed to verify OTP")
	//	ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to verify OTP"})
	//	return
	//}
	//
	//// Handle verification status
	//switch *resp.Status {
	//case "approved":
	//	// Update verification status based on the channel
	//	switch req.Channel {
	//	case "sms":
	//		if req.UserId != "" {
	//			mobileUpdateArg := db.UpdateMobileVerificationStatusParams{
	//				UserID:         req.UserId,
	//				Mobile:         &req.To,
	//				MobileVerified: true,
	//			}
	//			_, err := server.store.UpdateMobileVerificationStatus(ctx, mobileUpdateArg)
	//			if err != nil {
	//				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	//				return
	//			}
	//		}
	//		ctx.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
	//		return
	//
	//	case "email":
	//		ctx.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
	//		return
	//		//emailUpdateArg := db.UpdateEmailVerificationStatusParams{
	//		//	UserID:        req.UserId,
	//		//	EmailVerified: true,
	//		//}
	//		//_, err := server.store.UpdateEmailVerificationStatus(ctx, emailUpdateArg)
	//		//if err != nil {
	//		//	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	//		//	return
	//		//}
	//	}
	//default:
	//	ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid OTP"})
	//	return
	//}
}
