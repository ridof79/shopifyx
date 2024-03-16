package db

import (
	"database/sql"
	"fmt"
	"time"

	"shopifyx/config"

	_ "github.com/lib/pq"
)

var db *sql.DB

func GetDB() *sql.DB {
	return db
}

func InitDB(config config.Config) {
	var conn string
	if config.ENV != "production" {
		conn = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			config.DbUsername,
			config.DbPassword,
			config.DbHost,
			config.DbPort,
			config.DbName)
	} else {
		conn = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=verify-full&sslrootcert=ap-southeast-1-bundle.pem&timezone=UTC",
			config.DbUsername,
			config.DbPassword,
			config.DbHost,
			config.DbPort,
			config.DbName)
	}

	fmt.Println(conn)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err.Error())
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(25)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	db.Ping()
}

func CloseDB() {
	db.Close()
}
