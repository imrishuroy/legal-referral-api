package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
)

type saveLicenseRequest struct {
	UserID        string `json:"user_id" binding:"required"`
	Name          string `json:"name" binding:"required"`
	LicenseNumber string `json:"license_number" binding:"required"`
	IssueDate     string `json:"issue_date" binding:"required"`
	IssueState    string `json:"issue_state" binding:"required"`
}

func (server *Server) saveLicense(ctx *gin.Context) {
	var req db.SaveLicenseParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	arg := db.SaveLicenseParams{
		UserID:        req.UserID,
		Name:          req.Name,
		LicenseNumber: req.LicenseNumber,
		IssueDate:     req.IssueDate,
		IssueState:    req.IssueState,
	}

	license, err := server.store.SaveLicense(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// update the wizard step

	wizardStepArg := db.UpdateUserWizardStepParams{
		UserID:     req.UserID,
		WizardStep: 2,
	}

	_, err = server.store.UpdateUserWizardStep(ctx, wizardStepArg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, license)
}

type uploadLicenseRequest struct {
	UserId     string `json:"user_id" binding:"required"`
	LicensePdf string `json:"license_pdf" binding:"required"`
}

func (server *Server) uploadLicense(ctx *gin.Context) {
	var req uploadLicenseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	uploadLicenseArg := db.UploadLicenseParams{
		UserID:     req.UserId,
		LicensePdf: &req.LicensePdf,
	}

	_, err := server.store.UploadLicense(ctx, uploadLicenseArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// update the wizard step
	wizardStepArg := db.UpdateUserWizardStepParams{
		UserID:     req.UserId,
		WizardStep: 3,
	}

	_, err = server.store.UpdateUserWizardStep(ctx, wizardStepArg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "License uploaded successfully"})
}