package api

import (
	"github.com/gin-gonic/gin"
)

func (server *Server) setupRouter() {
	server.router = gin.Default()
	server.router.GET("/", server.ping).Use(CORSMiddleware())
	server.router.POST("/otp", server.sendOTP)
	server.router.POST("/otp/verify", server.verifyOTP)
	server.router.POST("/users", server.createUser)
	server.router.POST("/api/sign-up", server.signUp)

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
