package api

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	//"github.com/aws/aws-sdk-go/aws/credentials"

	//"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/credentials"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/imrishuroy/legal-referral/chat"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	//"github.com/redis/go-redis/v9"

	"github.com/imrishuroy/legal-referral/util"
	//"github.com/redis/go-redis/v9"
	"github.com/twilio/twilio-go"
	"google.golang.org/api/option"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

var (
	s3Client *s3.Client
)

type Server struct {
	config util.Config
	Store  db.Store
	//Router       *gin.Engine
	FirebaseAuth *auth.Client
	twilioClient *twilio.RestClient
	awsSession   *session.Session
	//SVC          *s3.S3
	S3Client *s3.Client
	//rdb          RedisClient
	Hub      *chat.Hub
	producer *kafka.Producer
}

//func NewServer(config util.Config, Store db.Store, hub *chat.Hub, producer *kafka.Producer, rdb RedisClient) (*Server, error) {

func NewServer(con util.Config, store db.Store, hub *chat.Hub, producer *kafka.Producer, ginLambda *ginadapter.GinLambda) (*Server, error) {

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
		Username: con.TwilioAccountSID,
		Password: con.TwilioAuthToken,
	})

	log.Info().Msg("Twilio client created")

	//awsSession, err := session.NewSession(&aws.Config{
	//	Region:      aws.String(config.AWSRegion),
	//	Credentials: credentials.NewStaticCredentials(config.AWSAccessKeyID, config.AWSSecretKey, ""),
	//})

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Error().Err(err).Msg("Failed to load AWS config")
		return nil, err
	}

	//if err != nil {
	//	log.Error().Err(err).Msg("Failed to create AWS session")
	//	return nil, err
	//}

	s3Client = s3.NewFromConfig(cfg)

	log.Info().Msg("AWS session created")

	// s3 session
	//svc := s3.New(awsSession)

	/// print svc
	fmt.Println("svc", s3Client)
	if s3Client == nil {
		fmt.Println("svc is nil")
	} else {
		fmt.Println("svc is not nil")
	}

	server := &Server{
		config:       con,
		Store:        store,
		FirebaseAuth: firebaseAuth,
		twilioClient: twilioClient,
		//awsSession:   awsSession,
		//SVC:          svc,
		//rdb:          rdb,
		S3Client: s3Client,
		Hub:      hub,
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

func (srv *Server) ping(ctx *gin.Context) {
	ctx.JSON(200, "OK")
}
