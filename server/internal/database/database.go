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

func Migrate(db *sql.DB) error {
	const create string = `
	CREATE TABLE IF NOT EXISTS user (
		id VARCHAR(36) PRIMARY KEY,
		username VARCHAR(50) UNIQUE,
		first_name VARCHAR(100),
		email VARCHAR(100),
		password VARCHAR(100),
		is_admin BOOL DEFAULT false,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	res, err := db.Exec(create)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	log.Println("Migration completed")
	return nil
}
