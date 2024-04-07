package api

import (
	"github.com/gin-gonic/gin"
)

func (server *Server) setupRouter() {
	server.router = gin.Default()
	server.router.GET("/", server.ping).Use(CORSMiddleware())

	server.router.POST("/api/sign-up", server.signUp)
	server.router.POST("/api/otp/send", server.sendOTP)
	server.router.POST("/api/otp/verify", server.verifyOTP)
	server.router.GET("/api/users/:user_id/wizardstep", server.getUserWizardStep)
	server.router.POST("/api/custom-signup", server.customTokenSignUp)

	auth := server.router.Group("/api").
		Use(authMiddleware(server.firebaseAuth))

	auth.POST("/users/:user_id/profile-image", server.updateUserImage)
	auth.POST("/users", server.createUser)
	auth.POST("/sign-in", server.signIn)
	auth.GET("/check-token", server.ping)
	auth.POST("/license", server.saveLicense)
	auth.POST("/license/upload", server.uploadLicense)
	auth.POST("/about-you", server.saveAboutYou)
	auth.POST("/experience", server.saveExperience)
	auth.GET("/users/:user_id", server.getUserById)

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
