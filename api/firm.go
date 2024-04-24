package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
)

type addFirmReq struct {
	Name     string `json:"name" binding:"required"`
	LogoUrl  string `json:"logo_url" binding:"required"`
	OrgType  string `json:"org_type" binding:"required"`
	Website  string `json:"website" binding:"required"`
	Location string `json:"location" binding:"required"`
	About    string `json:"about" binding:"required"`
}

func (server *Server) addFirm(ctx *gin.Context) {
	var req addFirmReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	arg := db.AddFirmParams{
		Name:     req.Name,
		LogoUrl:  req.LogoUrl,
		OrgType:  req.OrgType,
		Website:  req.Website,
		Location: req.Location,
		About:    req.About,
	}

	company, err := server.store.AddFirm(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, company)
}

type listFirmReq struct {
	Limit       int32  `form:"limit" binding:"required"`
	Offset      int32  `form:"offset" binding:"required"`
	SearchQuery string `form:"query"`
}

func (server *Server) listFirms(ctx *gin.Context) {

	var req listFirmReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListFirmsParams{
		Limit:  req.Limit,
		Offset: (req.Offset - 1) * req.Limit,
		Query:  req.SearchQuery,
	}

	firms, err := server.store.ListFirms(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, firms)
}
