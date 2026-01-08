// Export example - async export to JSONL or CSV.
//
// This example demonstrates:
// - Creating export jobs
// - Export formats (JSONL, CSV)
// - Export types (docs, victims)
// - Waiting for completion and downloading
//
// Run: OATHNET_API_KEY="your-key" go run main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
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

	// Create export job
	fmt.Println("=== Create and Download Export ===")

	job, err := client.Exports.Create("docs", &oathnet.ExportCreateOptions{
		Format: "jsonl",
		Limit:  100,
		Search: map[string]string{
			"query": "gmail.com",
		},
	})
	if err != nil {
		fmt.Printf("Error creating export: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created job: %s\n", job.Data.JobID)
	fmt.Printf("Initial status: %s\n", job.Data.Status)

	// Wait for completion
	fmt.Println("Waiting for completion...")
	result, err := client.Exports.WaitForCompletion(
		job.Data.JobID,
		time.Second,
		120*time.Second,
	)
	if err != nil {
		fmt.Printf("Error waiting: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Final status: %s\n", result.Data.Status)

	if result.Data.Status == "completed" {
		// Download to file
		tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("oathnet-export-%d.jsonl", time.Now().Unix()))
		err = client.Exports.Download(job.Data.JobID, tempPath)
		if err != nil {
			fmt.Printf("Error downloading: %v\n", err)
		} else {
			info, _ := os.Stat(tempPath)
			fmt.Printf("Downloaded to: %s (%d bytes)\n", tempPath, info.Size())

			// Preview first few lines
			data, _ := os.ReadFile(tempPath)
			preview := string(data)
			if len(preview) > 500 {
				preview = preview[:500]
			}
			fmt.Printf("Preview:\n%s...\n", preview)

			// Clean up
			os.Remove(tempPath)
		}
	}

	// CSV Export
	fmt.Println("\n=== CSV Export with Fields ===")

	csvJob, err := client.Exports.Create("docs", &oathnet.ExportCreateOptions{
		Format: "csv",
		Limit:  50,
		Fields: []string{"email", "password", "domain", "url"},
		Search: map[string]string{
			"query": "gmail.com",
		},
	})
	if err != nil {
		fmt.Printf("Error creating CSV export: %v\n", err)
	} else {
		fmt.Printf("Created CSV job: %s\n", csvJob.Data.JobID)

		csvResult, err := client.Exports.WaitForCompletion(
			csvJob.Data.JobID,
			time.Second,
			120*time.Second,
		)
		if err != nil {
			fmt.Printf("Error waiting: %v\n", err)
		} else {
			fmt.Printf("Status: %s\n", csvResult.Data.Status)

			if csvResult.Data.Status == "completed" {
				csvPath := filepath.Join(os.TempDir(), fmt.Sprintf("oathnet-export-%d.csv", time.Now().Unix()))
				err := client.Exports.Download(csvJob.Data.JobID, csvPath)
				if err != nil {
					fmt.Printf("Error downloading: %v\n", err)
				} else {
					csvData, _ := os.ReadFile(csvPath)
					preview := string(csvData)
					if len(preview) > 500 {
						preview = preview[:500]
					}
					fmt.Printf("CSV Preview:\n%s\n", preview)
					os.Remove(csvPath)
				}
			}
		}
	}

	// Victims Export
	fmt.Println("\n=== Export Victims Data ===")

	victimsJob, err := client.Exports.Create("victims", &oathnet.ExportCreateOptions{
		Format: "jsonl",
		Limit:  50,
		Search: map[string]string{
			"query": "gmail",
		},
	})
	if err != nil {
		fmt.Printf("Error creating victims export: %v\n", err)
	} else {
		fmt.Printf("Created victims export job: %s\n", victimsJob.Data.JobID)

		victimsResult, err := client.Exports.WaitForCompletion(
			victimsJob.Data.JobID,
			time.Second,
			120*time.Second,
		)
		if err != nil {
			fmt.Printf("Error waiting: %v\n", err)
		} else {
			fmt.Printf("Status: %s\n", victimsResult.Data.Status)

			if victimsResult.Data.Status == "completed" && victimsResult.Data.Result != nil {
				fmt.Printf("Export complete: %d records, %d bytes\n",
					victimsResult.Data.Result.Records,
					victimsResult.Data.Result.FileSize)
			}
		}
	}
}
