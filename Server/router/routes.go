package router

import (
	httpHandler "github.com/TuanNghia295/BE-FIT/internal/delivery/http/handler"
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
	// Setup public routes that don't require authentication
	r.SetupPublicRoutes()

	// Setup private routes that require authentication
	r.SetupPrivateRoutes()
}
