package api

import (
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/imrishuroy/legal-referral/chat"
	"github.com/rs/zerolog/log"
)

//func playgroundHandler() gin.HandlerFunc {
//	h := playground.Handler("GraphQL", "/api/query")
//	return func(c *gin.Context) {
//		h.ServeHTTP(c.Writer, c.Request)
//	}
//}

func (s *Server) setupRouter(ginLambda *ginadapter.GinLambda) {
	// Set Gin to release mode
	//gin.SetMode(gin.ReleaseMode)
	log.Info().Msg("Setting up router 1")
	gin.SetMode(gin.ReleaseMode)

	//srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
	//	Store: server.Store}}))

	//http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	//http.Handle("/query", srv)

	s.Router = gin.Default()
	log.Info().Msg("Setting up routes 2")

	//server.router.GET("/playground", playgroundHandler())

	s.Router.GET("/", s.ping).Use(CORSMiddleware())
	s.Router.GET("/health", s.ping).Use(CORSMiddleware())
	s.Router.GET("/check", s.ping).Use(CORSMiddleware())

	// auth
	s.Router.POST("/api/sign-in", s.SignIn)
	s.Router.POST("/api/sign-up", s.SignUp)
	s.Router.POST("/api/refresh-token", s.RefreshToken)
	s.Router.POST("/api/otp/send", s.SendOTP)
	s.Router.POST("/api/otp/verify", s.VerifyOTP)
	s.Router.POST("/api/reset-password", s.ResetPassword)

	s.Router.GET("/api/users/:user_id/wizardstep", s.GetUserWizardStep)
	s.Router.GET("/api/firms", s.SearchFirms)

	s.Router.POST("/api/sign-in/linkedin", s.LinkedinLogin)

	auth := s.Router.Group("/api").
		Use(s.AuthMiddleware(s.FirebaseAuth))

	// GRAPHQL
	//auth.POST("/query", gin.WrapH(srv))

	auth.GET("/check-token", s.ping)
	auth.POST("/users", s.CreateUser)

	auth.GET("/users/:user_id", s.GetUserById)
	auth.POST("/license", s.SaveLicense)
	auth.POST("/license/upload", s.uploadLicense)
	auth.POST("/about-you", s.saveAboutYou)
	auth.GET("/users/:user_id/profile", s.fetchUserProfile)

	auth.PUT("/users/info", s.updateUserInfo)
	auth.POST("/review", s.addReview)

	auth.POST("/price", s.addPrice)
	auth.PUT("/price/:price_id", s.updatePrice)
	auth.PUT("/users/:user_id/toggle-referral", s.toggleOpenToReferral)
	auth.PUT("/users/:user_id/banner", s.updateUserBannerImage)

	// profile/user
	auth.PUT("/users/:user_id/avatar", s.updateUserAvatar)

	// profile/socials
	auth.POST("/socials", s.addSocial)
	auth.PUT("/socials/:social_id", s.updateSocial)
	auth.GET("/socials/:entity_type/:entity_id", s.listSocials)
	auth.DELETE("/socials/:social_id", s.deleteSocial)

	// profile/experiences
	auth.POST("/users/:user_id/experiences", s.addExperience)
	auth.GET("/users/:user_id/experiences", s.listExperiences)
	auth.PUT("/users/:user_id/experiences/:experience_id", s.updateExperience)
	auth.DELETE("/users/:user_id/experiences/:experience_id", s.deleteExperience)

	// profile/educations
	auth.POST("/users/:user_id/educations", s.addEducation)
	auth.GET("/users/:user_id/educations", s.listEducations)
	auth.PUT("/users/:user_id/educations/:education_id", s.updateEducation)
	auth.DELETE("/users/:user_id/educations/:education_id", s.deleteEducation)

	// account
	auth.GET("/accounts/:user_id", s.getAccountInfo)

	// network
	auth.POST("/connections/send", s.sendConnection)
	auth.POST("/connections/:id/accept", s.acceptConnection)
	auth.POST("/connections/:id/reject", s.rejectConnection)
	auth.GET("/connections/invitations/:user_id", s.listConnectionInvitations)
	auth.GET("/connections/:user_id", s.listConnections)
	auth.GET("/recommendations/:user_id", s.listRecommendations)
	auth.POST("/recommendations/cancel", s.cancelRecommendation)
	auth.GET("/search/users", s.searchUsers)
	// check if user is connected to another user
	auth.GET("/connections/:user_id/:other_user_id", s.checkConnection)

	// chat
	auth.GET("/chat/:room_id", func(ctx *gin.Context) {
		roomId := ctx.Param("room_id")
		chat.ServeWS(ctx, roomId, s.hub)
	})
	auth.GET("/chat/:room_id/messages", s.listMessages)
	auth.GET("/chat/users/:user_id/rooms", s.listChatRooms)
	auth.POST("/chat/rooms", s.createChatRoom)

	// referral
	auth.POST("/referral", s.createReferral)
	auth.GET("/referrals/:user_id/active", s.listActiveReferrals)
	auth.GET("/referrals/users/:project_id", s.listReferredUsers)
	auth.GET("/users/:user_id/proposals", s.listActiveProposals)
	auth.POST("/proposals", s.createProposal)
	auth.PUT("/proposals/:proposal_id", s.updateProposal)
	auth.GET("users/:user_id/proposals/:project_id", s.getProposal)
	auth.POST("/projects/award", s.awardProject)
	auth.GET("/projects/awarded/:user_id", s.listAwardedProjects)
	auth.PUT("/projects/:project_id/accept", s.acceptProject)
	auth.PUT("/projects/:project_id/reject", s.rejectProject)
	auth.GET("/projects/active/:user_id", s.listActiveProjects)
	auth.PUT("/projects/:project_id/start", s.startProject)
	auth.PUT("/projects/:project_id/initiate-complete", s.initiateCompleteProject)
	auth.PUT("/projects/:project_id/cancel/initiate-complete", s.cancelInitiateCompleteProject)
	auth.PUT("/projects/:project_id/complete", s.completeProject)
	auth.GET("/projects/completed/:user_id", s.listCompletedProjects)

	auth.POST("projects/review", s.CreateProjectReview)
	auth.GET("projects/review/:project_id", s.getProjectReview)

	auth.GET("/users/:user_id/connected", s.listConnectedUsers)
	auth.GET("/users", s.ListUsers)

	auth.GET("/users/license-verified", s.listLicenseVerifiedUsers)
	auth.GET("/users/license-unverified", s.listLicenseUnverifiedUsers)

	// approve license
	auth.PUT("/users/:user_id/approve-license", s.approveLicense)
	auth.PUT("/users/:user_id/reject-license", s.rejectLicense)

	// posts
	auth.POST("/posts", s.createPost)
	auth.GET("/posts/:post_id", s.getPost)
	auth.DELETE("/posts/:post_id", s.deletePost)
	auth.GET("/search/posts", s.searchPosts)
	auth.GET("/posts/:post_id/is-featured", s.isPostFeatured)

	// news feed
	auth.GET("/feeds/:user_id", s.listNewsFeed)
	//auth.GET("/v2/feeds/:user_id", server.listNewsFeedV2)
	//auth.GET("/v3/feeds/:user_id", server.listNewsFeedV3)

	// like post
	auth.POST("/posts/:post_id/like", s.likePost)
	auth.DELETE("/posts/:post_id/like", s.unlikePost)
	auth.GET("/posts/:post_id/liked-users", s.listPostLikedUsers)

	// get post likes and comments count
	auth.GET("/posts/:post_id/likes-comments-count", s.postLikesAndCommentsCount)
	auth.GET("/posts/:post_id/is-liked", s.isPostLiked)

	// comments
	auth.POST("/posts/:post_id/comments", s.commentPost)
	auth.GET("/posts/:post_id/comments", s.listComments)
	auth.POST("/comments/:comment_id/like", s.likeComment)
	auth.DELETE("/comments/:comment_id/like", s.unlikeComment)

	// discussion
	auth.POST("/discussions", s.createDiscussion)

	// update discussion topic
	auth.PUT("/discussions/:discussion_id/topic", s.updateDiscussionTopic)
	auth.POST("/discussions/:discussion_id/invite", s.inviteUserToDiscussion)
	auth.POST("/discussions/:discussion_id/join", s.joinDiscussion)
	auth.POST("/discussions/:discussion_id/reject", s.rejectDiscussion)
	auth.GET("/discussions/invites/:user_id", s.listDiscussionInvites)
	auth.GET("/discussions/active/:user_id", s.listActiveDiscussions)
	auth.GET("/discussions/:discussion_id/participants", s.listDiscussionParticipants)
	auth.GET("/discussions/:discussion_id/uninvited", s.listUninvitedParticipants)

	// discussion messages
	auth.POST("/discussions/:discussion_id/messages", s.sendMessageToDiscussion)
	auth.GET("/discussions/:discussion_id/messages", s.listDiscussionMessages)

	// ads
	auth.POST("/ads", s.createAd)
	//playing ads
	auth.GET("/ads/playing", s.listPlayingAds)
	auth.GET("/ads/expired", s.listExpiredAds)
	// extend ad period
	auth.PUT("/ads/:ad_id/extend", s.extendAdPeriod)

	// admin
	auth.GET("/attorneys", s.listAttorneys)
	auth.GET("/lawyers", s.listLawyers)

	auth.GET("/referrals/:user_id", s.listAllReferralProjects)
	auth.GET("/referrals/completed/:user_id", s.listCompletedReferralProjects)
	auth.GET("/referrals/active/:user_id", s.listActiveReferralProjects)

	auth.POST("/faqs", s.createFAQ)
	auth.GET("/faqs", s.listFAQs)

	auth.POST("/firms", s.addFirm)
	auth.GET("/firms/owner/:owner_user_id", s.listFirmsByOwner)

	// save post
	auth.POST("/saved-posts", s.savePost)
	auth.DELETE("/saved-posts/:saved_post_id", s.unSavePost)
	auth.GET("/saved-posts/:user_id", s.listSavedPosts)

	// feature posts
	auth.POST("/feature-posts", s.featurePost)
	auth.DELETE("/feature-posts/:post_id", s.unFeaturePost)
	auth.GET("/feature-posts/:user_id", s.listFeaturePosts)

	// notifications
	auth.POST("/device-details", s.saveDevice)
	auth.POST("/notifications", s.createNotification)
	auth.GET("/notifications/:user_id", s.listNotifications)

	// post stats
	auth.GET("/posts/:post_id/stats", s.getPostStats)

	// report
	auth.POST("/report-post", s.reportPost)
	auth.GET("/posts/:post_id/reported-status/:user_id", s.isPostReported)
	auth.DELETE("/feeds/:feed_id/ignore", s.ignoreFeed)

	// activity
	auth.GET("/activity/posts/:user_id", s.listActivityPosts)
	auth.GET("/activity/comments/:user_id", s.listActivityComments)
	auth.GET("/users/:user_id/followers-count", s.getUserFollowersCount)

	// print ginAdapter
	log.Info().Msg("Setting up routes 3")

	log.Info().Msgf("GinAdapter: %+v", ginLambda)

	ginLambda = ginadapter.New(s.Router)
}

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
