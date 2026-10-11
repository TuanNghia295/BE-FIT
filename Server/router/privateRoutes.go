package router

import "github.com/TuanNghia295/BE-FIT/internal/delivery/http/middleware"

func (r *Router) SetupPrivateRoutes() {
	// Private routes that require authentication before accessing
	// apply middleware.AuthMiddleware(r.tokenService) to all private routes

	r.GET("/me", middleware.AuthMiddleware(r.tokenService), r.userHandler.Me)
	r.POST("/logout", middleware.AuthMiddleware(r.tokenService), r.userHandler.Logout)
	r.POST("/refresh", middleware.AuthMiddleware(r.tokenService), r.userHandler.Refresh)
	// r.POST("/update-profile", middleware.AuthMiddleware(r.tokenService), r.userHandler.UpdateProfile)
}
