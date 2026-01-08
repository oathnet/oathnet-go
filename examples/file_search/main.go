// File search example - async search within victim files.
//
// This example demonstrates:
// - Creating file search jobs
// - Different search modes (literal, regex, wildcard)
// - Polling for job completion
// - Processing search results
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

	client, err := oathnet.NewClient(apiKey)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		os.Exit(1)
	}

	// Method 1: All-in-one search (create + wait)
	fmt.Println("=== Simple File Search ===")
	result, err := client.FileSearch.Search("password", &oathnet.FileSearchCreateOptions{
		SearchMode: "literal",
		MaxMatches: 10,
	}, 60*time.Second)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Status: %s\n", result.Data.Status)
		if result.Data.Matches != nil {
			fmt.Printf("Found %d matches:\n", len(result.Data.Matches))
			for i, match := range result.Data.Matches {
				if i >= 5 {
					break
				}
				fmt.Printf("  - Log: %s\n", match.LogID)
				fmt.Printf("    File: %s\n", match.FileName)
				if match.MatchText != "" {
					preview := match.MatchText
					if len(preview) > 100 {
						preview = preview[:100]
					}
					fmt.Printf("    Match: %s...\n", preview)
				}
			}
		}
	}

	// Method 2: Manual job creation and polling
	fmt.Println("\n=== Manual Job Management ===")

	job, err := client.FileSearch.Create("api[_-]?key", &oathnet.FileSearchCreateOptions{
		SearchMode:     "regex",
		MaxMatches:     5,
		IncludeMatches: true,
		ContextLines:   2,
	})
	if err != nil {
		fmt.Printf("Error creating job: %v\n", err)
	} else {
		fmt.Printf("Created job: %s\n", job.Data.JobID)
		fmt.Printf("Initial status: %s\n", job.Data.Status)

		fmt.Println("\nPolling for completion...")
		completed, err := client.FileSearch.WaitForCompletion(
			job.Data.JobID,
			2*time.Second,
			60*time.Second,
		)
		if err != nil {
			fmt.Printf("Error waiting: %v\n", err)
		} else {
			fmt.Printf("Final status: %s\n", completed.Data.Status)
			if completed.Data.Matches != nil {
				fmt.Printf("\nFound %d regex matches\n", len(completed.Data.Matches))
			}
		}
	}

	// Search with specific log IDs
	fmt.Println("\n=== Search Within Specific Logs ===")

	victims, err := client.Victims.Search("gmail", &oathnet.VictimsSearchOptions{
		PageSize: 3,
	})
	if err != nil {
		fmt.Printf("Error searching victims: %v\n", err)
	} else {
		var logIDs []string
		for _, v := range victims.Data.Items {
			if v.LogID != "" {
				logIDs = append(logIDs, v.LogID)
				if len(logIDs) >= 3 {
					break
				}
			}
		}

		if len(logIDs) > 0 {
			fmt.Printf("Searching within %d victim logs...\n", len(logIDs))
			targetedResult, err := client.FileSearch.Search("*@gmail.com", &oathnet.FileSearchCreateOptions{
				SearchMode: "wildcard",
				LogIDs:     logIDs,
				MaxMatches: 10,
			}, 60*time.Second)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Status: %s\n", targetedResult.Data.Status)
				if targetedResult.Data.Matches != nil {
					fmt.Printf("Found %d matches in specified logs\n", len(targetedResult.Data.Matches))
				}
			}
		} else {
			fmt.Println("No log IDs available for targeted search")
		}
	}

	// Search with file pattern filter
	fmt.Println("\n=== Search with File Pattern Filter ===")
	patternResult, err := client.FileSearch.Search("token", &oathnet.FileSearchCreateOptions{
		SearchMode:  "literal",
		FilePattern: "*.txt",
		MaxMatches:  5,
	}, 60*time.Second)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		matchCount := 0
		if patternResult.Data.Matches != nil {
			matchCount = len(patternResult.Data.Matches)
		}
		fmt.Printf("Matches in .txt files: %d\n", matchCount)
	}
}
