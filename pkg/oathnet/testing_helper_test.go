package oathnet

import (
	"os"
	"testing"
)

// getAPIKey returns the API key from environment variable.
func getAPIKey(t *testing.T) string {
	apiKey := os.Getenv("OATHNET_API_KEY")
	if apiKey == "" {
		t.Skip("OATHNET_API_KEY environment variable required for integration tests")
	}
	return apiKey
}

// createTestClient creates a client for testing.
func createTestClient(t *testing.T) *Client {
	apiKey := getAPIKey(t)
	client, err := NewClient(apiKey)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	return client
}

// Test data constants
const (
	TestDiscordID            = "300760994454437890"
	TestDiscordIDWithHistory = "1375046349392974005"
	TestDiscordIDWithRoblox  = "1205957884584656927"
	TestSteamID              = "1100001586a2b38"
	TestXboxGamertag         = "ethan"
	TestRobloxUsername       = "chris"
	TestHoleheEmail          = "ethan_lewis_196@hotmail.co.uk"
	TestIP                   = "174.235.65.156"
	TestDomain               = "limabean.co.za"
	TestBreachQuery          = "winterfox"
	TestStealerQuery         = "gmail.com"
	TestVictimsQuery         = "gmail"
)
