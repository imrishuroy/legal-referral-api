package api

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/imrishuroy/legal-referral/util"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

//var (
//	// The signing key for the token.
//	signingKey = []byte("WCoaOyM7KJ4TCfeVziyPgW5RnO6qW6zY")
//
//	// Our token must be signed using this data.
//	keyFunc = func(ctx context.Context) (interface{}, error) {
//		return signingKey, nil
//	}
//)

// Validate if the posts scope is present in the token.
func (c CustomClaims) Validate(ctx context.Context) error {
	if c.Scope == "read:posts" {
		return errors.New("scope is required")
	}

	return nil

}

// HasScope checks whether our claims have a specific scope.
func (c CustomClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}

	return false
}

func VerifyAccessToken(config util.Config) gin.HandlerFunc {

	// The signing key for the token.
	signingKey := []byte(config.SigningKey)

	// Our token must be signed using this data.
	keyFunc := func(ctx context.Context) (interface{}, error) {
		return signingKey, nil
	}

	issuerURL, err := url.Parse("https://" + config.Auth0Domain + "/")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse Auth0 domain")
	}

	tokenValidator, err := validator.New(
		keyFunc,
		validator.HS256,
		issuerURL.String(),
		[]string{config.Auth0Audience},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{
					Scope: "read:posts",
				}
			},
		),
		// TODO: check this what it does
		validator.WithAllowedClockSkew(30*time.Second),
	)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create validator")
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Encountered error while validating JWT: %v", err)
	}

	middleware := jwtmiddleware.New(
		tokenValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(ctx *gin.Context) {
		encounteredError := true
		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			encounteredError = false
			ctx.Request = r
			ctx.Next()
		}

		middleware.CheckJWT(handler).ServeHTTP(ctx.Writer, ctx.Request)

		if encounteredError {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "Failed to validate JWT."})
		}
	}

}
