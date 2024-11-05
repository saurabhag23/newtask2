// This program generates a JWT token for a specified username using the utils package.
// It can be used to generate tokens for testing API endpoints that require authentication.

package main

import (
    "fmt"
    "log"
    "github.com/saurabhag23/receipt-processor/internal/utils" // Import utils package for JWT functions
)

func main() {
    // Generate a JWT token for the user "saurabh"
    // Replace "saurabh" with any username if you need a token for a different user
    token, err := utils.GenerateJWT("saurabh")
    if err != nil {
        // Log an error and terminate the program if token generation fails
        log.Fatal("Error generating token:", err)
    }

    // Print the generated token to the console for use in API requests
    fmt.Println("Generated JWT Token:", token)
}