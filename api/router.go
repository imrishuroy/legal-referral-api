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
	//server.router.POST("/api/upload", server.uploadfile)

	auth := server.router.Group("/api").
		Use(authMiddleware(server.firebaseAuth))

	auth.GET("/check-token", server.ping)
	auth.POST("/users", server.createUser)
	//auth.POST("/users/:user_id/image", server.updateUserImage)
	auth.GET("/users/:user_id", server.getUserById)
	auth.POST("/license", server.saveLicense)
	auth.POST("/license/upload", server.uploadLicense)
	auth.POST("/about-you", server.saveAboutYou)
	auth.POST("/experience", server.addExperience)
	auth.POST("/education", server.addEducation)
	auth.POST("/review", server.addReview)
	auth.POST("/social", server.addSocial)
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
