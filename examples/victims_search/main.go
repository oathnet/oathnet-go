// Victims search example with manifest and file access.
//
// This example demonstrates:
// - Victim profile search
// - Fetching victim manifest (file tree)
// - Accessing individual files
//
// Run: OATHNET_API_KEY="your-key" go run main.go
package main

import (
	"fmt"
	"os"
	"strings"

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

	// Search for victim profiles
	fmt.Println("=== Victims Search ===")
	result, err := client.Victims.Search("gmail", &oathnet.VictimsSearchOptions{
		PageSize: 5,
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

	if len(result.Data.Items) == 0 {
		fmt.Println("No results found")
		return
	}

	// Display victim profiles
	fmt.Println("\n=== Victim Profiles ===")
	for i, victim := range result.Data.Items {
		if i >= 3 {
			break
		}
		fmt.Printf("\n--- Profile %d ---\n", i+1)
		fmt.Printf("  Log ID: %s\n", victim.LogID)
		if len(victim.DeviceUsers) > 0 {
			fmt.Printf("  Device Users: %v\n", victim.DeviceUsers)
		}
		if len(victim.DeviceIPs) > 0 {
			fmt.Printf("  Device IPs: %v\n", victim.DeviceIPs)
		}
		if victim.TotalDocs > 0 {
			fmt.Printf("  Total Docs: %d\n", victim.TotalDocs)
		}
		if victim.PwnedAt != "" {
			fmt.Printf("  Pwned At: %s\n", victim.PwnedAt)
		}
	}

	// Get manifest for first victim with log_id
	var logID string
	for _, victim := range result.Data.Items {
		if victim.LogID != "" {
			logID = victim.LogID
			break
		}
	}

	if logID == "" {
		fmt.Println("\nNo log IDs available for manifest demo")
		return
	}

	fmt.Printf("\n=== Manifest for %s ===\n", logID)
	manifest, err := client.Victims.GetManifest(logID)
	if err != nil {
		fmt.Printf("Could not fetch manifest: %v\n", err)
		return
	}

	if manifest != nil && manifest.VictimTree != nil {
		fmt.Printf("Log Name: %s\n", manifest.LogName)
		fmt.Println("File tree:")
		printTree(manifest.VictimTree, 0)
	}
}

// printTree recursively prints the victim file tree
func printTree(node *oathnet.VictimManifestNode, depth int) {
	if node == nil {
		return
	}
	indent := strings.Repeat("  ", depth)
	typeIcon := ""
	if node.Type == "directory" {
		typeIcon = "[DIR]"
	} else {
		typeIcon = "[FILE]"
	}
	size := ""
	if node.SizeBytes > 0 {
		size = fmt.Sprintf(" (%d bytes)", node.SizeBytes)
	}
	fmt.Printf("%s%s %s%s\n", indent, typeIcon, node.Name, size)

	// Print children (limit depth to avoid too much output)
	if depth < 3 && len(node.Children) > 0 {
		for i, child := range node.Children {
			if i >= 5 {
				fmt.Printf("%s  ... and %d more\n", indent, len(node.Children)-5)
				break
			}
			printTree(&child, depth+1)
		}
	}
}
