package persistence

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// InitDBFunc is a type for the database initialization function
type InitDBFunc func() (*sql.DB, error)

// initDB is the current database initialization function
var initDB InitDBFunc = defaultInitDB

// SetInitDB allows overriding the database initialization function
func SetInitDB(f InitDBFunc) {
	initDB = f
}

// InitDB calls the currently set database initialization function
func InitDB() (*sql.DB, error) {
	return initDB()
}

// defaultInitDB is the default implementation of InitDB
func defaultInitDB() (*sql.DB, error) {
	// Read database configuration from environment variables
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	// Build the connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Connect to the database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v\n", err)
	}

	// Verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Unable to ping the database: %v\n", err)
	}

	log.Println("Successfully connected to the database.")
	return db, nil
}
