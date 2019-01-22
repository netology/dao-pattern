package postgresql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

// NewConnection is connection constructor
func NewConnection() *sql.DB {

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SCHEMA"))

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	return db
}
