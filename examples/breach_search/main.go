// Breach search example with advanced filters and pagination.
//
// This example demonstrates:
// - Breach search with filters
// - Database filtering by name
// - Cursor-based pagination
//
// Run: OATHNET_API_KEY="your-key" go run main.go
package main

import (
	"fmt"
	"os"

	"github.com/oathnet/oathnet-go/pkg/oathnet"
)

func main() {
	apiKey := os.Getenv("OATHNET_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: Set OATHNET_API_KEY environment variable")
		os.Exit(1)
	}

	client, err := oathnet.NewClient(apiKey)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		os.Exit(1)
	}

	// Basic breach search
	fmt.Println("=== Basic Breach Search ===")
	result, err := client.Search.Breach("winterfox", nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Total found: %d\n", result.Data.ResultsFound)
	fmt.Printf("Results: %d\n", len(result.Data.Results))

	// Search with database filter
	fmt.Println("\n=== Search with Database Filter ===")
	linkedinResult, err := client.Search.Breach("gmail.com", &oathnet.SearchOptions{
		DBNames: "linkedin",
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("LinkedIn results: %d\n", linkedinResult.Data.ResultsFound)
	}

	// Dynamic field display
	fmt.Println("\n=== Dynamic Field Display ===")
	dynamicResult, err := client.Search.Breach("winterfox", nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		for i, record := range dynamicResult.Data.Results {
			if i >= 3 {
				break
			}
			fmt.Printf("\n--- Record %d ---\n", i+1)
			if record.Email != "" {
				fmt.Printf("  Email: %s\n", record.Email)
			}
			if record.Password != "" {
				fmt.Printf("  Password: %s\n", record.Password)
			}
			if record.DBName != "" {
				fmt.Printf("  DBName: %s\n", record.DBName)
			}
			if record.PhoneNumber != "" {
				fmt.Printf("  Phone: %s\n", record.PhoneNumber)
			}
		}
	}

	// Pagination example with cursor
	fmt.Println("\n=== Pagination Example ===")
	cursor := ""
	totalFetched := 0
	maxPages := 3

	for page := 1; page <= maxPages; page++ {
		opts := &oathnet.SearchOptions{}
		if cursor != "" {
			opts.Cursor = cursor
		}

		pageResult, err := client.Search.Breach("gmail.com", opts)
		if err != nil {
			fmt.Printf("Error on page %d: %v\n", page, err)
			break
		}

		count := len(pageResult.Data.Results)
		totalFetched += count
		fmt.Printf("Page %d: %d results\n", page, count)

		cursor = pageResult.Data.Cursor
		if cursor == "" || count == 0 {
			break
		}
	}

	fmt.Printf("Total fetched: %d\n", totalFetched)
}
