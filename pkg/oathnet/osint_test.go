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

func TestOSINTService_OpenAPIResponseDecoding(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/service/ip-info":
			_, _ = w.Write([]byte(`{"success":true,"data":{"query":"8.8.8.8","as":"AS15169 Google LLC","asname":"GOOGLE","mobile":false,"partial":true,"fields_missing":["proxy"],"provider_summary":{"status":"partial","total":2,"completed":1,"failed":1},"provider_statuses":[{"provider":"primary","status":"error","duration_ms":50.25},{"provider":"fallback","status":"partial","duration_ms":20.5}],"provider_timing":{"sample_size":2,"average_ms":35.38,"p50_ms":20.5,"p95_ms":50.25,"p99_ms":50.25,"max_ms":50.25,"sweep_ms":71,"module_timeout_ms":10000,"overall_timeout_ms":30000},"_meta":{"service":{"id":"ip-info"}}}}`))
		case "/service/steam":
			_, _ = w.Write([]byte(`{"success":true,"data":{"username":"alice","id":"76561198000000000","avatar":"https://example.com/a.png","meta":{"raw_data":{"steamid":"76561198000000000","timecreated":123},"source":"steam_api"},"_meta":{"service":{"id":"steam"}}}}`))
		case "/service/xbox":
			_, _ = w.Write([]byte(`{"success":true,"data":{"username":"alice","id":"xuid-1","avatar":"https://example.com/x.png","partial":true,"warning":"fallback provider","provider_statuses":{"playerdb":"unknown","scraper":"found"},"provider_summary":{"status":"partial","total":2,"completed":1,"failed":1},"provider_timing":{"sample_size":2,"average_ms":40,"p95_ms":70,"p99_ms":70,"max_ms":70,"sweep_ms":72,"module_timeout_ms":12000},"meta":{"id":"xuid-1","meta":{"gamerscore":"1234","xboxonerep":"GoodPlayer"},"scraper_data":{"games_played":2,"game_history":[{"title":"Halo","scoreDetails":{"achieved":10,"total":20}}]}}}}`))
		case "/service/discord-userinfo":
			_, _ = w.Write([]byte(`{"success":true,"data":{"id":"123456789012345678","username":"alice","global_name":"Alice","avatar_url":"https://example.com/d.png","banner_url":null,"creation_date":"2017-04-09 22:36:48 UTC","badges":["staff"]}}`))
		case "/service/discord-username-history":
			_, _ = w.Write([]byte(`{"success":true,"data":{"success":true,"message":"Found.","history":[{"name":"alice","time":"2026-01-01T00:00:00Z"},{"name":["bob"],"time":["2026-01-02T00:00:00Z"]}],"lookups_left":9}}`))
		case "/service/discord-to-roblox":
			_, _ = w.Write([]byte(`{"success":true,"data":{"discord_id":"123456789012345678","roblox_id":"1","name":"builderman","displayName":"Builder","created":"2006-01-01T00:00:00Z","description":"profile","avatar":"https://example.com/r.png","badges":["Admin"],"groupCount":3,"cached":false,"disabled":true,"skipped":true,"results_found":0}}`))
		case "/service/roblox-userinfo":
			_, _ = w.Write([]byte(`{"success":true,"data":{"username":"builderman","Current Username":"builderman","Old Usernames":"None","Display Name":"Builder","user_id":"1","User ID":"1","Discord":"N/A","Join Date":"2006-01-01T00:00:00Z","Avatar URL":"https://example.com/r.png"}}`))
		case "/service/holehe":
			_, _ = w.Write([]byte(`{"success":true,"data":{"domains":["github.com","google.com"],"partial":true,"provider_summary":{"status":"partial","total":2,"completed":1,"failed":1},"provider_statuses":[{"provider":"alpha","status":"found","duration_ms":12.5},{"provider":"beta","status":"timeout","duration_ms":6001.2}],"provider_timing":{"sample_size":2,"average_ms":3006.85,"p50_ms":12.5,"p95_ms":6001.2,"p99_ms":6001.2,"max_ms":6001.2,"sweep_ms":6002,"module_timeout_ms":6000,"overall_timeout_ms":25000},"_meta":{"service":{"id":"holehe"}}}}`))
		case "/service/ghunt":
			_, _ = w.Write([]byte(`{"success":true,"data":{"status":"found","data":{"profile":{"Name":"Alice","Profile Picture":"https://example.com/g.png","Gaia ID":"gaia-1","Last Update":"2026-01-01"},"maps_reviews":"none","photos_url":"https://photos.example.com"},"_meta":{"service":{"id":"ghunt"}}},"errors":{"error":"","details":""}}`))
		case "/service/extract-subdomain":
			_, _ = w.Write([]byte(`{"success":true,"data":{"domain":"example.com","subdomains":["www.example.com",{"subdomain":"api.example.com","alive":true}],"count":2}}`))
		case "/service/mc-history":
			_, _ = w.Write([]byte(`{"success":true,"data":{"uuid":"uuid-1","username":"Notch","history":[{"username":"Notch","changed_at":"2026-01-01T00:00:00Z"}]}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ip, err := client.OSINT.IPInfo("8.8.8.8")
	if err != nil {
		t.Fatalf("IPInfo() error = %v", err)
	}
	if ip.Data == nil || ip.Data.AS != "AS15169 Google LLC" || ip.Data.ASName != "GOOGLE" || ip.Data.Meta.Service.ID != "ip-info" || !ip.Data.Partial || ip.Data.ProviderSummary == nil || ip.Data.ProviderSummary.Completed != 1 || ip.Data.ProviderTiming == nil || ip.Data.ProviderTiming.P99MS == nil || *ip.Data.ProviderTiming.P99MS != 50.25 {
		t.Fatalf("unexpected IPInfo response: %#v", ip.Data)
	}

	steam, err := client.OSINT.Steam("76561198000000000")
	if err != nil {
		t.Fatalf("Steam() error = %v", err)
	}
	if steam.Data == nil || steam.Data.ID == "" || steam.Data.Meta == nil || steam.Data.Meta.RawData.TimeCreated != 123 {
		t.Fatalf("unexpected Steam response: %#v", steam.Data)
	}

	xbox, err := client.OSINT.Xbox("alice")
	if err != nil {
		t.Fatalf("Xbox() error = %v", err)
	}
	if xbox.Data == nil || xbox.Data.Meta == nil || xbox.Data.Meta.ScraperData.GameHistory[0].ScoreDetails.Total != 20 || !xbox.Data.Partial || xbox.Data.ProviderStatuses["scraper"] != "found" || xbox.Data.ProviderSummary == nil || xbox.Data.ProviderSummary.Completed != 1 || xbox.Data.ProviderTiming == nil || xbox.Data.ProviderTiming.P99MS == nil || *xbox.Data.ProviderTiming.P99MS != 70 {
		t.Fatalf("unexpected Xbox response: %#v", xbox.Data)
	}

	discord, err := client.OSINT.DiscordUserinfo("123456789012345678")
	if err != nil {
		t.Fatalf("DiscordUserinfo() error = %v", err)
	}
	if discord.Data == nil || discord.Data.AvatarURL == "" || len(discord.Data.Badges) != 1 {
		t.Fatalf("unexpected Discord response: %#v", discord.Data)
	}

	history, err := client.OSINT.DiscordUsernameHistory("123456789012345678")
	if err != nil {
		t.Fatalf("DiscordUsernameHistory() error = %v", err)
	}
	if history.Data == nil || len(history.Data.History) != 2 || history.Data.History[0].Name[0] != "alice" || history.Data.History[1].Name[0] != "bob" || history.Data.LookupsLeft == nil {
		t.Fatalf("unexpected history response: %#v", history.Data)
	}

	linked, err := client.OSINT.DiscordToRoblox("123456789012345678")
	if err != nil {
		t.Fatalf("DiscordToRoblox() error = %v", err)
	}
	if linked.Data == nil || linked.Data.Name != "builderman" || linked.Data.GroupCount != 3 || !linked.Data.Disabled || !linked.Data.Skipped || linked.Data.ResultsFound != 0 || linked.Data.DiscordID == "" {
		t.Fatalf("unexpected DiscordToRoblox response: %#v", linked.Data)
	}

	roblox, err := client.OSINT.RobloxUserinfo(RobloxUserinfoOptions{Username: "builderman"})
	if err != nil {
		t.Fatalf("RobloxUserinfo() error = %v", err)
	}
	if roblox.Data == nil || roblox.Data.CurrentUsername != "builderman" || roblox.Data.UserIDAlt != "1" {
		t.Fatalf("unexpected Roblox response: %#v", roblox.Data)
	}

	holehe, err := client.OSINT.Holehe("person@example.com")
	if err != nil {
		t.Fatalf("Holehe() error = %v", err)
	}
	if holehe.Data == nil || len(holehe.Data.Domains) != 2 {
		t.Fatalf("unexpected Holehe response: %#v", holehe.Data)
	}
	if !holehe.Data.Partial || holehe.Data.ProviderTiming == nil || holehe.Data.ProviderTiming.P99MS == nil || *holehe.Data.ProviderTiming.P99MS != 6001.2 {
		t.Fatalf("unexpected Holehe provider timing: %#v", holehe.Data)
	}

	ghunt, err := client.OSINT.GHunt("person@example.com")
	if err != nil {
		t.Fatalf("GHunt() error = %v", err)
	}
	if ghunt.Data == nil || ghunt.Data.Status != "found" || ghunt.Data.Data.Profile.GaiaID != "gaia-1" {
		t.Fatalf("unexpected GHunt response: %#v", ghunt.Data)
	}

	subdomains, err := client.OSINT.ExtractSubdomain("example.com", nil)
	if err != nil {
		t.Fatalf("ExtractSubdomain() error = %v", err)
	}
	if subdomains.Data == nil || len(subdomains.Data.Subdomains) != 2 {
		t.Fatalf("unexpected subdomain response: %#v", subdomains.Data)
	}
	if _, ok := subdomains.Data.Subdomains[1].(map[string]interface{}); !ok {
		t.Fatalf("second subdomain = %#v, want object entry", subdomains.Data.Subdomains[1])
	}

	minecraft, err := client.OSINT.MinecraftHistory("Notch")
	if err != nil {
		t.Fatalf("MinecraftHistory() error = %v", err)
	}
	if minecraft.Data == nil || minecraft.Data.History[0].Username != "Notch" {
		t.Fatalf("unexpected Minecraft response: %#v", minecraft.Data)
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
