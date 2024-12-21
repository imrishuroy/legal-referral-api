package api

import (
	"errors"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"mime/multipart"
	"net/http"
)

func (server *Server) updateUserAvatar(ctx *gin.Context) {

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error parsing form"})
		return
	}

	files := form.File["file"]
	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "No file uploaded"})
		return
	}

	file, err := files[0].Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error opening file"})
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	fileName := generateUniqueFilename() + getFileExtension(files[0])
	url, err := server.uploadFile(file, fileName, files[0].Header.Get("Content-Type"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error uploading file"})
		return
	}

	arg := db.UpdateUserAvatarParams{
		UserID:    userID,
		AvatarUrl: &url,
	}

	err = server.store.UpdateUserAvatar(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	log.Info().Msgf("avatar url %v", url)

	ctx.String(http.StatusOK, url)

}

type toggleOpenToReferralReq struct {
	OpenToReferral bool `json:"open_to_referral"`
}

func (server *Server) toggleOpenToReferral(ctx *gin.Context) {

	var req toggleOpenToReferralReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}
	log.Info().Msgf("open to referral %v", req.OpenToReferral)

	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.ToggleOpenToRefferalParams{
		UserID:         userID,
		OpenToReferral: req.OpenToReferral,
	}

	err := server.store.ToggleOpenToRefferal(ctx, arg)
	if err != nil {
		log.Err(err).Msg("error changing open to referral")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

type userProfile struct {
	UserID                  string         `json:"user_id"`
	FirstName               string         `json:"first_name"`
	LastName                string         `json:"last_name"`
	PracticeArea            *string        `json:"practice_area"`
	AvatarUrl               *string        `json:"avatar_url"`
	BannerUrl               *string        `json:"banner_url"`
	AverageBillingPerClient *int32         `json:"average_billing_per_client"`
	CaseResolutionRate      *int32         `json:"case_resolution_rate"`
	OpenToReferral          bool           `json:"open_to_referral"`
	About                   *string        `json:"about"`
	PriceID                 *int64         `json:"price_id"`
	ServiceType             *string        `json:"service_type"`
	PerHourPrice            pgtype.Numeric `json:"per_hour_price"`
	PerHearingPrice         pgtype.Numeric `json:"per_hearing_price"`
	ContingencyPrice        *string        `json:"contingency_price"`
	HybridPrice             *string        `json:"hybrid_price"`
	RatingInfo              *ratingInfo    `json:"rating_info"`
	FollowersCount          int64          `json:"followers_count"`
	ConnectionsCount        int64          `json:"connections_count"`
}

func (server *Server) fetchUserProfile(ctx *gin.Context) {

	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	profile, err := server.store.FetchUserProfile(ctx, userID)
	log.Error().Err(err).Msg("error fetching user profile")
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusOK, profile)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	avgRating, ok := profile.AverageRating.(float64)
	if !ok {
		avgRating = 0.0
	}

	attorneys, ok := profile.Attorneys.(int64)
	if !ok {
		attorneys = 0
	}

	followers, ok := profile.FollowersCount.(int64)
	if !ok {
		followers = 0
	}

	connections, ok := profile.ConnectionsCount.(int64)
	if !ok {
		connections = 0
	}

	userProfile := userProfile{
		UserID:                  profile.UserID,
		FirstName:               profile.FirstName,
		LastName:                profile.LastName,
		PracticeArea:            profile.PracticeArea,
		AvatarUrl:               profile.AvatarUrl,
		BannerUrl:               profile.BannerUrl,
		AverageBillingPerClient: profile.AverageBillingPerClient,
		CaseResolutionRate:      profile.CaseResolutionRate,
		OpenToReferral:          profile.OpenToReferral,
		About:                   profile.About,
		PriceID:                 profile.PriceID,
		ServiceType:             profile.ServiceType,
		PerHourPrice:            profile.PerHourPrice,
		PerHearingPrice:         profile.PerHearingPrice,
		ContingencyPrice:        profile.ContingencyPrice,
		HybridPrice:             profile.HybridPrice,
		RatingInfo: &ratingInfo{
			AverageRating: avgRating,
			Attorneys:     attorneys,
		},
		FollowersCount:   followers,
		ConnectionsCount: connections,
	}

	ctx.JSON(http.StatusOK, userProfile)
}

func (server *Server) updateUserBannerImage(ctx *gin.Context) {

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error parsing form"})
		return
	}

	files := form.File["file"]
	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "No file uploaded"})
		return
	}

	file, err := files[0].Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error opening file"})
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	userID := ctx.Param("user_id")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID != userID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	fileName := generateUniqueFilename() + getFileExtension(files[0])

	url, err := server.uploadFile(file, fileName, files[0].Header.Get("Content-Type"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error uploading file"})
		return
	}

	arg := db.UpdateUserBannerImageParams{
		UserID:    userID,
		BannerUrl: &url,
	}

	err = server.store.UpdateUserBannerImage(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	log.Info().Msgf("banner url %v", url)

	// send the url as string
	ctx.String(http.StatusOK, url)
}
