package api

import (
	"github.com/gin-gonic/gin"
)

func (server *Server) setupRouter() {
	server.router = gin.Default()
	server.router.GET("/", server.ping).Use(CORSMiddleware())

	server.router.POST("/api/otp/send", server.sendOTP)
	server.router.POST("/api/otp/verify", server.verifyOTP)

	server.router.POST("/api/reset-password", server.resetPassword)
	server.router.GET("/api/users/:user_id/wizardstep", server.getUserWizardStep)
	server.router.POST("/api/firm", server.addFirm)
	server.router.GET("/api/firms", server.listFirms)

	auth := server.router.Group("/api").
		Use(authMiddleware(server.firebaseAuth))

	auth.GET("/check-token", server.ping)
	auth.POST("/users", server.createUser)
	//auth.POST("/users/:user_id/image", server.updateUserImage)
	auth.GET("/users/:user_id", server.getUserById)
	auth.POST("/license", server.saveLicense)
	auth.POST("/license/upload", server.uploadLicense)
	auth.POST("/about-you", server.saveAboutYou)
	auth.GET("/users/:user_id/profile", server.fetchUserProfile)

	auth.PUT("/users/info", server.updateUserInfo)
	auth.POST("/review", server.addReview)
	auth.POST("/socials", server.addSocial)
	auth.PUT("/socials/:social_id", server.updateSocial)
	auth.GET("/socials/:entity_type/:entity_id", server.listSocials)
	auth.POST("/price", server.addPrice)
	auth.PUT("/price/:price_id", server.updatePrice)
	auth.PUT("/users/:user_id/toggle-referral", server.toggleOpenToReferral)
	auth.PUT("/users/:user_id/banner", server.updateUserBannerImage)

	// profile/user
	auth.PUT("/users/:user_id/avatar", server.updateUserAvatar)

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
