package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "github.com/saurabhag23/receipt-processor/internal/database"
    "github.com/saurabhag23/receipt-processor/internal/handlers"
)

func main() {
    // Load environment variables from a .env file if present
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    // Initialize database connection
    database.Initialize()
    defer database.Close()

	database.InitializeMongoDB()
	defer database.DisconnectMongoDB()
    // Create a new router
    r := mux.NewRouter()

    // Define the HTTP route for processing receipts
    r.HandleFunc("/receipts/process", handlers.ProcessReceipt).Methods("POST")

    // Define the HTTP route for retrieving points for a specific receipt by ID
    r.HandleFunc("/receipts/{id}/points", handlers.GetPoints).Methods("GET")

    // Define server address and port, could be configured via environment variables
    addr := os.Getenv("SERVER_ADDRESS")
    port := os.Getenv("SERVER_PORT")
    if addr == "" {
        addr = "localhost" // Default address
    }
    if port == "" {
        port = "8080" // Default port
    }

    // Start the HTTP server
    log.Printf("Server starting on %s:%s...\n", addr, port)
    log.Fatal(http.ListenAndServe(addr+":"+port, r))
}
