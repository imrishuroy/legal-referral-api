package api

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"net/http"
)

type searchUserRequest struct {
	Query  string `form:"query" binding:"required"`
	Filter string `form:"filter" binding:"required"`
	Limit  int32  `form:"limit" binding:"required"`
	Offset int32  `form:"offset" binding:"required"`
}

func (s *Server) searchUsers(ctx *gin.Context) {
	var req searchUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Token)
	if authPayload.UID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}

	switch req.Filter {
	case "All":
		users, err := s.store.SearchAllUsers(ctx, req.Query)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, users)
		return

	case "1st":
		args := db.Search1stDegreeConnectionsParams{
			CurrentUserID: authPayload.UID,
			Query:         req.Query,
		}
		users, err := s.store.Search1stDegreeConnections(ctx, args)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, users)
		return

	case "2nd":
		args := db.Search2ndDegreeConnectionsParams{
			CurrentUserID: authPayload.UID,
			Query:         req.Query,
		}
		users, err := s.store.Search2ndDegreeConnections(ctx, args)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, users)
		return

	default:
		users := make([]db.User, 0)
		ctx.JSON(http.StatusOK, users)
		return
	}

}
