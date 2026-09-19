package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/TuanNghia295/BE-FIT/package/database"
	"github.com/joho/godotenv"
)

func main() {
	// config env file
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Erro loading .env")
	}
	// Connect to the database
	database.ConnectDB()

	port := os.Getenv("PORT")
	// Start API server
	fmt.Printf("BE-FIT API running on: http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
