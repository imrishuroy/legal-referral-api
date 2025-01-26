package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
	"github.com/imrishuroy/legal-referral/api"
	"github.com/imrishuroy/legal-referral/chat"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/imrishuroy/legal-referral/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// TODO: CHANGE THE ROUTER LIKE THIS https://github.com/build-on-aws/golang-gin-app-on-aws-lambda/blob/main/function/main.go

var ginLambda *ginadapter.GinLambda
var server *api.Server
var pool *pgxpool.Pool

func init() {

	log.Info().Msg("Welcome to LegalReferral")

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config : " + err.Error())
	}

	// Load connection details from environment variables
	proxyEndpoint := "legal-referral-db-proxy.proxy-ct2smiqa0pnv.us-east-1.rds.amazonaws.com"
	//dbHost := config.DBHost
	dbUser := config.DBUser         // Database username
	dbPassword := config.DBPassword // Database password
	dbName := config.DBName         // Database name

	if dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatal().Msg("Missing required environment variables")
	}

	// Construct the connection string
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=verify-full", dbUser, dbPassword, proxyEndpoint, dbName)

	log.Info().Msg("DB URL: " + dbURL)

	// Parse the pool config
	dbConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot parse db config")
	}

	// Configure pool settings (optional)
	dbConfig.MaxConns = 10
	dbConfig.MinConns = 2

	// Create the connection pool
	pool, err = pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create connection pool")
	}
	//defer pool.Close()

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("cannot connect to database")
	}

	log.Info().Msg("Connected to the database")

	// db connection
	//connPool, err := pgxpool.New(context.Background(), config.DBSource)

	//if err != nil {
	//	fmt.Println("cannot connect to db:", err)
	//}

	//fmt.Println("Connection Pool: ", connPool)

	//defer pool.Close()

	store := db.NewStore(pool)

	hub := chat.NewHub(store)
	go hub.Run()

	//setup producer
	conf := kafka.ConfigMap{
		// User-specific properties that you must set
		"bootstrap.servers": config.BootStrapServers,
		"sasl.username":     config.SASLUsername,
		"sasl.password":     config.SASLPassword,

		// Fixed properties
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "PLAIN",
		"acks":              "all"}

	producer, err := kafka.NewProducer(&conf)
	if err != nil {
		log.Error().Err(err).Msg("cannot create producer")
	}

	//ctx := context.Background()

	//rdb := GetRedisClient(config)

	//pong, err := rdb.Ping(ctx).Result()
	//if err != nil {
	//	log.Error().Err(err).Msg("cannot connect to redis")
	//} else {
	//	log.Info().Msg("Connected to Redis with TLS: " + pong)
	//}

	// api server setup
	//server, err := api.NewServer(config, store, hub, producer, rdb)
	server, err = api.NewServer(config, store, hub, producer, ginLambda)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}

	log.Info().Msg("Server created")

}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, request)
}

func ping(ctx *gin.Context) {
	ctx.JSON(200, "OK")
}

func main() {

	defer pool.Close()

	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	r.GET("/ping", ping)

	// auth routes
	r.POST("/api/sign-in", server.SignIn)
	r.POST("/api/sign-up", server.SignUp)
	r.POST("/api/refresh-token", server.RefreshToken)
	r.POST("/api/otp/send", server.SendOTP)
	r.POST("/api/otp/verify", server.VerifyOTP)
	r.POST("/api/reset-password", server.ResetPassword)
	r.GET("/api/users/:user_id/wizardstep", server.GetUserWizardStep)
	r.GET("/api/firms", server.SearchFirms)
	r.POST("/api/sign-in/linkedin", server.LinkedinLogin)

	auth := r.Group("/api").
		Use(server.AuthMiddleware(server.FirebaseAuth))

	auth.POST("/users", server.CreateUser)
	auth.GET("/users/:user_id", server.GetUserById)
	auth.POST("/license", server.SaveLicense)
	//auth.POST("/license/upload", server.uploadLicense)
	//auth.POST("/about-you", server.saveAboutYou)
	//auth.GET("/users/:user_id/profile", server.fetchUserProfile)
	//
	//auth.PUT("/users/info", server.updateUserInfo)
	//auth.POST("/review", server.addReview)
	//
	//auth.POST("/price", server.addPrice)
	//auth.PUT("/price/:price_id", server.updatePrice)
	//auth.PUT("/users/:user_id/toggle-referral", server.toggleOpenToReferral)
	//auth.PUT("/users/:user_id/banner", server.updateUserBannerImage)
	//
	//// profile/user
	//auth.PUT("/users/:user_id/avatar", server.updateUserAvatar)
	//
	//// profile/socials
	//auth.POST("/socials", server.addSocial)
	//auth.PUT("/socials/:social_id", server.updateSocial)
	//auth.GET("/socials/:entity_type/:entity_id", server.listSocials)
	//auth.DELETE("/socials/:social_id", server.deleteSocial)
	//
	//// profile/experiences
	//auth.POST("/users/:user_id/experiences", server.addExperience)
	//auth.GET("/users/:user_id/experiences", server.listExperiences)
	//auth.PUT("/users/:user_id/experiences/:experience_id", server.updateExperience)
	//auth.DELETE("/users/:user_id/experiences/:experience_id", server.deleteExperience)
	//
	//// profile/educations
	//auth.POST("/users/:user_id/educations", server.addEducation)
	//auth.GET("/users/:user_id/educations", server.listEducations)
	//auth.PUT("/users/:user_id/educations/:education_id", server.updateEducation)
	//auth.DELETE("/users/:user_id/educations/:education_id", server.deleteEducation)
	//
	//// account
	//auth.GET("/accounts/:user_id", server.getAccountInfo)
	//
	//// network
	//auth.POST("/connections/send", server.sendConnection)
	//auth.POST("/connections/:id/accept", server.acceptConnection)
	//auth.POST("/connections/:id/reject", server.rejectConnection)
	//auth.GET("/connections/invitations/:user_id", server.listConnectionInvitations)
	//auth.GET("/connections/:user_id", server.listConnections)
	//auth.GET("/recommendations/:user_id", server.listRecommendations)
	//auth.POST("/recommendations/cancel", server.cancelRecommendation)
	//auth.GET("/search/users", server.searchUsers)
	//// check if user is connected to another user
	//auth.GET("/connections/:user_id/:other_user_id", server.checkConnection)
	//
	//// chat
	//auth.GET("/chat/:room_id", func(ctx *gin.Context) {
	//	roomId := ctx.Param("room_id")
	//	chat.ServeWS(ctx, roomId, server.hub)
	//})
	//auth.GET("/chat/:room_id/messages", server.listMessages)
	//auth.GET("/chat/users/:user_id/rooms", server.listChatRooms)
	//auth.POST("/chat/rooms", server.createChatRoom)
	//
	//// referral
	//auth.POST("/referral", server.createReferral)
	//auth.GET("/referrals/:user_id/active", server.listActiveReferrals)
	//auth.GET("/referrals/users/:project_id", server.listReferredUsers)
	//auth.GET("/users/:user_id/proposals", server.listActiveProposals)
	//auth.POST("/proposals", server.createProposal)
	//auth.PUT("/proposals/:proposal_id", server.updateProposal)
	//auth.GET("users/:user_id/proposals/:project_id", server.getProposal)
	//auth.POST("/projects/award", server.awardProject)
	//auth.GET("/projects/awarded/:user_id", server.listAwardedProjects)
	//auth.PUT("/projects/:project_id/accept", server.acceptProject)
	//auth.PUT("/projects/:project_id/reject", server.rejectProject)
	//auth.GET("/projects/active/:user_id", server.listActiveProjects)
	//auth.PUT("/projects/:project_id/start", server.startProject)
	//auth.PUT("/projects/:project_id/initiate-complete", server.initiateCompleteProject)
	//auth.PUT("/projects/:project_id/cancel/initiate-complete", server.cancelInitiateCompleteProject)
	//auth.PUT("/projects/:project_id/complete", server.completeProject)
	//auth.GET("/projects/completed/:user_id", server.listCompletedProjects)
	//
	//auth.POST("projects/review", server.createProjectReview)
	//auth.GET("projects/review/:project_id", server.getProjectReview)
	//
	//auth.GET("/users/:user_id/connected", server.listConnectedUsers)
	auth.GET("/users", server.ListUsers)
	//
	//auth.GET("/users/license-verified", server.listLicenseVerifiedUsers)
	//auth.GET("/users/license-unverified", server.listLicenseUnverifiedUsers)
	//
	//// approve license
	//auth.PUT("/users/:user_id/approve-license", server.approveLicense)
	//auth.PUT("/users/:user_id/reject-license", server.rejectLicense)
	//
	//// posts
	//auth.POST("/posts", server.createPost)
	//auth.GET("/posts/:post_id", server.getPost)
	//auth.DELETE("/posts/:post_id", server.deletePost)
	//auth.GET("/search/posts", server.searchPosts)
	//auth.GET("/posts/:post_id/is-featured", server.isPostFeatured)
	//
	//// news feed
	//auth.GET("/feeds/:user_id", server.listNewsFeed)
	////auth.GET("/v2/feeds/:user_id", server.listNewsFeedV2)
	////auth.GET("/v3/feeds/:user_id", server.listNewsFeedV3)
	//
	//// like post
	//auth.POST("/posts/:post_id/like", server.likePost)
	//auth.DELETE("/posts/:post_id/like", server.unlikePost)
	//auth.GET("/posts/:post_id/liked-users", server.listPostLikedUsers)
	//
	//// get post likes and comments count
	//auth.GET("/posts/:post_id/likes-comments-count", server.postLikesAndCommentsCount)
	//auth.GET("/posts/:post_id/is-liked", server.isPostLiked)
	//
	//// comments
	//auth.POST("/posts/:post_id/comments", server.commentPost)
	//auth.GET("/posts/:post_id/comments", server.listComments)
	//auth.POST("/comments/:comment_id/like", server.likeComment)
	//auth.DELETE("/comments/:comment_id/like", server.unlikeComment)
	//
	//// discussion
	//auth.POST("/discussions", server.createDiscussion)
	//
	//// update discussion topic
	//auth.PUT("/discussions/:discussion_id/topic", server.updateDiscussionTopic)
	//auth.POST("/discussions/:discussion_id/invite", server.inviteUserToDiscussion)
	//auth.POST("/discussions/:discussion_id/join", server.joinDiscussion)
	//auth.POST("/discussions/:discussion_id/reject", server.rejectDiscussion)
	//auth.GET("/discussions/invites/:user_id", server.listDiscussionInvites)
	//auth.GET("/discussions/active/:user_id", server.listActiveDiscussions)
	//auth.GET("/discussions/:discussion_id/participants", server.listDiscussionParticipants)
	//auth.GET("/discussions/:discussion_id/uninvited", server.listUninvitedParticipants)
	//
	//// discussion messages
	//auth.POST("/discussions/:discussion_id/messages", server.sendMessageToDiscussion)
	//auth.GET("/discussions/:discussion_id/messages", server.listDiscussionMessages)
	//
	//// ads
	//auth.POST("/ads", server.createAd)
	////playing ads
	//auth.GET("/ads/playing", server.listPlayingAds)
	//auth.GET("/ads/expired", server.listExpiredAds)
	//// extend ad period
	//auth.PUT("/ads/:ad_id/extend", server.extendAdPeriod)
	//
	//// admin
	//auth.GET("/attorneys", server.listAttorneys)
	//auth.GET("/lawyers", server.listLawyers)
	//
	//auth.GET("/referrals/:user_id", server.listAllReferralProjects)
	//auth.GET("/referrals/completed/:user_id", server.listCompletedReferralProjects)
	//auth.GET("/referrals/active/:user_id", server.listActiveReferralProjects)
	//
	//auth.POST("/faqs", server.createFAQ)
	//auth.GET("/faqs", server.listFAQs)
	//
	//auth.POST("/firms", server.addFirm)
	//auth.GET("/firms/owner/:owner_user_id", server.listFirmsByOwner)
	//
	//// save post
	//auth.POST("/saved-posts", server.savePost)
	//auth.DELETE("/saved-posts/:saved_post_id", server.unSavePost)
	//auth.GET("/saved-posts/:user_id", server.listSavedPosts)
	//
	//// feature posts
	//auth.POST("/feature-posts", server.featurePost)
	//auth.DELETE("/feature-posts/:post_id", server.unFeaturePost)
	//auth.GET("/feature-posts/:user_id", server.listFeaturePosts)
	//
	//// notifications
	//auth.POST("/device-details", server.saveDevice)
	//auth.POST("/notifications", server.createNotification)
	//auth.GET("/notifications/:user_id", server.listNotifications)
	//
	//// post stats
	//auth.GET("/posts/:post_id/stats", server.getPostStats)
	//
	//// report
	//auth.POST("/report-post", server.reportPost)
	//auth.GET("/posts/:post_id/reported-status/:user_id", server.isPostReported)
	//auth.DELETE("/feeds/:feed_id/ignore", server.ignoreFeed)
	//
	//// activity
	//auth.GET("/activity/posts/:user_id", server.listActivityPosts)
	//auth.GET("/activity/comments/:user_id", server.listActivityComments)
	//auth.GET("/users/:user_id/followers-count", server.getUserFollowersCount)
	//
	//// print ginAdapter
	//log.Info().Msg("Setting up routes 3")
	//
	//log.Info().Msgf("GinAdapter: %+v", ginLambda)

	//go server.CreateConsumer(context.Background())

	// start the server
	//err = server.Start(config.ServerAddress)
	//err = server.Start(config.ServerAddress)
	//if err != nil {
	//	log.Fatal().Err(err).Msg("cannot create server:")
	//}

	//// GRAPHQL
	////auth.POST("/query", gin.WrapH(srv))
	//

	ginLambda = ginadapter.New(r)
	lambda.Start(Handler)
}

//func GetRedisClient(config util.Config) api.RedisClient {
//	redisURL := fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort)
//	if config.Env == "prod" {
//		clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
//			Addrs:           []string{redisURL},
//			Password:        "",
//			PoolSize:        50,
//			MinIdleConns:    20,
//			DialTimeout:     3 * time.Second,
//			ReadTimeout:     1 * time.Second,
//			WriteTimeout:    1 * time.Second,
//			PoolTimeout:     2 * time.Second,
//			MaxRetries:      3,
//			MinRetryBackoff: 8 * time.Millisecond,
//			MaxRetryBackoff: 256 * time.Millisecond,
//			TLSConfig: &tls.Config{
//				InsecureSkipVerify: false,
//			},
//			ReadOnly:       false,
//			RouteByLatency: true,  // Prioritize low-latency nodes
//			RouteRandomly:  false, // Avoid random routing to improve predictability
//		})
//		return clusterClient
//	} else {
//		client := redis.NewClient(&redis.Options{
//			Addr:     redisURL,
//			Password: "",
//			DB:       0,
//		})
//		return client
//	}
//}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
