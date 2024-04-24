package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

type resetPasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (server *Server) resetPassword(ctx *gin.Context) {
	var req resetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	if req.Email == "" || req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Email and Password are required"})
		return
	}

	user, err := server.firebaseAuth.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// authorization
	// authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	//if authPayload.UID != user.UID {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
	//	return
	//}

	// delete the user
	err = server.firebaseAuth.DeleteUser(ctx, user.UID)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete user")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// create the user
	userArg := &auth.UserToCreate{}
	userArg.Email(req.Email)
	userArg.EmailVerified(false)
	userArg.Password(req.Password)
	userArg.UID(user.UID)

	_, err = server.firebaseAuth.CreateUser(ctx, userArg)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
	return
}
