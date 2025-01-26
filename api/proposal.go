package api

import (
	"errors"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
	"strconv"
)

func (s *Server) createProposal(ctx *gin.Context) {
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

	proposal, err := s.store.CreateProposal(ctx, *req)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}
	ctx.JSON(200, proposal)
}

func (s *Server) updateProposal(ctx *gin.Context) {
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

	proposal, err := s.store.UpdateProposal(ctx, *req)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}
	ctx.JSON(200, proposal)
}

func (s *Server) getProposal(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	projectIDStr := ctx.Param("project_id")
	projectID, err := strconv.Atoi(projectIDStr)

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.GetProposalParams{
		ProjectID: int32(projectID),
		UserID:    userID,
	}

	proposal, err := s.store.GetProposal(ctx, arg)
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
