package api

import (
	"errors"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
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

type fetchUserProfileRes struct {
	User  db.User    `json:"user"`
	Price db.Pricing `json:"price"`
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

	ctx.JSON(http.StatusOK, profile)
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
