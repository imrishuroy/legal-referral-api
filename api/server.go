package api

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
	"github.com/imrishuroy/legal-referral/chat"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/imrishuroy/legal-referral/util"
	"github.com/rs/zerolog/log"
	"github.com/twilio/twilio-go"
	"github.com/valkey-io/valkey-go"
	"google.golang.org/api/option"
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
	Config        util.Config
	Store         db.Store
	FirebaseAuth  *auth.Client
	TwilioClient  *twilio.RestClient
	S3Client      *s3.Client
	Hub           *chat.Hub
	KafkaProducer *kafka.Producer
	ValkeyClient  valkey.Client
	SQS           *sqs.SQS
}

func NewServer(con util.Config, store db.Store, hub *chat.Hub, producer *kafka.Producer, valkeyClient valkey.Client, sqs *sqs.SQS) (*Server, error) {

	opt := option.WithCredentialsFile("./service-account-key.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal().Msg("Failed to create Firebase app")
	}

	firebaseAuth, err := app.Auth(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Firebase auth client")
	}

	var twilioClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: con.TwilioAccountSID,
		Password: con.TwilioAuthToken,
	})

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Error().Err(err).Msg("Failed to load AWS Config")
		return nil, err
	}

	s3Client = s3.NewFromConfig(cfg)

	server := &Server{
		Config:        con,
		Store:         store,
		FirebaseAuth:  firebaseAuth,
		TwilioClient:  twilioClient,
		S3Client:      s3Client,
		Hub:           hub,
		KafkaProducer: producer,
		ValkeyClient:  valkeyClient,
		SQS:           sqs,
	}

	return server, nil
}

func successResponse() gin.H {
	return gin.H{"result": "success"}
}

func errorResponse(err error) gin.H {
	return gin.H{"message": err.Error()}
}

func (srv *Server) ping(ctx *gin.Context) {
	ctx.JSON(200, "OK")
}
