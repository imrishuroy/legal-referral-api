package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRouter() {
	server.router = gin.Default()
	server.router.GET("/", server.ping)
	server.router.POST("/users", server.createUser)

	authRoutes := server.router.Group("/api").Use(VerifyAccessToken(server.config))

	authRoutes.GET("/check-token", server.ping)
	authRoutes.GET("/check-scope", server.checkScope)

}
