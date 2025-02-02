package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func (s *Server) AuthMiddleware(auth *auth.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		log.Info().Msgf("Authorization header: %s", authorizationHeader)
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		log.Info().Msgf("Fields: %v", fields)

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authrorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		log.Info().Msgf("Authorization type: %s", authorizationType)
		accessToken := fields[1]
		log.Info().Msgf("Access token: %s", accessToken)
		idToken, err := auth.VerifyIDToken(context.Background(), accessToken)
		if err != nil {
			log.Info().Msgf("Error verifying ID token: %v", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		log.Info().Msgf("ID token: %v", idToken)
		ctx.Set(authorizationPayloadKey, idToken)
		ctx.Next()

	}
}
