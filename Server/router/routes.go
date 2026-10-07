package router

import (
	httpHandler "github.com/TuanNghia295/BE-FIT/internal/delivery/http/handler"
	"github.com/TuanNghia295/BE-FIT/internal/delivery/http/middleware"
	"github.com/TuanNghia295/BE-FIT/internal/infrastructure/postgres/auth"
	"github.com/gin-gonic/gin"
)

type Router struct {
	*gin.Engine
	userHandler  *httpHandler.UserHandler
	tokenService *auth.JWTService
}

// Create a new router instance
func NewRouter(userHandler *httpHandler.UserHandler, tokenService *auth.JWTService) *Router {
	return &Router{
		Engine:       gin.Default(),
		userHandler:  userHandler,
		tokenService: tokenService,
	}
}

func (r *Router) SetupRoutes() {
	r.Use(middleware.ErrorHandler())

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Server is running"})
	})

	r.POST("/login", r.userHandler.Login)
	r.POST("/register", r.userHandler.Register)
	r.POST("/refresh", r.userHandler.Refresh)
	r.POST("/logout", r.userHandler.Logout)
	r.GET("/me", middleware.AuthMiddleware(r.tokenService), r.userHandler.Me)
}
