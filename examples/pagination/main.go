// Pagination example - cursor-based pagination patterns.
//
// This example demonstrates:
// - Cursor-based pagination (breach search)
// - Cursor-based pagination (V2 APIs)
// - Efficient iteration patterns
// - Collecting all results
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

	// Cursor-based pagination (breach search)
	fmt.Println("=== Cursor-Based Pagination (Breach) ===")
	cursor := ""
	totalFetched := 0
	maxPages := 3

	for page := 1; page <= maxPages; page++ {
		opts := &oathnet.SearchOptions{}
		if cursor != "" {
			opts.Cursor = cursor
		}

		result, err := client.Search.Breach("gmail.com", opts)
		if err != nil {
			fmt.Printf("Error on page %d: %v\n", page, err)
			break
		}

		count := len(result.Data.Results)
		totalFetched += count

		fmt.Printf("Page %d: %d results (Total: %d)\n", page, count, result.Data.ResultsFound)

		cursor = result.Data.Cursor
		if cursor == "" || count == 0 {
			fmt.Println("Reached end of results")
			break
		}
	}

	fmt.Printf("Fetched %d records across %d pages\n\n", totalFetched, maxPages)

	// Cursor-based pagination (V2 Stealer)
	fmt.Println("=== Cursor-Based Pagination (V2 Stealer) ===")
	cursor = ""
	pageSize := 25
	totalFetched = 0
	pages := 0

	for pages < maxPages {
		result, err := client.Stealer.Search("gmail.com", &oathnet.StealerSearchOptions{
			PageSize: pageSize,
			Cursor:   cursor,
		})
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			break
		}

		count := len(result.Data.Items)
		totalFetched += count
		pages++

		total := 0
		if result.Data.Meta != nil {
			total = result.Data.Meta.Total
		}
		fmt.Printf("Page %d: %d results (Total: %d)\n", pages, count, total)

		cursor = result.Data.NextCursor
		if cursor == "" {
			fmt.Println("No more pages")
			break
		}
	}

	fmt.Printf("Fetched %d stealer records across %d pages\n\n", totalFetched, pages)

	// Cursor-based pagination (V2 Victims)
	fmt.Println("=== Cursor-Based Pagination (V2 Victims) ===")
	cursor = ""
	totalFetched = 0
	pages = 0

	for pages < maxPages {
		result, err := client.Victims.Search("gmail", &oathnet.VictimsSearchOptions{
			PageSize: 10,
			Cursor:   cursor,
		})
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			break
		}

		count := len(result.Data.Items)
		totalFetched += count
		pages++

		total := 0
		if result.Data.Meta != nil {
			total = result.Data.Meta.Total
		}
		fmt.Printf("Page %d: %d victims (Total: %d)\n", pages, count, total)

		cursor = result.Data.NextCursor
		if cursor == "" {
			fmt.Println("No more pages")
			break
		}
	}

	fmt.Printf("Fetched %d victim profiles across %d pages\n\n", totalFetched, pages)

	// Collect all results helper pattern
	fmt.Println("=== Collect All Results Pattern ===")

	allResults := collectStealerResults(client, "gmail.com", 50)
	fmt.Printf("Collected %d total results\n", len(allResults))

	// Show unique domains from collected results
	domains := make(map[string]bool)
	for _, r := range allResults {
		for _, domain := range r.Domain {
			domains[domain] = true
		}
	}
	fmt.Printf("Unique domains: %d\n", len(domains))
}

func collectStealerResults(client *oathnet.Client, query string, maxResults int) []oathnet.V2StealerResult {
	var results []oathnet.V2StealerResult
	cursor := ""

	for len(results) < maxResults {
		remaining := maxResults - len(results)
		pageSize := 25
		if remaining < pageSize {
			pageSize = remaining
		}

		response, err := client.Stealer.Search(query, &oathnet.StealerSearchOptions{
			PageSize: pageSize,
			Cursor:   cursor,
		})
		if err != nil {
			break
		}

		results = append(results, response.Data.Items...)
		cursor = response.Data.NextCursor

		if cursor == "" || len(response.Data.Items) == 0 {
			break
		}
	}

	if len(results) > maxResults {
		results = results[:maxResults]
	}

	return results
}
