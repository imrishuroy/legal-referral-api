package api

import (
	"errors"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
	"strconv"
)

type addSocialReq struct {
	EntityType string `json:"entity_type" binding:"required"`
	Platform   string `json:"platform" binding:"required"`
	Link       string `json:"link" binding:"required"`
}

func (s *Server) addSocial(ctx *gin.Context) {
	var req addSocialReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.AddSocialParams{
		EntityID:   authPayload.UID,
		EntityType: req.EntityType,
		Platform:   req.Platform,
		Link:       req.Link,
	}

	social, err := s.store.AddSocial(ctx, arg)
	if err != nil {
		// check if error the duplicate key error is returned
		if !errors.Is(err, db.ErrUniqueViolation) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Social link already exists"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, social)
}

type updateSocialReq struct {
	Platform string `json:"platform" binding:"required"`
	Link     string `json:"link" binding:"required"`
}

func (s *Server) updateSocial(ctx *gin.Context) {
	var req updateSocialReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	socialIDParam := ctx.Param("social_id")
	socialID, err := strconv.ParseInt(socialIDParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid entity id"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.UpdateSocialParams{
		SocialID: socialID,
		Platform: req.Platform,
		Link:     req.Link,
	}

	social, err := s.store.UpdateSocial(ctx, arg)
	if err != nil {
		if !errors.Is(err, db.ErrUniqueViolation) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Social link already exists"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, social)
}

func (s *Server) listSocials(ctx *gin.Context) {

	entityID := ctx.Param("entity_id")
	entityType := ctx.Param("entity_type")

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)

	if authPayload.UID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}

	arg := db.ListSocialsParams{
		EntityID:   entityID,
		EntityType: entityType,
	}

	socials, err := s.store.ListSocials(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, socials)
}

func (s *Server) deleteSocial(ctx *gin.Context) {

	socialIDParam := ctx.Param("social_id")
	socialID, err := strconv.ParseInt(socialIDParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid entity id"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}

	err = s.store.DeleteSocial(ctx, socialID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
