package api

import (
	"github.com/gin-gonic/gin"
)

func (server *Server) setupRouter() {
	server.router = gin.Default()
	server.router.GET("/", server.ping).Use(CORSMiddleware())
	server.router.POST("/api/otp/email", server.sendEmailOTP)
	server.router.POST("/api/otp/email/verify", server.verifyEmailOTP)
	server.router.POST("/api/otp/mobile", server.sendMobileOTP)
	server.router.POST("/api/otp/mobile/verify", server.verifyMobileOTP)

	server.router.GET("/api/users/:user_id", server.getUser)
	server.router.POST("/api/users", server.createUser)
	server.router.POST("/api/sign-up", server.signUp)
	server.router.POST("/api/license", server.saveLicense)
	server.router.POST("/api/about-you", server.saveAboutYou)
	server.router.POST("/api/experience", server.saveExperience)
	server.router.GET("/api/users/:user_id/wizardstep", server.getUserWizardStep)

	auth := server.router.Group("/api").
		Use(authMiddleware(server.firebaseAuth))

	// auth.POST("/sign-up", server.signUp)
	auth.POST("/sign-in", server.signIn)
	auth.GET("/check-token", server.ping)

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
