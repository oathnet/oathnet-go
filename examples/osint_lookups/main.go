// OSINT lookups example - all OSINT methods.
//
// This example demonstrates:
// - Discord user info and username history
// - Discord to Roblox mapping
// - Steam and Xbox profile lookups
// - Roblox user info
// - Email analysis with Holehe
// - IP geolocation
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

	// Discord User Info
	fmt.Println("=== Discord User Info ===")
	discordResult, err := client.OSINT.DiscordUserinfo("300760994454437890")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Username: %s\n", discordResult.Data.Username)
		fmt.Printf("Global Name: %s\n", discordResult.Data.GlobalName)
		fmt.Printf("Avatar: %s\n", discordResult.Data.Avatar)
	}

	// Discord Username History
	fmt.Println("\n=== Discord Username History ===")
	historyResult, err := client.OSINT.DiscordUsernameHistory("1375046349392974005")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else if historyResult.Data != nil && len(historyResult.Data.History) > 0 {
		fmt.Printf("Found %d username changes:\n", len(historyResult.Data.History))
		for i, entry := range historyResult.Data.History {
			if i >= 5 {
				break
			}
			fmt.Printf("  - %s (changed: %s)\n", entry.Username, entry.ChangedAt)
		}
	} else {
		fmt.Println("No username history found")
	}

	// Discord to Roblox
	fmt.Println("\n=== Discord to Roblox ===")
	robloxResult, err := client.OSINT.DiscordToRoblox("1205957884584656927")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Discord ID: %s\n", robloxResult.Data.DiscordID)
		fmt.Printf("Roblox ID: %s\n", robloxResult.Data.RobloxID)
		fmt.Printf("Roblox Username: %s\n", robloxResult.Data.RobloxUsername)
	}

	// Steam Profile
	fmt.Println("\n=== Steam Profile ===")
	steamResult, err := client.OSINT.Steam("1100001586a2b38")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Username: %s\n", steamResult.Data.Username)
		fmt.Printf("Profile URL: %s\n", steamResult.Data.ProfileURL)
		if steamResult.Data.RealName != "" {
			fmt.Printf("Real Name: %s\n", steamResult.Data.RealName)
		}
	}

	// Xbox Profile
	fmt.Println("\n=== Xbox Profile ===")
	xboxResult, err := client.OSINT.Xbox("ethan")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Gamertag: %s\n", xboxResult.Data.Username)
		fmt.Printf("Gamerscore: %d\n", xboxResult.Data.Gamerscore)
	}

	// Roblox User Info
	fmt.Println("\n=== Roblox User Info ===")
	robloxUserResult, err := client.OSINT.RobloxUserinfo(oathnet.RobloxUserinfoOptions{
		Username: "chris",
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("User ID: %d\n", robloxUserResult.Data.UserID)
		fmt.Printf("Username: %s\n", robloxUserResult.Data.Username)
		fmt.Printf("Display Name: %s\n", robloxUserResult.Data.DisplayName)
	}

	// Holehe Email Check
	fmt.Println("\n=== Holehe Email Check ===")
	holeheResult, err := client.OSINT.Holehe("ethan_lewis_196@hotmail.co.uk")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Email: ethan_lewis_196@hotmail.co.uk")
		if holeheResult.Data != nil && len(holeheResult.Data.Domains) > 0 {
			fmt.Printf("Found on %d services:\n", len(holeheResult.Data.Domains))
			for i, domain := range holeheResult.Data.Domains {
				if i >= 10 {
					break
				}
				status := "not found"
				if domain.Exists {
					status = "exists"
				}
				fmt.Printf("  - %s: %s\n", domain.Domain, status)
			}
		} else {
			fmt.Println("No registered services found")
		}
	}

	// IP Info
	fmt.Println("\n=== IP Info ===")
	ipResult, err := client.OSINT.IPInfo("174.235.65.156")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("IP: 174.235.65.156")
		fmt.Printf("Country: %s\n", ipResult.Data.Country)
		fmt.Printf("City: %s\n", ipResult.Data.City)
		fmt.Printf("Region: %s\n", ipResult.Data.RegionName)
		fmt.Printf("ISP: %s\n", ipResult.Data.ISP)
		fmt.Printf("AS Name: %s\n", ipResult.Data.ASName)
	}

	// Subdomain Extraction
	fmt.Println("\n=== Subdomain Extraction ===")
	subdomainResult, err := client.OSINT.ExtractSubdomain("limabean.co.za", nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Domain: limabean.co.za")
		if subdomainResult.Data != nil {
			fmt.Printf("Subdomain count: %d\n", subdomainResult.Data.Count)
			if len(subdomainResult.Data.Subdomains) > 0 {
				fmt.Println("Subdomains found:")
				for i, sub := range subdomainResult.Data.Subdomains {
					if i >= 10 {
						break
					}
					fmt.Printf("  - %s\n", sub)
				}
			}
		}
	}
}
