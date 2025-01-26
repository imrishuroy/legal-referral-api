package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/rs/zerolog/log"
	"mime/multipart"
	"net/http"
)

type socialReq struct {
	Platform string `json:"platform" binding:"required"`
	Link     string `json:"link" binding:"required"`
}

type addFirmReq struct {
	Name        string                  `form:"name" binding:"required"`
	OwnerUserID string                  `form:"owner_user_id" binding:"required"`
	Files       []*multipart.FileHeader `form:"file"`
	OrgType     string                  `form:"org_type" binding:"required"`
	Website     string                  `form:"website" binding:"required"`
	Location    string                  `form:"location" binding:"required"`
	About       string                  `form:"about" binding:"required"`
}

func (s *Server) addFirm(ctx *gin.Context) {

	var req addFirmReq

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	log.Info().Msgf("Request: %+v", req)

	// Check if the authenticated user is the owner
	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != req.OwnerUserID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	urls, err := s.handleFilesUpload(req.Files)

	if err != nil && len(urls) == 0 {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.AddFirmParams{
		Name:        req.Name,
		OwnerUserID: req.OwnerUserID,
		LogoUrl:     urls[0],
		OrgType:     req.OrgType,
		Website:     req.Website,
		Location:    req.Location,
		About:       req.About,
	}

	firm, err := s.store.AddFirm(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, firm)
}

type searchFirmsReq struct {
	Limit       int32  `form:"limit" binding:"required"`
	Offset      int32  `form:"offset" binding:"required"`
	SearchQuery string `form:"query"`
}

func (s *Server) SearchFirms(ctx *gin.Context) {

	var req searchFirmsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListFirmsParams{
		Limit:  req.Limit,
		Offset: (req.Offset - 1) * req.Limit,
		Query:  req.SearchQuery,
	}

	firms, err := s.store.ListFirms(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, firms)
}

type listFirmsByOwnerReq struct {
	OwnerUserID string `uri:"owner_user_id" binding:"required"`
}

func (s *Server) listFirmsByOwner(ctx *gin.Context) {

	var req listFirmsByOwnerReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Err(err).Msg("error binding uri")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	firms, err := s.store.ListFirmsByOwner(ctx, "YLFPbwsDBqOpMNdP3C04GC6iEdW2")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, firms)

}
