package main

import (
	"fmt"
	"log"

	"github.com/TuanNghia295/BE-FIT/config"
	"github.com/TuanNghia295/BE-FIT/package/database"
	"github.com/gin-gonic/gin"
)

func main() {
	appConfig := config.GetConfig()

	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Server is running"})
	})

	// Connect to the database
	database.ConnectDB(appConfig)

	// Start API server
	PORT := appConfig.Server.Port
	fmt.Printf("BE-FIT API running on: http://localhost:%d\n", PORT)
	if err := r.Run(fmt.Sprintf(":%d", PORT)); err != nil {
		log.Fatal(err)
	}

}
