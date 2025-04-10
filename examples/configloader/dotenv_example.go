package main

import (
	"fmt"
	"github.com/intezya/pkglib/configloader"
	"log"
)

// Run with: go run main.go --env-file=.env
func main() {
	// Load environment variables from .env file
	configloader.LoadEnv()

	// Example 1: Get required environment variable
	apiKey, err := configloader.GetEnv("API_KEY")
	if err != nil {
		// Handle missing required variable
		log.Printf("Error: %v", err)
		// Use fallback instead
		apiKey = "default-api-key"
	}
	fmt.Println("API Key:", apiKey)

	// Example 2: Get environment variable with fallback
	dbHost := configloader.GetEnvOrFallback("DB_HOST", "localhost")
	fmt.Println("Database Host:", dbHost)

	// Example 3: Get required integer environment variable
	port := configloader.GetEnvIntOrFallback("PORT", 8080)
	fmt.Println("Port:", port)

	// Example 4: Get required environment variable (panic version)
	// Note: This will panic if DEBUG_MODE is not set
	// debugMode := configloader.GetEnvOrPanic("DEBUG_MODE")

	// Example 5: Get required integer environment variable (panic version)
	// Note: This will panic if WORKER_COUNT is not set or not an integer
	// workerCount := configloader.GetEnvIntOrPanic("WORKER_COUNT")

	// Start your application...
	fmt.Println("Application running...")
}
