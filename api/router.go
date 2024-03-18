package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRouter() {
	server.router = gin.Default()
	server.router.GET("/", server.ping)
	server.router.POST("/users", server.createUser)

	server.router.POST("/auth/signup", server.SignUp)

	//idTokenRoutes := server.router.Group("/auth").Use(VerifyIDToken())
	//idTokenRoutes.GET("/check-id-token", server.ping)

	authRoutes := server.router.Group("/api").Use(VerifyAccessToken(server.config))

	// authRoutes.POST("/signup", server.SignUp)
	authRoutes.GET("/check-token", server.ping)
	authRoutes.GET("/check-scope", server.checkScope)

}
