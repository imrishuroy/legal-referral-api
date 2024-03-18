package api

import (
	db "github.com/imrishuroy/legal-referral/db/sqlc"

	"github.com/imrishuroy/legal-referral/util"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config util.Config
	store  db.Store
	router *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {

	server := &Server{config: config, store: store}
	server.setupRouter()
	return server, nil
}

// Start HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func successResponse() gin.H {
	return gin.H{"result": "success"}
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) ping(c *gin.Context) {
	c.JSON(200, "pong")
}

func (server *Server) checkScope(ctx *gin.Context) {

	token := ctx.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)

	claims := token.CustomClaims.(*CustomClaims)
	if !claims.HasScope("read:posts") {
		ctx.JSON(403, gin.H{"error": "Insufficient scope"})
		return
	}

	ctx.JSON(200, successResponse())

}
