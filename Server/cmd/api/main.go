package main

import (
	"log"

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
}
