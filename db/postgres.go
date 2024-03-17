package db

import (
	"database/sql"
	"fmt"
	"time"

	"shopifyx/config"

	_ "github.com/lib/pq"
)

var db *sql.DB
var c config.Config

func GetDB() *sql.DB {
	return InitDB(c)
}

func SetDBConfig(config config.Config) {
	c = config
}

func InitDB(config config.Config) *sql.DB {
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
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err.Error())
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(25)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)
	return db
}

func CloseDB() {
	db.Close()
}
