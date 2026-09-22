package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func ConnectDB() *pgx.Conn {
	// Context use to manage request life cycle or goroutine.Especially Cancellation, timeout/deadline and sent request scoped data between layer
	/*
		context.Background() use to create root context. It has no cancellation, timout or deadline
	*/
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully connected to database")
	return conn
}
