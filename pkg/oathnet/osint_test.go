package oathnet

import (
	"testing"
)

func TestOSINTService_DiscordUserinfo(t *testing.T) {
	client := createTestClient(t)

	t.Run("get Discord user info", func(t *testing.T) {
		result, err := client.OSINT.DiscordUserinfo(TestDiscordID)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
		if result.Data == nil || result.Data.Username == "" {
			t.Error("Expected username to be set")
		}
	})
}

func TestOSINTService_DiscordUsernameHistory(t *testing.T) {
	client := createTestClient(t)

	t.Run("get Discord username history", func(t *testing.T) {
		result, err := client.OSINT.DiscordUsernameHistory(TestDiscordIDWithHistory)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
	})
}

func TestOSINTService_DiscordToRoblox(t *testing.T) {
	client := createTestClient(t)

	t.Run("get Discord to Roblox mapping", func(t *testing.T) {
		result, err := client.OSINT.DiscordToRoblox(TestDiscordIDWithRoblox)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
		if result.Data == nil || result.Data.RobloxID == "" {
			t.Error("Expected Roblox ID to be set")
		}
	})
}

func TestOSINTService_Steam(t *testing.T) {
	client := createTestClient(t)

	t.Run("get Steam profile", func(t *testing.T) {
		result, err := client.OSINT.Steam(TestSteamID)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
		if result.Data == nil || result.Data.Username == "" {
			t.Error("Expected username to be set")
		}
	})
}

func TestOSINTService_Xbox(t *testing.T) {
	client := createTestClient(t)

	t.Run("get Xbox profile", func(t *testing.T) {
		result, err := client.OSINT.Xbox(TestXboxGamertag)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
		if result.Data == nil || result.Data.Username == "" {
			t.Error("Expected username to be set")
		}
	})
}

func TestOSINTService_RobloxUserinfo(t *testing.T) {
	client := createTestClient(t)

	t.Run("get Roblox user by username", func(t *testing.T) {
		result, err := client.OSINT.RobloxUserinfo(RobloxUserinfoOptions{
			Username: TestRobloxUsername,
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
		if result.Data == nil || result.Data.UserID == "" {
			t.Error("Expected user ID to be set")
		}
	})
}

func TestOSINTService_Holehe(t *testing.T) {
	client := createTestClient(t)

	t.Run("check email registration", func(t *testing.T) {
		result, err := client.OSINT.Holehe(TestHoleheEmail)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
		if result.Data == nil {
			t.Error("Expected data to be non-nil")
		}
	})
}

func TestOSINTService_IPInfo(t *testing.T) {
	client := createTestClient(t)

	t.Run("get IP geolocation", func(t *testing.T) {
		result, err := client.OSINT.IPInfo(TestIP)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
		if result.Data == nil || result.Data.Country == "" {
			t.Error("Expected country to be set")
		}
		if result.Data.City == "" {
			t.Error("Expected city to be set")
		}
	})
}

func TestOSINTService_ExtractSubdomain(t *testing.T) {
	client := createTestClient(t)

	t.Run("extract subdomains", func(t *testing.T) {
		result, err := client.OSINT.ExtractSubdomain(TestDomain, nil)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
	})
}
