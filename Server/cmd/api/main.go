package main

import (
	"fmt"
	"log"

	"github.com/TuanNghia295/BE-FIT/config"
	http "github.com/TuanNghia295/BE-FIT/internal/delivery/http/handler"
	"github.com/TuanNghia295/BE-FIT/internal/delivery/http/middleware"
	"github.com/TuanNghia295/BE-FIT/internal/infrastructure/postgres"
	"github.com/TuanNghia295/BE-FIT/internal/infrastructure/postgres/auth"
	"github.com/TuanNghia295/BE-FIT/internal/usecase/user"
	"github.com/TuanNghia295/BE-FIT/package/database"
	"github.com/gin-gonic/gin"
)

func main() {
	appConfig := config.GetConfig()

	// Connect to the database
	db := database.ConnectDB(appConfig)

	userRepo := postgres.NewUserRepository(db)
	refreshRepo := postgres.NewRefreshTokenRepository(db)
	if len(appConfig.JWTSecret) < 32 {
		log.Fatal("JWT_SECRET must contain at least 32 bytes")
	}
	tokenService := auth.NewJWTService(appConfig.JWTSecret)

	registerUsecase := user.NewRegisterUsecase(userRepo)
	loginUsecase := user.NewLoginUseCase(userRepo, refreshRepo, tokenService)
	refreshUsecase := user.NewRefreshUseCase(refreshRepo, tokenService)

	userHandler := http.NewUserHandler(registerUsecase, loginUsecase, refreshUsecase, appConfig.CookieSecure)

	r := gin.Default()

	r.Use(middleware.ErrorHandler())

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Server is running"})
	})

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.POST("/refresh", userHandler.Refresh)

	// Start API server
	PORT := appConfig.Server.Port
	fmt.Printf("BE-FIT API running on: http://localhost:%d\n", PORT)
	if err := r.Run(fmt.Sprintf(":%d", PORT)); err != nil {
		log.Fatal(err)
	}

}
