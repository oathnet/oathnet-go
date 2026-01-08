// Error handling example - exception patterns and best practices.
//
// This example demonstrates:
// - Different error types
// - Proper error handling patterns
// - Retry logic for transient errors
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
	apiKey := os.Getenv("OATHNET_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: Set OATHNET_API_KEY environment variable")
		os.Exit(1)
	}

	// Example 1: Basic error handling
	fmt.Println("=== Basic Error Handling ===")
	client, _ := oathnet.NewClient(apiKey)

	result, err := client.Search.Breach("example.com", nil)
	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("Search successful: %d results\n", result.Data.ResultsFound)
	}

	// Example 2: Invalid API key handling
	fmt.Println("\n=== Invalid API Key ===")
	badClient, _ := oathnet.NewClient("invalid_api_key")

	_, err = badClient.Search.Breach("test", nil)
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}

	// Example 3: Validation errors
	fmt.Println("\n=== Validation Errors ===")
	_, err = client.OSINT.DiscordUserinfo("123") // Too short
	if err != nil {
		handleError(err)
	}

	// Example 4: No results handling
	fmt.Println("\n=== No Results Handling ===")
	noResult, err := client.Search.Breach("xyznonexistent123456789abcdef", nil)
	if err != nil {
		// Check if it's a "no results" error
		if validErr, ok := err.(*oathnet.ValidationError); ok {
			fmt.Printf("No results found (API returned validation error): %v\n", validErr)
		} else {
			handleError(err)
		}
	} else if noResult.Data.ResultsFound == 0 {
		fmt.Println("No results found (API returned success with empty results)")
	}

	// Example 5: Retry logic for transient errors
	fmt.Println("\n=== Retry Logic ===")
	retryResult := searchWithRetry(client, "gmail.com", 3, time.Second)
	if retryResult != nil {
		fmt.Printf("Search succeeded: %d results\n", retryResult.Data.ResultsFound)
	}

	// Example 6: Error type detection
	fmt.Println("\n=== Error Type Detection ===")
	_, err = client.OSINT.DiscordUserinfo("invalid")
	if err != nil {
		fmt.Printf("Error type: %T\n", err)
		fmt.Printf("Error message: %v\n", err)
	}

	fmt.Println("\nError handling examples complete!")
}

func handleError(err error) {
	switch e := err.(type) {
	case *oathnet.AuthenticationError:
		fmt.Printf("Authentication failed - check your API key: %v\n", e)
	case *oathnet.ValidationError:
		fmt.Printf("Invalid input: %v\n", e)
	case *oathnet.NotFoundError:
		fmt.Printf("Resource not found: %v\n", e)
	case *oathnet.RateLimitError:
		fmt.Printf("Rate limited - slow down requests: %v\n", e)
	case *oathnet.OathNetError:
		fmt.Printf("API error: %v\n", e)
	default:
		fmt.Printf("Unknown error: %v\n", err)
	}
}

func searchWithRetry(client *oathnet.Client, query string, maxRetries int, backoff time.Duration) *oathnet.BreachSearchResponse {
	for attempt := 0; attempt < maxRetries; attempt++ {
		result, err := client.Search.Breach(query, nil)
		if err == nil {
			return result
		}

		switch err.(type) {
		case *oathnet.RateLimitError:
			if attempt < maxRetries-1 {
				waitTime := backoff * time.Duration(1<<attempt)
				fmt.Printf("Rate limited, waiting %v...\n", waitTime)
				time.Sleep(waitTime)
			}
		case *oathnet.OathNetError:
			if attempt < maxRetries-1 {
				waitTime := backoff * time.Duration(1<<attempt)
				fmt.Printf("Server error, retrying in %v...\n", waitTime)
				time.Sleep(waitTime)
			}
		default:
			// Non-retryable error
			fmt.Printf("Non-retryable error: %v\n", err)
			return nil
		}
	}

	fmt.Println("Max retries exceeded")
	return nil
}
