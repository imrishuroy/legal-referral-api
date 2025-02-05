package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
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

	// api srv setup
	srv, err := api.NewServer(config, store, hub, producer, ginLambda)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create srv:")
	}

	log.Info().Msg("Server created")

	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	r.GET("/ping", ping)

	// auth routes
	r.POST("/api/sign-in", srv.SignIn)
	r.POST("/api/sign-up", srv.SignUp)
	r.POST("/api/refresh-token", srv.RefreshToken)
	r.POST("/api/otp/send", srv.SendOTP)
	r.POST("/api/otp/verify", srv.VerifyOTP)
	r.POST("/api/reset-password", srv.ResetPassword)
	r.GET("/api/users/:user_id/wizardstep", srv.GetUserWizardStep)
	r.GET("/api/firms", srv.SearchFirms)
	r.POST("/api/sign-in/linkedin", srv.LinkedinLogin)

	auth := r.Group("/api").
		Use(srv.AuthMiddleware(srv.FirebaseAuth))

	auth.POST("/users", srv.CreateUser)
	auth.GET("/users/:user_id", srv.GetUserById)
	auth.POST("/license", srv.SaveLicense)
	auth.POST("/license/upload", srv.UploadLicense)
	auth.POST("/about-you", srv.SaveAboutYou)
	auth.GET("/users/:user_id/profile", srv.FetchUserProfile)
	auth.PUT("/users/info", srv.UpdateUserInfo)
	auth.POST("/review", srv.AddReview)

	auth.POST("/price", srv.AddPrice)
	auth.PUT("/price/:price_id", srv.UpdatePrice)
	auth.PUT("/users/:user_id/toggle-referral", srv.ToggleOpenToReferral)
	auth.PUT("/users/:user_id/banner", srv.UpdateUserBannerImage)

	//// profile/user
	auth.PUT("/users/:user_id/avatar", srv.UpdateUserAvatar)

	//// profile/socials
	auth.POST("/socials", srv.AddSocial)
	auth.PUT("/socials/:social_id", srv.UpdateSocial)
	auth.GET("/socials/:entity_type/:entity_id", srv.ListSocials)
	auth.DELETE("/socials/:social_id", srv.DeleteSocial)

	//// profile/experiences
	auth.POST("/users/:user_id/experiences", srv.AddExperience)
	auth.GET("/users/:user_id/experiences", srv.ListExperiences)
	auth.PUT("/users/:user_id/experiences/:experience_id", srv.UpdateExperience)
	auth.DELETE("/users/:user_id/experiences/:experience_id", srv.DeleteExperience)

	//// profile/educations
	auth.POST("/users/:user_id/educations", srv.AddEducation)
	auth.GET("/users/:user_id/educations", srv.ListEducations)
	auth.PUT("/users/:user_id/educations/:education_id", srv.UpdateEducation)
	auth.DELETE("/users/:user_id/educations/:education_id", srv.DeleteEducation)

	//// account
	auth.GET("/accounts/:user_id", srv.GetAccountInfo)

	//// network
	auth.POST("/connections/send", srv.SendConnection)
	auth.POST("/connections/:id/accept", srv.AcceptConnection)
	auth.POST("/connections/:id/reject", srv.RejectConnection)
	auth.GET("/connections/invitations/:user_id", srv.ListConnectionInvitations)
	auth.GET("/connections/:user_id", srv.ListConnections)
	auth.GET("/recommendations/:user_id", srv.ListRecommendations)
	auth.POST("/recommendations/cancel", srv.CancelRecommendation)
	auth.GET("/search/users", srv.SearchUsers)
	// check if user is connected to another user
	auth.GET("/connections/:user_id/:other_user_id", srv.CheckConnection)

	//// chat
	auth.GET("/chat/:room_id", func(ctx *gin.Context) {
		roomId := ctx.Param("room_id")
		chat.ServeWS(ctx, roomId, srv.Hub)
	})
	auth.GET("/chat/:room_id/messages", srv.ListMessages)
	auth.GET("/chat/users/:user_id/rooms", srv.ListChatRooms)
	auth.POST("/chat/rooms", srv.CreateChatRoom)

	//// referral
	auth.POST("/referral", srv.CreateReferral)
	auth.GET("/referrals/:user_id/active", srv.ListActiveReferrals)
	auth.GET("/referrals/users/:project_id", srv.ListReferredUsers)
	auth.GET("/users/:user_id/proposals", srv.ListActiveProposals)
	auth.POST("/proposals", srv.CreateProposal)
	auth.PUT("/proposals/:proposal_id", srv.UpdateProposal)
	auth.GET("users/:user_id/proposals/:project_id", srv.GetProposal)
	auth.POST("/projects/award", srv.AwardProject)
	auth.GET("/projects/awarded/:user_id", srv.ListAwardedProjects)
	auth.PUT("/projects/:project_id/accept", srv.AcceptProject)
	auth.PUT("/projects/:project_id/reject", srv.RejectProject)
	auth.GET("/projects/active/:user_id", srv.ListActiveProjects)
	auth.PUT("/projects/:project_id/start", srv.StartProject)
	auth.PUT("/projects/:project_id/initiate-complete", srv.InitiateCompleteProject)
	auth.PUT("/projects/:project_id/cancel/initiate-complete", srv.CancelInitiateCompleteProject)
	auth.PUT("/projects/:project_id/complete", srv.CompleteProject)
	auth.GET("/projects/completed/:user_id", srv.ListCompletedProjects)
	auth.POST("projects/review", srv.CreateProjectReview)
	auth.GET("projects/review/:project_id", srv.GetProjectReview)
	//
	auth.GET("/users/:user_id/connected", srv.ListConnectedUsers)
	auth.GET("/users", srv.ListUsers)
	auth.GET("/users/license-verified", srv.ListLicenseVerifiedUsers)
	auth.GET("/users/license-unverified", srv.ListLicenseUnverifiedUsers)

	//// approve license
	auth.PUT("/users/:user_id/approve-license", srv.ApproveLicense)
	auth.PUT("/users/:user_id/reject-license", srv.RejectLicense)

	//// posts
	auth.POST("/posts", srv.CreatePost)
	auth.GET("/posts/:post_id", srv.GetPost)
	auth.DELETE("/posts/:post_id", srv.DeletePost)
	auth.GET("/search/posts", srv.SearchPosts)
	auth.GET("/posts/:post_id/is-featured", srv.IsPostFeatured)

	//// news feed
	auth.GET("/feeds/:user_id", srv.ListNewsFeed)
	////auth.GET("/v2/feeds/:user_id", srv.listNewsFeedV2)
	////auth.GET("/v3/feeds/:user_id", srv.listNewsFeedV3)

	//// like post
	auth.POST("/posts/:post_id/like", srv.LikePost)
	auth.DELETE("/posts/:post_id/like", srv.UnlikePost)
	auth.GET("/posts/:post_id/liked-users", srv.ListPostLikedUsers)

	//// get post likes and comments count
	auth.GET("/posts/:post_id/likes-comments-count", srv.PostLikesAndCommentsCount)
	auth.GET("/posts/:post_id/is-liked", srv.IsPostLiked)

	//// comments
	auth.POST("/posts/:post_id/comments", srv.CommentPost)
	auth.GET("/posts/:post_id/comments", srv.ListComments)
	auth.POST("/comments/:comment_id/like", srv.LikeComment)
	auth.DELETE("/comments/:comment_id/like", srv.UnlikeComment)

	//// discussion
	auth.POST("/discussions", srv.CreateDiscussion)

	//// update discussion topic
	auth.PUT("/discussions/:discussion_id/topic", srv.UpdateDiscussionTopic)
	auth.POST("/discussions/:discussion_id/invite", srv.InviteUserToDiscussion)
	auth.POST("/discussions/:discussion_id/join", srv.JoinDiscussion)
	auth.POST("/discussions/:discussion_id/reject", srv.RejectDiscussion)
	auth.GET("/discussions/invites/:user_id", srv.ListDiscussionInvites)
	auth.GET("/discussions/active/:user_id", srv.ListActiveDiscussions)
	auth.GET("/discussions/:discussion_id/participants", srv.ListDiscussionParticipants)
	auth.GET("/discussions/:discussion_id/uninvited", srv.ListUninvitedParticipants)

	//// discussion messages
	auth.POST("/discussions/:discussion_id/messages", srv.SendMessageToDiscussion)
	auth.GET("/discussions/:discussion_id/messages", srv.ListDiscussionMessages)

	//// ads
	auth.POST("/ads", srv.CreateAd)
	//playing ads
	auth.GET("/ads/playing", srv.ListPlayingAds)
	auth.GET("/ads/expired", srv.ListExpiredAds)
	// extend ad period
	auth.PUT("/ads/:ad_id/extend", srv.ExtendAdPeriod)

	//// admin
	auth.GET("/attorneys", srv.ListAttorneys)
	auth.GET("/lawyers", srv.ListLawyers)

	auth.GET("/referrals/:user_id", srv.ListAllReferralProjects)
	auth.GET("/referrals/completed/:user_id", srv.ListCompletedReferralProjects)
	auth.GET("/referrals/active/:user_id", srv.ListActiveReferralProjects)

	auth.POST("/faqs", srv.CreateFAQ)
	auth.GET("/faqs", srv.ListFAQs)

	auth.POST("/firms", srv.AddFirm)
	auth.GET("/firms/owner/:owner_user_id", srv.ListFirmsByOwner)

	//// save post
	auth.POST("/saved-posts", srv.SavePost)
	auth.DELETE("/saved-posts/:saved_post_id", srv.UnSavePost)
	auth.GET("/saved-posts/:user_id", srv.ListSavedPosts)

	//// feature posts
	auth.POST("/feature-posts", srv.FeaturePost)
	auth.DELETE("/feature-posts/:post_id", srv.UnFeaturePost)
	auth.GET("/feature-posts/:user_id", srv.ListFeaturePosts)

	//// notifications
	auth.POST("/device-details", srv.SaveDevice)
	auth.POST("/notifications", srv.CreateNotification)
	auth.GET("/notifications/:user_id", srv.ListNotifications)

	//// post stats
	auth.GET("/posts/:post_id/stats", srv.GetPostStats)

	//// report
	auth.POST("/report-post", srv.ReportPost)
	auth.GET("/posts/:post_id/reported-status/:user_id", srv.IsPostReported)
	auth.DELETE("/feeds/:feed_id/ignore", srv.IgnoreFeed)

	//// activity
	auth.GET("/activity/posts/:user_id", srv.ListActivityPosts)
	auth.GET("/activity/comments/:user_id", srv.ListActivityComments)
	auth.GET("/users/:user_id/followers-count", srv.GetUserFollowersCount)

	// to run local
	err = r.Run(config.ServerAddress)
	log.Info().Err(err).Msg("cannot create srv:")

	// to run on lambda
	//ginLambda = ginadapter.New(r)
	//lambda.Start(Handler)
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
