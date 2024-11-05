package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// DB is the global database connection pool.
var DB *sql.DB

// Initialize sets up the database connection using environment variables.
func Initialize() {
	var err error

	// Construct the connection string from environment variables
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	// Open a connection to the database
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Verify the connection before proceeding
	if err = DB.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	log.Println("Successfully connected to the database.")
}

// Close safely closes the database connection.
func Close() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Fatalf("Error closing database: %v", err)
		}
	}
}

func CheckReceiptExists(receiptHash string) (bool, error) {
    var exists bool
    err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM receipts WHERE receipt_hash = $1)", receiptHash).Scan(&exists)
    if err != nil {
        return false, err
    }
    return exists, nil
}

func InsertReceipt(id string, points int, hash string) error {
    _, err := DB.Exec("INSERT INTO receipts (id, points, receipt_hash) VALUES ($1, $2, $3)", id, points, hash)
    return err
}