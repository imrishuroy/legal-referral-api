package api

import (
	"errors"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
)

type addSocialReq struct {
	EntityType string `json:"entity_type" binding:"required"`
	Platform   string `json:"platform" binding:"required"`
	Link       string `json:"link" binding:"required"`
}

func (server *Server) addSocial(ctx *gin.Context) {
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

	social, err := server.store.AddSocial(ctx, arg)
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
