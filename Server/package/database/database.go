package database

import (
	"context"
	"fmt"
	"os"

	"github.com/TuanNghia295/BE-FIT/config"
	"github.com/jackc/pgx/v5"
)

func ConnectDB(conf *config.Config) *pgx.Conn {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		conf.Db.Host,
		conf.Db.User,
		conf.Db.Password,
		conf.Db.DBName,
		conf.Db.Port,
		conf.Db.SSLMode,
		conf.Db.TimeZone,
	)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully connected to database")
	return conn
}
