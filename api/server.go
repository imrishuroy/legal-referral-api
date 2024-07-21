package api

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/imrishuroy/legal-referral/chat"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/twilio/twilio-go"
	"google.golang.org/api/option"

	"github.com/imrishuroy/legal-referral/util"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Server struct {
	config       util.Config
	store        db.Store
	router       *gin.Engine
	firebaseAuth *auth.Client
	twilioClient *twilio.RestClient
	awsSession   *session.Session
	//redis        *redis.Client
	hub      *chat.Hub
	producer *kafka.Producer
}

func NewServer(config util.Config, store db.Store, hub *chat.Hub, producer *kafka.Producer) (*Server, error) {

	//func NewServer(config util.Config, store db.Store, hub *chat.Hub) (*Server, error) {

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

	var twilioClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.TwilioAccountSID,
		Password: config.TwilioAuthToken,
	})

	awsSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.AWSRegion),
		Credentials: credentials.NewStaticCredentials(config.AWSAccessKeyID, config.AWSSecretKey, ""),
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to create AWS session")
		return nil, err
	}

	server := &Server{
		config:       config,
		store:        store,
		firebaseAuth: firebaseAuth,
		twilioClient: twilioClient,
		awsSession:   awsSession,
		//redis:        redis,
		hub:      hub,
		producer: producer,
	}

	server.setupRouter()
	return server, nil
}

// Start HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

//func successResponse() gin.H {
//	return gin.H{"result": "success"}
//}

func errorResponse(err error) gin.H {
	return gin.H{"message": err.Error()}
}

func (server *Server) ping(ctx *gin.Context) {
	ctx.JSON(200, "UP")
}
