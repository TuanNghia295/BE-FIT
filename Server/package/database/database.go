package database

import (
	"fmt"
	"os"

	"github.com/TuanNghia295/BE-FIT/config"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

func ConnectDB(conf *config.Config) *gorm.DB {
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

	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully connected to database")
	return db
}
