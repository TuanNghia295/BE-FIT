package router

import "github.com/gin-gonic/gin"

func (r *Router) SetupPublicRoutes() {
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Server is running"})
	})

	r.POST("/login", r.userHandler.Login)
	r.POST("/register", r.userHandler.Register)
}
