package main

import (
	"fmt"
	"log"

	"github.com/TuanNghia295/BE-FIT/config"
	http "github.com/TuanNghia295/BE-FIT/internal/delivery/http/handler"
	"github.com/TuanNghia295/BE-FIT/internal/infrastructure/postgres"
	"github.com/TuanNghia295/BE-FIT/internal/infrastructure/postgres/auth"
	"github.com/TuanNghia295/BE-FIT/internal/usecase"
	"github.com/TuanNghia295/BE-FIT/package/database"
	"github.com/TuanNghia295/BE-FIT/router"
)

func main() {
	appConfig := config.GetConfig()

	// Connect to the database
	db := database.ConnectDB(appConfig)

	configRepo := postgres.ConfigRepo(db)
	tokenService, err := auth.LoadJWTService(appConfig.JWTPrivateKeyPath, appConfig.JWTPublicKeyPath)
	if err != nil {
		log.Fatalf("load JWT key pair: %v", err)
	}

	useCases := usecase.NewUseCase(
		configRepo.UserRepository,
		configRepo.RefreshTokenRepository,
		tokenService,
	)
	userHandler := http.NewUserHandler(
		useCases.RegisterUser,
		useCases.LoginUser,
		useCases.RefreshUser,
		useCases.MeUser,
		appConfig.CookieSecure,
	)

	r := router.NewRouter(userHandler, tokenService)
	r.SetupRoutes()

	// Start API server
	PORT := appConfig.Server.Port
	fmt.Printf("BE-FIT API running on: http://localhost:%d\n", PORT)
	if err := r.Run(fmt.Sprintf(":%d", PORT)); err != nil {
		log.Fatal(err)
	}

}
