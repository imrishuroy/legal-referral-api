package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/golang-jwt/jwt/v5"
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

// Validate if the posts scope is present in the token.
func (c CustomClaims) Validate(ctx context.Context) error {
	log.Info().Msg("Validating the id token")
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

	issuerURL, err := url.Parse("https://" + config.Auth0Domain + "/")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse Auth0 domain")
	}
	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	tokenValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{config.Auth0Audience},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{
					Scope: "read:posts",
				}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
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

		middleware.
			CheckJWT(handler).ServeHTTP(ctx.Writer, ctx.Request)

		if encounteredError {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "Failed to validate JWT."})
		}
	}
}

func ExtractEmailFromIDToken(ctx *gin.Context) string {
	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return ""
	}
	bearerToken := authHeader[7:]
	if bearerToken == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token is required"})
		return "Error getting token"
	}

	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return "", nil
	})
	if err != nil {

	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		email, ok := claims["email"].(string)
		if !ok {
			return ""
		}
		return email
	}
	return ""
}

//func VerifyIDToken(ctx *gin.Context, config util.Config) string {
//
//	authHeader := ctx.Request.Header.Get("Authorization")
//	if authHeader == "" {
//		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
//		return ""
//	}
//	bearerToken := authHeader[7:]
//	if bearerToken == "" {
//		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token is required"})
//		return ""
//	}
//	log.Info().Msgf("Bearer token: %s", bearerToken)
//
//	// The signing key for the token.
//	signingKey := []byte(config.SigningKey)
//
//	// Our token must be signed using this data.
//	//keyFunc := func(token *jwt.Token) (interface{}, error) {
//	//	return signingKey, nil
//	//}
//
//	token1, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
//		signedString, err := token.SignedString(signingKey)
//		if err != nil {
//			return nil, err
//		}
//		log.Info().Msgf("Signed string: %s", signedString)
//		token.Signature = signedString
//
//		sign, err := token.Method.Sign(signedString, signingKey)
//		if err != nil {
//			return nil, err
//		}
//		log.Info().Msgf("Sign: %s", sign)
//
//		//if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
//		//	return nil, errors.New("Unexpected signing method")
//		//}
//		return "Email", nil
//	})
//	if err != nil {
//		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
//		return ""
//	}
//
//	// get the email from the token
//	claims, ok := token1.Claims.(jwt.MapClaims)
//	if !ok {
//		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to get claims from token"})
//		return ""
//	}
//	email, ok := claims["email"].(string)
//	if !ok {
//		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to get email from token"})
//		return ""
//	}
//	return email
//}

//func VerifyIdToken(ctx *gin.Context, config util.Config) string {
//	authHeader := ctx.Request.Header.Get("Authorization")
//	if authHeader == "" {
//		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
//		return ""
//	}
//	bearerToken := authHeader[7:]
//	if bearerToken == "" {
//		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token is required"})
//		return "Error getting token"
//	}
//
//	toValidate := map[string]string{}
//	toValidate["nonce"] = "8v1LGYS8WZ3CnMcDxG-RzAaGT_wn7KKHqCQbpRmZAFY="
//	toValidate["aud"] = "YXTkygnU2SepreFXxlY5THnX5Vz3EuwN"
//
//	issuerURL, err := url.Parse("https://" + config.Auth0Domain + "/")
//
//	jwtVerifierSetup := jwtverifier.JwtVerifier{
//		Issuer:           issuerURL.String(),
//		ClaimsToValidate: toValidate,
//	}
//
//	verifier, error := jwtVerifierSetup.New()
//	if error != nil {
//		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": error.Error()})
//		return ""
//	}
//
//	token, err := verifier.VerifyIdToken(bearerToken)
//	if err != nil {
//		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
//		return ""
//	}
//
//	sub := token.Claims["sub"]
//	log.Info().Msgf("Sub: %s", sub)
//	return fmt.Sprintf("%v", sub)
//}
