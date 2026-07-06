package oathnet

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestStealerV2Service_Search(t *testing.T) {
	client := createTestClient(t)

	t.Run("basic V2 stealer search", func(t *testing.T) {
		result, err := client.Stealer.Search(TestStealerQuery, &StealerSearchOptions{
			PageSize: 5,
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
		if result.Data == nil {
			t.Error("Expected data to be non-nil")
		}
	})

	t.Run("search with domain filter", func(t *testing.T) {
		result, err := client.Stealer.Search("", &StealerSearchOptions{
			Domain:   "google.com",
			PageSize: 5,
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
	})

	t.Run("search with has_log_id filter", func(t *testing.T) {
		result, err := client.Stealer.Search(TestStealerQuery, &StealerSearchOptions{
			HasLogID: true,
			PageSize: 5,
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
	})

	t.Run("cursor pagination", func(t *testing.T) {
		t.Skip("Skipping: V2 stealer cursor pagination has known issues")
	})

	t.Run("wildcard search", func(t *testing.T) {
		result, err := client.Stealer.Search("gmail", &StealerSearchOptions{
			Wildcard: true,
			PageSize: 5,
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
	})
}

func TestStealerV2Service_SearchRequestConstruction(t *testing.T) {
	var gotPath string
	var gotQuery url.Values

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.EscapedPath()
		gotQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"success": true,
			"data": {
				"items": [{
					"id": "credential-1",
					"log_id": "log-123",
					"url_str": "https://accounts.example.com/login",
					"domain": ["example.com"],
					"username": "alice",
					"password": "secret"
				}],
				"meta": {"total": 1, "count": 1, "took_ms": 5},
				"next_cursor": "next-stealer"
			}
		}`))
	}))
	defer server.Close()

	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	filter := map[string]interface{}{
		"and": []interface{}{
			map[string]interface{}{"field": "domain", "operator": "eq", "value": "example.com"},
		},
	}

	resp, err := client.Stealer.Search("alice@example.com", &StealerSearchOptions{
		Cursor:                 "cursor/with space",
		PageSize:               50,
		Sort:                   "-indexed_at",
		From:                   "2026-01-01T00:00:00Z",
		To:                     "2026-01-31T00:00:00Z",
		DateField:              "indexed_at",
		LogID:                  "log-123",
		HasLogID:               true,
		Wildcard:               true,
		Logic:                  "and",
		Filter:                 filter,
		FilterID:               "0123456789abcdef01234567",
		Domains:                []string{"example.com"},
		Subdomains:             []string{"accounts.example.com"},
		Usernames:              []string{"alice"},
		Passwords:              []string{"secret"},
		PasswordHashes:         []string{"sha256:value"},
		Paths:                  []string{"/Users/Alice/AppData"},
		Emails:                 []string{"alice@example.com"},
		EmailDomains:           []string{"example.com"},
		IPs:                    []string{"203.0.113.10"},
		HWIDs:                  []string{"hwid-1"},
		DiscordIDs:             []string{"1234567890"},
		SourceTypes:            []string{"stealer"},
		ArchiveHashes:          []string{"archive-hash"},
		CanonicalCredentialIDs: []string{"cred-1"},
		Fields:                 []string{"url_str", "domain"},
		SearchID:               "session-123",
		View:                   "enriched",
	})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	wantQuery := url.Values{
		"q":                         {"alice@example.com"},
		"cursor":                    {"cursor/with space"},
		"page_size":                 {"50"},
		"sort":                      {"-indexed_at"},
		"from":                      {"2026-01-01T00:00:00Z"},
		"to":                        {"2026-01-31T00:00:00Z"},
		"date_field":                {"indexed_at"},
		"log_id":                    {"log-123"},
		"has_log_id":                {"true"},
		"wildcard":                  {"true"},
		"logic":                     {"and"},
		"filter":                    {`{"and":[{"field":"domain","operator":"eq","value":"example.com"}]}`},
		"filter_id":                 {"0123456789abcdef01234567"},
		"domain[]":                  {"example.com"},
		"subdomain[]":               {"accounts.example.com"},
		"username[]":                {"alice"},
		"password[]":                {"secret"},
		"password_hash[]":           {"sha256:value"},
		"path[]":                    {"/Users/Alice/AppData"},
		"email[]":                   {"alice@example.com"},
		"email_domain[]":            {"example.com"},
		"ip[]":                      {"203.0.113.10"},
		"hwid[]":                    {"hwid-1"},
		"discord_id[]":              {"1234567890"},
		"source_type[]":             {"stealer"},
		"archive_hash[]":            {"archive-hash"},
		"canonical_credential_id[]": {"cred-1"},
		"fields[]":                  {"url_str", "domain"},
		"search_id":                 {"session-123"},
		"view":                      {"enriched"},
	}
	if gotPath != "/service/v2/stealer/search" {
		t.Fatalf("path = %s, want /service/v2/stealer/search", gotPath)
	}
	if !reflect.DeepEqual(gotQuery, wantQuery) {
		t.Fatalf("query = %#v, want %#v", gotQuery, wantQuery)
	}
	if resp == nil || !resp.Success || resp.Data == nil || resp.Data.NextCursor != "next-stealer" {
		t.Fatalf("unexpected response: %#v", resp)
	}
	if len(resp.Data.Items) != 1 || resp.Data.Items[0].LogID != "log-123" {
		t.Fatalf("unexpected items: %#v", resp.Data.Items)
	}
}

func TestStealerV2Service_Subdomain(t *testing.T) {
	client := createTestClient(t)

	t.Run("extract subdomains from stealer data", func(t *testing.T) {
		result, err := client.Stealer.Subdomain("google.com", "")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
	})

	t.Run("subdomain with query filter", func(t *testing.T) {
		result, err := client.Stealer.Subdomain("google.com", "mail")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
	})
}

func TestStealerV2Service_SubdomainRequestConstruction(t *testing.T) {
	alive := true
	isAlive := false
	var gotPath string
	var gotQuery url.Values

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"success": true,
			"message": "ok",
			"data": {
				"domain": "example.com",
				"subdomains": ["app.example.com", {"subdomain": "api.example.com", "alive": true}],
				"count": 2,
				"alive_results": {"app.example.com": {"alive": false}},
				"source": "stealer"
			},
			"_meta": {"service": {"id": "stealer-subdomain"}}
		}`))
	}))
	defer server.Close()

	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	result, err := client.Stealer.Subdomain("example.com", "legacy", &SubdomainOptions{
		Query:    "mail",
		Alive:    &alive,
		IsAlive:  &isAlive,
		SearchID: "sess_123",
	})
	if err != nil {
		t.Fatalf("Subdomain() error = %v", err)
	}

	wantQuery := url.Values{
		"domain":    []string{"example.com"},
		"q":         []string{"mail"},
		"alive":     []string{"true"},
		"is_alive":  []string{"false"},
		"search_id": []string{"sess_123"},
	}
	if gotPath != "/service/v2/stealer/subdomain" {
		t.Fatalf("path = %s, want /service/v2/stealer/subdomain", gotPath)
	}
	if !reflect.DeepEqual(gotQuery, wantQuery) {
		t.Fatalf("query = %#v, want %#v", gotQuery, wantQuery)
	}
	if result.Data == nil || result.Data.Source != "stealer" || len(result.Data.Subdomains) != 2 {
		t.Fatalf("unexpected subdomain response: %#v", result.Data)
	}
	if _, ok := result.Data.Subdomains[1].(map[string]interface{}); !ok {
		t.Fatalf("second subdomain = %#v, want object entry", result.Data.Subdomains[1])
	}
	if result.Data.AliveResults["app.example.com"].(map[string]interface{})["alive"] != false {
		t.Fatalf("alive_results = %#v", result.Data.AliveResults)
	}

	_, err = client.Stealer.ExtractSubdomainV2("example.com", &SubdomainOptions{
		Query:    "mail",
		SearchID: "sess_123",
	})
	if err != nil {
		t.Fatalf("ExtractSubdomainV2() error = %v", err)
	}
}
