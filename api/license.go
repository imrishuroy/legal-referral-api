package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"mime/multipart"
	"net/http"
)

type saveLicenseRequest struct {
	Name          string      `json:"name" binding:"required"`
	LicenseNumber string      `json:"license_number" binding:"required"`
	IssueDate     pgtype.Date `json:"issue_date" binding:"required"`
	IssueState    string      `json:"issue_state" binding:"required"`
}

func (srv *Server) SaveLicense(ctx *gin.Context) {
	var req saveLicenseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind JSON")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.SaveLicenseParams{
		UserID:        authPayload.UID,
		Name:          req.Name,
		LicenseNumber: req.LicenseNumber,
		IssueDate:     req.IssueDate,
		IssueState:    req.IssueState,
	}

	license, err := srv.Store.SaveLicense(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// update the wizard step
	wizardStepArg := db.UpdateUserWizardStepParams{
		UserID:     authPayload.UID,
		WizardStep: 1,
	}

	_, err = srv.Store.UpdateUserWizardStep(ctx, wizardStepArg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, license)
}

func (srv *Server) UploadLicense(ctx *gin.Context) {

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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	fileName := generateUniqueFilename() + getFileExtension(files[0])
	url, err := srv.uploadFile(ctx, file, fileName, files[0].Header.Get("Content-Type"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error uploading file"})
		return
	}

	uploadLicenseArg := db.UploadLicenseParams{
		UserID:     authPayload.UID,
		LicenseUrl: &url,
	}
	_, err = srv.Store.UploadLicense(ctx, uploadLicenseArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// update the wizard step
	wizardStepArg := db.UpdateUserWizardStepParams{
		UserID:     authPayload.UID,
		WizardStep: 2,
	}

	_, err = srv.Store.UpdateUserWizardStep(ctx, wizardStepArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "License uploaded successfully"})
}
