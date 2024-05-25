package api

import (
	"errors"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
	"strconv"
)

func (server *Server) createProposal(ctx *gin.Context) {
	var req *db.CreateProposalParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != req.UserID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	proposal, err := server.store.CreateProposal(ctx, *req)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}
	ctx.JSON(200, proposal)
}

func (server *Server) updateProposal(ctx *gin.Context) {
	var req *db.UpdateProposalParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != req.UserID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	proposal, err := server.store.UpdateProposal(ctx, *req)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}
	ctx.JSON(200, proposal)
}

func (server *Server) getProposal(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	referralID := ctx.Param("referral_id")
	proposalID, err := strconv.Atoi(referralID)

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.GetProposalParams{
		ReferralID: int32(proposalID),
		UserID:     userID,
	}

	proposal, err := server.store.GetProposal(ctx, arg)
	if err != nil {

		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(200, nil)
			return
		}
		ctx.JSON(400, errorResponse(err))
		return
	}
	ctx.JSON(200, proposal)
}
