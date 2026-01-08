// Stealer search example with V2 API.
//
// This example demonstrates:
// - V2 stealer search with various filters
// - Domain and subdomain filtering
// - Log ID access for victim profile linking
// - Subdomain extraction
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

	// Basic stealer search
	fmt.Println("=== V2 Stealer Search ===")
	result, err := client.Stealer.Search("gmail.com", &oathnet.StealerSearchOptions{
		PageSize: 10,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Success: %v\n", result.Success)
	total := 0
	if result.Data.Meta != nil {
		total = result.Data.Meta.Total
	}
	fmt.Printf("Total found: %d\n", total)
	fmt.Printf("Results: %d\n", len(result.Data.Items))

	// Show results with log IDs
	fmt.Println("\n=== Results with Log IDs ===")
	var logIDs []string
	for i, item := range result.Data.Items {
		if i >= 5 {
			break
		}
		fmt.Printf("\n--- Record %d ---\n", i+1)
		if item.URL != "" {
			fmt.Printf("  URL: %s\n", item.URL)
		}
		if item.Username != "" {
			fmt.Printf("  Username: %s\n", item.Username)
		}
		if item.Password != "" {
			fmt.Println("  Password: ********")
		}
		if item.LogID != "" {
			fmt.Printf("  Log ID: %s\n", item.LogID)
			logIDs = append(logIDs, item.LogID)
		}
	}

	// Unique log IDs
	if len(logIDs) > 0 {
		uniqueLogIDs := make(map[string]bool)
		for _, lid := range logIDs {
			uniqueLogIDs[lid] = true
		}
		fmt.Println("\n=== Unique Log IDs (for victim profile access) ===")
		fmt.Printf("Found %d unique victim profiles:\n", len(uniqueLogIDs))
		count := 0
		for lid := range uniqueLogIDs {
			if count >= 10 {
				break
			}
			fmt.Printf("  - %s\n", lid)
			count++
		}
	}

	// Domain-specific search
	fmt.Println("\n=== Domain-Specific Search ===")
	domainResult, err := client.Stealer.Search("", &oathnet.StealerSearchOptions{
		Domain:   "google.com",
		PageSize: 5,
		Wildcard: true,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		domainTotal := 0
		if domainResult.Data.Meta != nil {
			domainTotal = domainResult.Data.Meta.Total
		}
		fmt.Printf("Google domain results: %d\n", domainTotal)
	}

	// Search with has_log_id filter
	fmt.Println("\n=== Search with Log ID Filter ===")
	logIDResult, err := client.Stealer.Search("gmail.com", &oathnet.StealerSearchOptions{
		HasLogID: true,
		PageSize: 5,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		logTotal := 0
		if logIDResult.Data.Meta != nil {
			logTotal = logIDResult.Data.Meta.Total
		}
		fmt.Printf("Results with victim profiles: %d\n", logTotal)
	}

	// Subdomain extraction
	fmt.Println("\n=== Subdomain Extraction ===")
	subdomainResult, err := client.Stealer.Subdomain("google.com", "")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success: %v\n", subdomainResult.Success)
		if subdomainResult.Data != nil && len(subdomainResult.Data.Subdomains) > 0 {
			fmt.Printf("Found %d subdomains:\n", subdomainResult.Data.Count)
			for i, sub := range subdomainResult.Data.Subdomains {
				if i >= 10 {
					break
				}
				fmt.Printf("  - %s\n", sub)
			}
		}
	}
}
