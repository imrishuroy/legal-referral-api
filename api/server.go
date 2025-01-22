package api

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/imrishuroy/legal-referral/chat"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	//"github.com/redis/go-redis/v9"

	"github.com/imrishuroy/legal-referral/util"
	//"github.com/redis/go-redis/v9"
	"github.com/twilio/twilio-go"
	"google.golang.org/api/option"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// RedisClient is a common interface for both standalone and cluster clients
//type RedisClient interface {
//	Ping(ctx context.Context) *redis.StatusCmd
//	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
//	Get(ctx context.Context, key string) *redis.StringCmd
//	MGet(ctx context.Context, keys ...string) *redis.SliceCmd
//	Pipeline() redis.Pipeliner
//}

type Server struct {
	config       util.Config
	store        db.Store
	Router       *gin.Engine
	firebaseAuth *auth.Client
	twilioClient *twilio.RestClient
	awsSession   *session.Session
	svc          *s3.S3
	//rdb          RedisClient
	hub      *chat.Hub
	producer *kafka.Producer
}

//func NewServer(config util.Config, store db.Store, hub *chat.Hub, producer *kafka.Producer, rdb RedisClient) (*Server, error) {

func NewServer(config util.Config, store db.Store, hub *chat.Hub, producer *kafka.Producer, ginLambda *ginadapter.GinLambda) (*Server, error) {

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

	log.Info().Msg("Firebase auth client created")

	var twilioClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.TwilioAccountSID,
		Password: config.TwilioAuthToken,
	})

	log.Info().Msg("Twilio client created")

	awsSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.AWSRegion),
		Credentials: credentials.NewStaticCredentials(config.AWSAccessKeyID, config.AWSSecretKey, ""),
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to create AWS session")
		return nil, err
	}

	log.Info().Msg("AWS session created")

	// s3 session
	svc := s3.New(awsSession)

	server := &Server{
		config:       config,
		store:        store,
		firebaseAuth: firebaseAuth,
		twilioClient: twilioClient,
		awsSession:   awsSession,
		svc:          svc,
		//rdb:          rdb,
		hub:      hub,
		producer: producer,
	}

	log.Info().Msg("Server created")

	//server.setupRouter(ginLambda)

	log.Info().Msg("Router setup done")
	return server, nil
}

// Start HTTP server on a specific address
//func (server *Server) Start(address string) error {
//	return server.router.Run(address)
//}

//func successResponse() gin.H {
//	return gin.H{"result": "success"}
//}

func errorResponse(err error) gin.H {
	return gin.H{"message": err.Error()}
}

func (server *Server) ping(ctx *gin.Context) {
	ctx.JSON(200, "OK")
}
