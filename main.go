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

var ginLambda *ginadapter.GinLambda

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, request)
}

func ping(ctx *gin.Context) {
	ctx.JSON(200, "OK")
}

func main() {

	log.Info().Msg("Welcome to Legal Referral API")

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config : " + err.Error())
	}

	// db connection
	pool, err := pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		fmt.Println("cannot connect to db:", err)
	}
	defer pool.Close()

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("cannot connect to database")
	}

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

	//rdb := GetRedisClient(config)

	//pong, err := rdb.Ping(ctx).Result()
	//if err != nil {
	//	log.Error().Err(err).Msg("cannot connect to redis")
	//} else {
	//	log.Info().Msg("Connected to Redis with TLS: " + pong)
	//}

	// api server setup
	server, err := api.NewServer(config, store, hub, producer, ginLambda)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}

	log.Info().Msg("Server created")

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
	auth.POST("/license/upload", server.UploadLicense)
	auth.POST("/about-you", server.SaveAboutYou)
	auth.GET("/users/:user_id/profile", server.FetchUserProfile)
	auth.PUT("/users/info", server.UpdateUserInfo)
	auth.POST("/review", server.AddReview)

	auth.POST("/price", server.AddPrice)
	auth.PUT("/price/:price_id", server.UpdatePrice)
	auth.PUT("/users/:user_id/toggle-referral", server.ToggleOpenToReferral)
	auth.PUT("/users/:user_id/banner", server.UpdateUserBannerImage)

	//// profile/user
	auth.PUT("/users/:user_id/avatar", server.UpdateUserAvatar)

	//// profile/socials
	auth.POST("/socials", server.AddSocial)
	auth.PUT("/socials/:social_id", server.UpdateSocial)
	auth.GET("/socials/:entity_type/:entity_id", server.ListSocials)
	auth.DELETE("/socials/:social_id", server.DeleteSocial)

	//// profile/experiences
	auth.POST("/users/:user_id/experiences", server.AddExperience)
	auth.GET("/users/:user_id/experiences", server.ListExperiences)
	auth.PUT("/users/:user_id/experiences/:experience_id", server.UpdateExperience)
	auth.DELETE("/users/:user_id/experiences/:experience_id", server.DeleteExperience)

	//// profile/educations
	auth.POST("/users/:user_id/educations", server.AddEducation)
	auth.GET("/users/:user_id/educations", server.ListEducations)
	auth.PUT("/users/:user_id/educations/:education_id", server.UpdateEducation)
	auth.DELETE("/users/:user_id/educations/:education_id", server.DeleteEducation)

	//// account
	auth.GET("/accounts/:user_id", server.GetAccountInfo)

	//// network
	auth.POST("/connections/send", server.SendConnection)
	auth.POST("/connections/:id/accept", server.AcceptConnection)
	auth.POST("/connections/:id/reject", server.RejectConnection)
	auth.GET("/connections/invitations/:user_id", server.ListConnectionInvitations)
	auth.GET("/connections/:user_id", server.ListConnections)
	auth.GET("/recommendations/:user_id", server.ListRecommendations)
	auth.POST("/recommendations/cancel", server.CancelRecommendation)
	auth.GET("/search/users", server.SearchUsers)
	// check if user is connected to another user
	auth.GET("/connections/:user_id/:other_user_id", server.CheckConnection)

	//// chat
	auth.GET("/chat/:room_id", func(ctx *gin.Context) {
		roomId := ctx.Param("room_id")
		chat.ServeWS(ctx, roomId, server.Hub)
	})
	auth.GET("/chat/:room_id/messages", server.ListMessages)
	auth.GET("/chat/users/:user_id/rooms", server.ListChatRooms)
	auth.POST("/chat/rooms", server.CreateChatRoom)

	//// referral
	auth.POST("/referral", server.CreateReferral)
	auth.GET("/referrals/:user_id/active", server.ListActiveReferrals)
	auth.GET("/referrals/users/:project_id", server.ListReferredUsers)
	auth.GET("/users/:user_id/proposals", server.ListActiveProposals)
	auth.POST("/proposals", server.CreateProposal)
	auth.PUT("/proposals/:proposal_id", server.UpdateProposal)
	auth.GET("users/:user_id/proposals/:project_id", server.GetProposal)
	auth.POST("/projects/award", server.AwardProject)
	auth.GET("/projects/awarded/:user_id", server.ListAwardedProjects)
	auth.PUT("/projects/:project_id/accept", server.AcceptProject)
	auth.PUT("/projects/:project_id/reject", server.RejectProject)
	auth.GET("/projects/active/:user_id", server.ListActiveProjects)
	auth.PUT("/projects/:project_id/start", server.StartProject)
	auth.PUT("/projects/:project_id/initiate-complete", server.InitiateCompleteProject)
	auth.PUT("/projects/:project_id/cancel/initiate-complete", server.CancelInitiateCompleteProject)
	auth.PUT("/projects/:project_id/complete", server.CompleteProject)
	auth.GET("/projects/completed/:user_id", server.ListCompletedProjects)
	auth.POST("projects/review", server.CreateProjectReview)
	auth.GET("projects/review/:project_id", server.GetProjectReview)
	//
	auth.GET("/users/:user_id/connected", server.ListConnectedUsers)
	auth.GET("/users", server.ListUsers)
	auth.GET("/users/license-verified", server.ListLicenseVerifiedUsers)
	auth.GET("/users/license-unverified", server.ListLicenseUnverifiedUsers)

	//// approve license
	auth.PUT("/users/:user_id/approve-license", server.ApproveLicense)
	auth.PUT("/users/:user_id/reject-license", server.RejectLicense)

	//// posts
	auth.POST("/posts", server.CreatePost)
	auth.GET("/posts/:post_id", server.GetPost)
	auth.DELETE("/posts/:post_id", server.DeletePost)
	auth.GET("/search/posts", server.SearchPosts)
	auth.GET("/posts/:post_id/is-featured", server.IsPostFeatured)

	//// news feed
	auth.GET("/feeds/:user_id", server.ListNewsFeed)
	////auth.GET("/v2/feeds/:user_id", server.listNewsFeedV2)
	////auth.GET("/v3/feeds/:user_id", server.listNewsFeedV3)

	//// like post
	auth.POST("/posts/:post_id/like", server.LikePost)
	auth.DELETE("/posts/:post_id/like", server.UnlikePost)
	auth.GET("/posts/:post_id/liked-users", server.ListPostLikedUsers)

	//// get post likes and comments count
	auth.GET("/posts/:post_id/likes-comments-count", server.PostLikesAndCommentsCount)
	auth.GET("/posts/:post_id/is-liked", server.IsPostLiked)

	//// comments
	auth.POST("/posts/:post_id/comments", server.CommentPost)
	auth.GET("/posts/:post_id/comments", server.ListComments)
	auth.POST("/comments/:comment_id/like", server.LikeComment)
	auth.DELETE("/comments/:comment_id/like", server.UnlikeComment)

	//// discussion
	auth.POST("/discussions", server.CreateDiscussion)

	//// update discussion topic
	auth.PUT("/discussions/:discussion_id/topic", server.UpdateDiscussionTopic)
	auth.POST("/discussions/:discussion_id/invite", server.InviteUserToDiscussion)
	auth.POST("/discussions/:discussion_id/join", server.JoinDiscussion)
	auth.POST("/discussions/:discussion_id/reject", server.RejectDiscussion)
	auth.GET("/discussions/invites/:user_id", server.ListDiscussionInvites)
	auth.GET("/discussions/active/:user_id", server.ListActiveDiscussions)
	auth.GET("/discussions/:discussion_id/participants", server.ListDiscussionParticipants)
	auth.GET("/discussions/:discussion_id/uninvited", server.ListUninvitedParticipants)

	//// discussion messages
	auth.POST("/discussions/:discussion_id/messages", server.SendMessageToDiscussion)
	auth.GET("/discussions/:discussion_id/messages", server.ListDiscussionMessages)

	//// ads
	auth.POST("/ads", server.CreateAd)
	//playing ads
	auth.GET("/ads/playing", server.ListPlayingAds)
	auth.GET("/ads/expired", server.ListExpiredAds)
	// extend ad period
	auth.PUT("/ads/:ad_id/extend", server.ExtendAdPeriod)

	//// admin
	auth.GET("/attorneys", server.ListAttorneys)
	auth.GET("/lawyers", server.ListLawyers)

	auth.GET("/referrals/:user_id", server.ListAllReferralProjects)
	auth.GET("/referrals/completed/:user_id", server.ListCompletedReferralProjects)
	auth.GET("/referrals/active/:user_id", server.ListActiveReferralProjects)

	auth.POST("/faqs", server.CreateFAQ)
	auth.GET("/faqs", server.ListFAQs)

	auth.POST("/firms", server.AddFirm)
	auth.GET("/firms/owner/:owner_user_id", server.ListFirmsByOwner)

	//// save post
	auth.POST("/saved-posts", server.SavePost)
	auth.DELETE("/saved-posts/:saved_post_id", server.UnSavePost)
	auth.GET("/saved-posts/:user_id", server.ListSavedPosts)

	//// feature posts
	auth.POST("/feature-posts", server.FeaturePost)
	auth.DELETE("/feature-posts/:post_id", server.UnFeaturePost)
	auth.GET("/feature-posts/:user_id", server.ListFeaturePosts)

	//// notifications
	auth.POST("/device-details", server.SaveDevice)
	auth.POST("/notifications", server.CreateNotification)
	auth.GET("/notifications/:user_id", server.ListNotifications)

	//// post stats
	auth.GET("/posts/:post_id/stats", server.GetPostStats)

	//// report
	auth.POST("/report-post", server.ReportPost)
	auth.GET("/posts/:post_id/reported-status/:user_id", server.IsPostReported)
	auth.DELETE("/feeds/:feed_id/ignore", server.IgnoreFeed)

	//// activity
	auth.GET("/activity/posts/:user_id", server.ListActivityPosts)
	auth.GET("/activity/comments/:user_id", server.ListActivityComments)
	auth.GET("/users/:user_id/followers-count", server.GetUserFollowersCount)

	// start the server
	//err = server.Start(config.ServerAddress)
	//err = server.Start(config.ServerAddress)
	//if err != nil {
	//	log.Fatal().Err(err).Msg("cannot create server:")
	//}

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
