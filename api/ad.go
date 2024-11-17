package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/rs/zerolog/log"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type AdType string

const (
	AdTypeImage AdType = "image"
	AdTypeVideo AdType = "video"
)

type PaymentCycle string

const (
	PaymentCycleWeekly  PaymentCycle = "weekly"
	PaymentCycleMonthly PaymentCycle = "monthly"
)

type createAdReq struct {
	Title        string                  `form:"title" binding:"required"`
	Description  string                  `form:"description"`
	Link         string                  `form:"link"`
	AuthorID     string                  `form:"author_id" binding:"required"`
	PaymentCycle PaymentCycle            `form:"payment_cycle" binding:"required"`
	AdType       AdType                  `form:"ad_type" binding:"required"`
	Files        []*multipart.FileHeader `form:"files"`
	StartDate    *time.Time              `form:"start_date"`
	EndDate      *time.Time              `form:"end_date"`
}

func (server *Server) createAd(ctx *gin.Context) {

	var req createAdReq

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != req.AuthorID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	imageUrls := make([]string, 0)

	if req.AdType == AdTypeImage || req.AdType == AdTypeVideo {
		//	var bucketName string
		//	if req.AdType == AdTypeImage {
		//		bucketName = "post-images"
		//	} else {
		//		bucketName = "post-videos"
		//	}

		urls, err := server.handleFilesUpload(req.Files)

		//urls, err := server.handleFilesUpload(req.Files, bucketName)
		if err != nil {
			log.Error().Msgf("Error uploading files: %v", err)
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		log.Info().Msgf("URLs: %+v", urls)
		imageUrls = append(imageUrls, urls...)

	}

	arg := db.CreateAdParams{
		Title:        req.Title,
		Description:  req.Description,
		Link:         req.Link,
		AuthorID:     req.AuthorID,
		AdType:       db.AdType(req.AdType),
		PaymentCycle: db.PaymentCycle(req.PaymentCycle),
		Media:        imageUrls,
		StartDate:    *req.StartDate,
		EndDate:      *req.EndDate,
	}

	_, err := server.store.CreateAd(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Ad created successfully"})
}

func (server *Server) listPlayingAds(ctx *gin.Context) {

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ads, err := server.store.ListPlayingAds(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, ads)
}

func (server *Server) listExpiredAds(ctx *gin.Context) {

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ads, err := server.store.ListExpiredAds(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, ads)
}

func (server *Server) extendAdPeriod(ctx *gin.Context) {
	adIDStr := ctx.Param("ad_id")
	adID, err := strconv.Atoi(adIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ad ID"})
		return
	}

	var req db.ExtendAdPeriodParams

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	req.AdID = int32(adID)

	ad, err := server.store.ExtendAdPeriod(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, ad)
}
