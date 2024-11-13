package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func GetConnection() (*sql.DB, error) {
	DB_NAME := os.Getenv("DB_NAME")
	if DB_NAME == "" {
		log.Fatalln("Error, database name not set")
	}

	db, err := sql.Open("sqlite3", DB_NAME)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connect to SQLITE database")
	return db, nil
}
