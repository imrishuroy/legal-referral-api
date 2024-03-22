package api

import (
	"errors"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/imrishuroy/legal-referral/util"
	"net/http"
)

const provider = "provider"

func signupMiddleware(firebaseAuth *auth.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		provider := ctx.GetHeader(provider)
		if provider == "" {
			err := errors.New("provider header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse(err))
			ctx.Abort()
			return
		}
		signUpMethod, err := getSignUpMethod(provider)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid provider"})
			ctx.Abort()
			return
		}
		switch signUpMethod {
		case Email:
			//if err := VerifyAuth0AccessToken(config); err != nil {
			//	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			//	ctx.Abort()
			//	return
			//}
			ctx.Next()
		case Google, Apple:
			// Call firebaseAuthMiddleware as middleware
			firebaseAuthMiddlewareHandler := firebaseAuthMiddleware(firebaseAuth)
			firebaseAuthMiddlewareHandler(ctx)
		default:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid provider"})
			ctx.Abort()
			return
		}
	}
}

func authMiddleware(firebaseAuth *auth.Client, config util.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		provider := ctx.GetHeader(provider)
		if provider == "" {
			err := errors.New("provider header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse(err))
			ctx.Abort()
			return
		}
		signUpMethod, err := getSignUpMethod(provider)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid provider"})
			ctx.Abort()
			return
		}
		switch signUpMethod {
		case Email:
			if err := VerifyAuth0AccessToken(config); err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				ctx.Abort()
				return
			}
			ctx.Next()
		case Google, Apple:
			// Call firebaseAuthMiddleware as middleware
			firebaseAuthMiddlewareHandler := firebaseAuthMiddleware(firebaseAuth)
			firebaseAuthMiddlewareHandler(ctx)
		default:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid provider"})
			ctx.Abort()
			return
		}
	}
}

func getSignUpMethod(provider string) (SignUpMethod, error) {
	switch provider {
	case "email":
		return Email, nil
	case "google":
		return Google, nil
	case "apple":
		return Apple, nil
	default:
		return -1, errors.New("invalid provider")
	}
}
