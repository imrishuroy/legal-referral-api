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
	"net/http"
)

// TODO: CHANGE THE ROUTER LIKE THIS https://github.com/build-on-aws/golang-gin-app-on-aws-lambda/blob/main/function/main.go

var ginLambda *ginadapter.GinLambda

func init() {

	log.Info().Msg("Welcome to LegalReferral")

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config : " + err.Error())
	}

	// db connection
	connPool, err := pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		fmt.Println("cannot connect to db:", err)
	}
	defer connPool.Close()

	store := db.NewStore(connPool)

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
	_, err = api.NewServer(config, store, hub, producer, ginLambda)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}

	log.Info().Msg("Server created")

	//server.Router = gin.Default()

	//server.router.GET("/playground", playgroundHandler())

	//r.GET("/", server.ping).Use(CORSMiddleware())
	//server.Router.GET("/health", server.ping).Use(CORSMiddleware())
	//server.Router.GET("/check", server.ping).Use(CORSMiddleware())
	//
	//// auth
	//server.Router.POST("/api/sign-in", server.signIn)
	//server.Router.POST("/api/sign-up", server.signUp)
	//server.Router.POST("/api/refresh-token", server.refreshToken)
	//server.Router.POST("/api/otp/send", server.sendOTP)
	//server.Router.POST("/api/otp/verify", server.verifyOTP)
	//server.Router.POST("/api/reset-password", server.resetPassword)
	//
	//server.Router.GET("/api/users/:user_id/wizardstep", server.getUserWizardStep)
	//server.Router.GET("/api/firms", server.searchFirms)
	//
	//server.Router.POST("/api/sign-in/linkedin", server.linkedinLogin)
	//
	//auth := server.Router.Group("/api").
	//	Use(authMiddleware(server.firebaseAuth))
	//
	//// GRAPHQL
	////auth.POST("/query", gin.WrapH(srv))
	//
	//auth.GET("/check-token", server.ping)
	//auth.POST("/users", server.createUser)
	//
	//auth.GET("/users/:user_id", server.getUserById)
	//auth.POST("/license", server.saveLicense)
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
	//auth.GET("/users", server.listUsers)
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
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, request)
}

func main() {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		log.Info().Msg("ping fun invoked")
		c.String(http.StatusOK, "pong")
	})

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
