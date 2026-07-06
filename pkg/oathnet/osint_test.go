package oathnet

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestOSINTService_RequestConstruction(t *testing.T) {
	searchID := "search-123"
	alive := false

	tests := []struct {
		name      string
		call      func(*Client) error
		wantPath  string
		wantQuery url.Values
	}{
		{
			name: "IPInfo includes search_id",
			call: func(client *Client) error {
				_, err := client.OSINT.IPInfo("8.8.8.8", OSINTOptions{SearchID: searchID})
				return err
			},
			wantPath: "/service/ip-info",
			wantQuery: url.Values{
				"ip":        []string{"8.8.8.8"},
				"search_id": []string{searchID},
			},
		},
		{
			name: "Steam omits empty search_id",
			call: func(client *Client) error {
				_, err := client.OSINT.Steam("steam-1", OSINTOptions{})
				return err
			},
			wantPath: "/service/steam",
			wantQuery: url.Values{
				"steam_id": []string{"steam-1"},
			},
		},
		{
			name: "Xbox includes search_id",
			call: func(client *Client) error {
				_, err := client.OSINT.Xbox("xbox-1", OSINTOptions{SearchID: searchID})
				return err
			},
			wantPath: "/service/xbox",
			wantQuery: url.Values{
				"xbl_id":    []string{"xbox-1"},
				"search_id": []string{searchID},
			},
		},
		{
			name: "DiscordUserinfo includes search_id",
			call: func(client *Client) error {
				_, err := client.OSINT.DiscordUserinfo("discord-1", OSINTOptions{SearchID: searchID})
				return err
			},
			wantPath: "/service/discord-userinfo",
			wantQuery: url.Values{
				"discord_id": []string{"discord-1"},
				"search_id":  []string{searchID},
			},
		},
		{
			name: "DiscordUsernameHistory includes search_id",
			call: func(client *Client) error {
				_, err := client.OSINT.DiscordUsernameHistory("discord-2", OSINTOptions{SearchID: searchID})
				return err
			},
			wantPath: "/service/discord-username-history",
			wantQuery: url.Values{
				"discord_id": []string{"discord-2"},
				"search_id":  []string{searchID},
			},
		},
		{
			name: "DiscordToRoblox includes search_id",
			call: func(client *Client) error {
				_, err := client.OSINT.DiscordToRoblox("discord-3", OSINTOptions{SearchID: searchID})
				return err
			},
			wantPath: "/service/discord-to-roblox",
			wantQuery: url.Values{
				"discord_id": []string{"discord-3"},
				"search_id":  []string{searchID},
			},
		},
		{
			name: "RobloxUserinfo includes search_id",
			call: func(client *Client) error {
				_, err := client.OSINT.RobloxUserinfo(RobloxUserinfoOptions{
					Username: "builderman",
					SearchID: searchID,
				})
				return err
			},
			wantPath: "/service/roblox-userinfo",
			wantQuery: url.Values{
				"username":  []string{"builderman"},
				"search_id": []string{searchID},
			},
		},
		{
			name: "Holehe includes search_id",
			call: func(client *Client) error {
				_, err := client.OSINT.Holehe("person@example.com", OSINTOptions{SearchID: searchID})
				return err
			},
			wantPath: "/service/holehe",
			wantQuery: url.Values{
				"email":     []string{"person@example.com"},
				"search_id": []string{searchID},
			},
		},
		{
			name: "GHunt includes search_id",
			call: func(client *Client) error {
				_, err := client.OSINT.GHunt("person@example.com", OSINTOptions{SearchID: searchID})
				return err
			},
			wantPath: "/service/ghunt",
			wantQuery: url.Values{
				"email":     []string{"person@example.com"},
				"search_id": []string{searchID},
			},
		},
		{
			name: "ExtractSubdomain includes alive flag and search_id",
			call: func(client *Client) error {
				_, err := client.OSINT.ExtractSubdomain("example.com", &alive, OSINTOptions{SearchID: searchID})
				return err
			},
			wantPath: "/service/extract-subdomain",
			wantQuery: url.Values{
				"domain":    []string{"example.com"},
				"is_alive":  []string{"false"},
				"search_id": []string{searchID},
			},
		},
		{
			name: "MinecraftHistory uses canonical path and search_id",
			call: func(client *Client) error {
				_, err := client.OSINT.MinecraftHistory("Notch", OSINTOptions{SearchID: searchID})
				return err
			},
			wantPath: "/service/mc-history",
			wantQuery: url.Values{
				"username":  []string{"Notch"},
				"search_id": []string{searchID},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath string
			var gotQuery url.Values

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
				}
				gotPath = r.URL.Path
				gotQuery = r.URL.Query()
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"success":true}`))
			}))
			defer server.Close()

			client, err := NewClient("test-key", WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			if err := tt.call(client); err != nil {
				t.Fatalf("OSINT call error = %v", err)
			}
			if gotPath != tt.wantPath {
				t.Fatalf("path = %s, want %s", gotPath, tt.wantPath)
			}
			if !reflect.DeepEqual(gotQuery, tt.wantQuery) {
				t.Fatalf("query = %v, want %v", gotQuery, tt.wantQuery)
			}
		})
	}
}

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
