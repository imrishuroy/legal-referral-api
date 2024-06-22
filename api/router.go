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

	server.router.POST("/api/otp/send", server.sendOTP)
	server.router.POST("/api/otp/verify", server.verifyOTP)

	server.router.POST("/api/reset-password", server.resetPassword)
	server.router.GET("/api/users/:user_id/wizardstep", server.getUserWizardStep)
	server.router.POST("/api/firm", server.addFirm)
	server.router.GET("/api/firms", server.listFirms)

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
