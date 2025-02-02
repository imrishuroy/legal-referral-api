package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
)

func (s *Server) ListLicenseVerifiedUsers(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var req listUsersReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	arg := db.ListLicenseVerifiedUsersParams{
		Limit:  req.Limit,
		Offset: (req.Offset - 1) * req.Limit,
	}

	users, err := s.Store.ListLicenseVerifiedUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (s *Server) ListLicenseUnverifiedUsers(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var req listUsersReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	arg := db.ListLicenseUnVerifiedUsersParams{
		Limit:  req.Limit,
		Offset: (req.Offset - 1) * req.Limit,
	}

	users, err := s.Store.ListLicenseUnVerifiedUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

type listAttorneysReq struct {
	Limit  int32 `form:"limit"`
	Offset int32 `form:"offset"`
}

func (s *Server) ListAttorneys(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var req listAttorneysReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	arg := db.ListAttorneysParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	users, err := s.Store.ListAttorneys(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// lawyers

func (s *Server) ListLawyers(ctx *gin.Context) {

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	lawyers, err := s.Store.ListLawyers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, lawyers)
}

// referrals

func (s *Server) ListAllReferralProjects(ctx *gin.Context) {
	userId := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	referrals, err := s.Store.ListAllReferralProjects(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, referrals)
}

func (s *Server) ListCompletedReferralProjects(ctx *gin.Context) {
	userId := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	referrals, err := s.Store.ListCompletedReferralProjects(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, referrals)
}

func (s *Server) ListActiveReferralProjects(ctx *gin.Context) {
	userId := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	referrals, err := s.Store.ListActiveReferralProjects(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, referrals)
}
