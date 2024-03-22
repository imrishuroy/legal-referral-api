package api

import (
	"github.com/gin-gonic/gin"
)

func (server *Server) setupRouter() {
	server.router = gin.Default()
	server.router.GET("/", server.ping)
	server.router.POST("/users", server.createUser)

	auth := server.router.Group("/api/auth").
		Use(signupMiddleware(server.firebaseAuth))

	auth.POST("/sign-up", server.SignUp)
	//auth.POST("/sign-in", server.SignIn)

	authorizedRoutes := server.router.Group("/api").
		Use(authMiddleware(server.firebaseAuth, server.config))

	authorizedRoutes.POST("/auth/sign-in", server.SignIn)
	authorizedRoutes.GET("/check-token", server.ping)
	authorizedRoutes.GET("/check-scope", server.checkScope)

}
