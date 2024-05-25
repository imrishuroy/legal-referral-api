package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
	"strconv"
)

type addReferralRequest struct {
	ReferredUserIDs           []string `json:"referred_user_ids"`
	ReferrerUserID            string   `json:"referrer_user_id"`
	Title                     string   `json:"title"`
	PreferredPracticeArea     string   `json:"preferred_practice_area"`
	PreferredPracticeLocation string   `json:"preferred_practice_location"`
	CaseDescription           string   `json:"case_description"`
}

func (server *Server) addReferral(ctx *gin.Context) {
	var req addReferralRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != req.ReferrerUserID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	//var referrals []db.Referral
	for _, referredUserID := range req.ReferredUserIDs {
		arg := db.CreateReferralParams{
			ReferredUserID:            referredUserID,
			ReferrerUserID:            req.ReferrerUserID,
			Title:                     req.Title,
			PreferredPracticeArea:     req.PreferredPracticeArea,
			PreferredPracticeLocation: req.PreferredPracticeLocation,
			CaseDescription:           req.CaseDescription,
		}

		_, err := server.store.CreateReferral(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		//referrals = append(referrals, referral)
	}

	ctx.String(http.StatusCreated, "Referral created")
}

func (server *Server) listActiveReferrals(ctx *gin.Context) {

	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	referrals, err := server.store.ListActiveReferrals(ctx, authPayload.UID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, referrals)
}

func (server *Server) listReferredUsers(ctx *gin.Context) {
	referralIDStr := ctx.Param("referral_id")
	referralID, err := strconv.ParseInt(referralIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	referredUsers, err := server.store.ListReferredUsers(ctx, int32(referralID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, referredUsers)
}

func (server *Server) listProposals(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	proposals, err := server.store.ListProposals(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, proposals)

}
