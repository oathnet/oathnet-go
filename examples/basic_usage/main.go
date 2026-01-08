// Basic usage example for OathNet SDK.
//
// This example demonstrates:
// - Client initialization
// - Simple breach search
// - Basic error handling
//
// Run: OATHNET_API_KEY="your-key" go run main.go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/oathnet/oathnet-go/pkg/oathnet"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("OATHNET_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: Set OATHNET_API_KEY environment variable")
		os.Exit(1)
	}

	// Initialize client with longer timeout
	client, err := oathnet.NewClient(apiKey, oathnet.WithTimeout(60*time.Second))
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		os.Exit(1)
	}

	// Simple breach search
	fmt.Println("Searching for 'example.com' in breach database...")
	result, err := client.Search.Breach("example.com", nil)
	if err != nil {
		// Handle different error types
		switch e := err.(type) {
		case *oathnet.AuthenticationError:
			fmt.Printf("Authentication failed: %v\n", e)
		case *oathnet.ValidationError:
			fmt.Printf("Invalid input: %v\n", e)
		default:
			fmt.Printf("API error: %v\n", err)
		}
		os.Exit(1)
	}

	if result.Success {
		fmt.Printf("Found %d results\n", result.Data.ResultsFound)
		fmt.Printf("Retrieved %d records\n", len(result.Data.Results))

		// Print first few results
		for i, record := range result.Data.Results {
			if i >= 5 {
				break
			}
			fmt.Printf("\n--- Result %d ---\n", i+1)
			if record.Email != "" {
				fmt.Printf("Email: %s\n", record.Email)
			}
			if record.DBName != "" {
				fmt.Printf("Database: %s\n", record.DBName)
			}
			if record.Password != "" {
				fmt.Println("Password: ********")
			}
		}
	}
}
