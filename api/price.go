package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"net/http"
)

type addPriceReq struct {
	ServiceType      string         `json:"service_type" binding:"required"`
	PerHourPrice     pgtype.Numeric `json:"per_hour_price"`
	PerHearingPrice  pgtype.Numeric `json:"per_hearing_price"`
	ContingencyPrice *string        `json:"contingency_price"`
	HybridPrice      *string        `json:"hybrid_price"`
}

func (server *Server) addPrice(ctx *gin.Context) {
	var req addPriceReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid json body"})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	priceArg := db.AddPriceParams{
		UserID:           authPayload.UID,
		ServiceType:      req.ServiceType,
		PerHourPrice:     req.PerHourPrice,
		PerHearingPrice:  req.PerHearingPrice,
		ContingencyPrice: req.ContingencyPrice,
		HybridPrice:      req.HybridPrice,
	}

	price, err := server.store.AddPrice(ctx, priceArg)
	if err != nil {
		log.Error().Err(err).Msg("failed to add price")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(200, price)
}
