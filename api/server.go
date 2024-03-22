package api

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"google.golang.org/api/option"

	"github.com/imrishuroy/legal-referral/util"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Server struct {
	config       util.Config
	store        db.Store
	router       *gin.Engine
	firebaseAuth *auth.Client
}

func NewServer(config util.Config, store db.Store) (*Server, error) {

	opt := option.WithCredentialsFile("./service-account-key.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal().Msg("Failed to create Firebase app")
	}

	fmt.Println("fb connection done ", app)

	firebaseAuth, err := app.Auth(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Firebase auth client")
	}

	server := &Server{config: config, store: store, firebaseAuth: firebaseAuth}
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

func ErrorResponse(err error) gin.H {
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
