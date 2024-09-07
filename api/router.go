package api

import (
	"github.com/gin-gonic/gin"
	"github.com/imrishuroy/legal-referral/chat"
)

func (server *Server) setupRouter() {
	// Set Gin to release mode
	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.ReleaseMode)

	server.router = gin.Default()
	server.router.GET("/", server.ping).Use(CORSMiddleware())
	server.router.GET("/health", server.ping).Use(CORSMiddleware())
	server.router.GET("/check", server.ping).Use(CORSMiddleware())

	server.router.POST("/api/otp/send", server.sendOTP)
	server.router.POST("/api/otp/verify", server.verifyOTP)

	server.router.POST("/api/reset-password", server.resetPassword)
	server.router.GET("/api/users/:user_id/wizardstep", server.getUserWizardStep)

	server.router.GET("/api/firms", server.searchFirms)

	server.router.POST("/api/sign-in/linkedin", server.linkedinLogin)

	auth := server.router.Group("/api").
		Use(authMiddleware(server.firebaseAuth))

	auth.GET("/check-token", server.ping)
	auth.POST("/users", server.createUser)

	auth.GET("/users/:user_id", server.getUserById)
	auth.POST("/license", server.saveLicense)
	auth.POST("/license/upload", server.uploadLicense)
	auth.POST("/about-you", server.saveAboutYou)
	auth.GET("/users/:user_id/profile", server.fetchUserProfile)

	auth.PUT("/users/info", server.updateUserInfo)
	auth.POST("/review", server.addReview)

	auth.POST("/price", server.addPrice)
	auth.PUT("/price/:price_id", server.updatePrice)
	auth.PUT("/users/:user_id/toggle-referral", server.toggleOpenToReferral)
	auth.PUT("/users/:user_id/banner", server.updateUserBannerImage)

	// profile/user
	auth.PUT("/users/:user_id/avatar", server.updateUserAvatar)

	// profile/socials
	auth.POST("/socials", server.addSocial)
	auth.PUT("/socials/:social_id", server.updateSocial)
	auth.GET("/socials/:entity_type/:entity_id", server.listSocials)
	auth.DELETE("/socials/:social_id", server.deleteSocial)

	// profile/experiences
	auth.POST("/users/:user_id/experiences", server.addExperience)
	auth.GET("/users/:user_id/experiences", server.listExperiences)
	auth.PUT("/users/:user_id/experiences/:experience_id", server.updateExperience)
	auth.DELETE("/users/:user_id/experiences/:experience_id", server.deleteExperience)

	// profile/educations
	auth.POST("/users/:user_id/educations", server.addEducation)
	auth.GET("/users/:user_id/educations", server.listEducations)
	auth.PUT("/users/:user_id/educations/:education_id", server.updateEducation)
	auth.DELETE("/users/:user_id/educations/:education_id", server.deleteEducation)

	// network
	auth.POST("/connections/send", server.sendConnection)
	auth.POST("/connections/:id/accept", server.acceptConnection)
	auth.POST("/connections/:id/reject", server.rejectConnection)
	auth.GET("/connections/invitations/:user_id", server.listConnectionInvitations)
	auth.GET("/connections/:user_id", server.listConnections)
	auth.GET("/recommendations/:user_id", server.listRecommendations)
	auth.POST("/recommendations/cancel", server.cancelRecommendation)
	auth.GET("/search/users", server.searchUsers)
	// check if user is connected to another user
	auth.GET("/connections/:user_id/:other_user_id", server.checkConnection)

	// chat
	auth.GET("/chat/:room_id", func(ctx *gin.Context) {
		roomId := ctx.Param("room_id")
		chat.ServeWS(ctx, roomId, server.hub)
	})
	auth.GET("/chat/:room_id/messages", server.listMessages)
	auth.GET("/chat/users/:user_id/rooms", server.listChatRooms)
	auth.POST("/chat/rooms", server.createChatRoom)

	// referral
	auth.POST("/referral", server.createReferral)
	auth.GET("/referrals/:user_id/active", server.listActiveReferrals)
	auth.GET("/referrals/users/:project_id", server.listReferredUsers)
	auth.GET("/users/:user_id/proposals", server.listActiveProposals)
	auth.POST("/proposals", server.createProposal)
	auth.PUT("/proposals/:proposal_id", server.updateProposal)
	auth.GET("users/:user_id/proposals/:project_id", server.getProposal)
	auth.POST("/projects/award", server.awardProject)
	auth.GET("/projects/awarded/:user_id", server.listAwardedProjects)
	auth.PUT("/projects/:project_id/accept", server.acceptProject)
	auth.PUT("/projects/:project_id/reject", server.rejectProject)
	auth.GET("/projects/active/:user_id", server.listActiveProjects)
	auth.PUT("/projects/:project_id/start", server.startProject)
	auth.PUT("/projects/:project_id/initiate-complete", server.initiateCompleteProject)
	auth.PUT("/projects/:project_id/cancel/initiate-complete", server.cancelInitiateCompleteProject)
	auth.PUT("/projects/:project_id/complete", server.completeProject)
	auth.GET("/projects/completed/:user_id", server.listCompletedProjects)

	auth.POST("projects/review", server.createProjectReview)
	auth.GET("projects/review/:project_id", server.getProjectReview)

	auth.GET("/users/:user_id/connected", server.listConnectedUsers)
	auth.GET("/users", server.listUsers)

	auth.GET("/users/license-verified", server.listLicenseVerifiedUsers)
	auth.GET("/users/license-unverified", server.listLicenseUnverifiedUsers)

	// approve license
	auth.PUT("/users/:user_id/approve-license", server.approveLicense)
	auth.PUT("/users/:user_id/reject-license", server.rejectLicense)

	// posts
	auth.POST("/posts", server.createPost)

	// news feed
	auth.GET("/feeds/:user_id", server.listNewsFeed)

	// like post
	auth.POST("/posts/:post_id/like", server.likePost)
	auth.DELETE("/posts/:post_id/like", server.unlikePost)

	auth.GET("/posts/:post_id/liked-users", server.listPostLikedUsers)

	// comments
	auth.POST("/posts/:post_id/comments", server.commentPost)
	auth.GET("/posts/:post_id/comments", server.listComments)
	auth.POST("/comments/:comment_id/like", server.likeComment)
	auth.DELETE("/comments/:comment_id/like", server.unlikeComment)

	// discussion
	auth.POST("/discussions", server.createDiscussion)

	// update discussion topic
	auth.PUT("/discussions/:discussion_id/topic", server.updateDiscussionTopic)
	auth.POST("/discussions/:discussion_id/invite", server.inviteUserToDiscussion)
	auth.POST("/discussions/:discussion_id/join", server.joinDiscussion)
	auth.POST("/discussions/:discussion_id/reject", server.rejectDiscussion)
	auth.GET("/discussions/invites/:user_id", server.listDiscussionInvites)
	auth.GET("/discussions/active/:user_id", server.listActiveDiscussions)
	auth.GET("/discussions/:discussion_id/participants", server.listDiscussionParticipants)
	auth.GET("/discussions/:discussion_id/uninvited", server.listUninvitedParticipants)

	// discussion messages
	auth.POST("/discussions/:discussion_id/messages", server.sendMessageToDiscussion)
	auth.GET("/discussions/:discussion_id/messages", server.listDiscussionMessages)

	// ads
	auth.POST("/ads", server.createAd)
	//playing ads
	auth.GET("/ads/playing", server.listPlayingAds)
	auth.GET("/ads/expired", server.listExpiredAds)
	// extend ad period
	auth.PUT("/ads/:ad_id/extend", server.extendAdPeriod)

	// admin
	auth.GET("/attorneys", server.listAttorneys)
	auth.GET("/lawyers", server.listLawyers)

	auth.GET("/referrals/:user_id", server.listAllReferralProjects)
	auth.GET("/referrals/completed/:user_id", server.listCompletedReferralProjects)
	auth.GET("/referrals/active/:user_id", server.listActiveReferralProjects)

	auth.POST("/faqs", server.createFAQ)
	auth.GET("/faqs", server.listFAQs)

	auth.POST("/firms", server.addFirm)
	auth.GET("/firms/owner/:owner_user_id", server.listFirmsByOwner)

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
