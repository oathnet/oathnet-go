package oathnet

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestSearchService_InitSessionRequestConstruction(t *testing.T) {
	var gotMethod string
	var gotPath string
	var gotAPIKey string
	var gotBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.EscapedPath()
		gotAPIKey = r.Header.Get("x-api-key")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("ReadAll() error = %v", err)
		}
		if err := json.Unmarshal(body, &gotBody); err != nil {
			t.Fatalf("Unmarshal request body error = %v; body=%s", err, string(body))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"success": true,
			"message": "Search session initialized successfully",
			"data": {
				"session": {
					"id": "sess_abc123",
					"query": "test@example.com",
					"search_type": "email",
					"status": "active",
					"created_at": "2026-07-06T00:00:00Z",
					"expires_at": "2026-07-06T01:00:00Z",
					"duration_minutes": 60
				},
				"user": {
					"plan": "Pro",
					"plan_type": "pro",
					"daily_lookups": {
						"used": 1,
						"remaining": 999,
						"limit": 1000,
						"is_unlimited": false
					}
				},
				"services": {
					"breach": {
						"name": "Breach Search",
						"service_id": "breach",
						"category": "search",
						"is_available": true,
						"is_premium": false,
						"session_quota": 100,
						"today_usage": 1,
						"recommended_quota": 50
					}
				},
				"summary": {
					"total_services": 1,
					"available_services": 1,
					"session_expires_in_minutes": 60
				}
			}
		}`))
	}))
	defer server.Close()

	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.Search.InitSession("test@example.com", &SearchSessionOptions{
		SearchType: "email",
	})
	if err != nil {
		t.Fatalf("InitSession() error = %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Fatalf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/service/search/init" {
		t.Fatalf("path = %s, want /service/search/init", gotPath)
	}
	if gotAPIKey != "test-key" {
		t.Fatalf("x-api-key = %q, want test-key", gotAPIKey)
	}
	wantBody := map[string]interface{}{
		"query":       "test@example.com",
		"search_type": "email",
	}
	if !reflect.DeepEqual(gotBody, wantBody) {
		t.Fatalf("body = %#v, want %#v", gotBody, wantBody)
	}
	if resp == nil || !resp.Success || resp.Data == nil || resp.Data.Session == nil {
		t.Fatalf("unexpected response: %#v", resp)
	}
	if resp.Data.Session.Status != "active" || resp.Data.Session.DurationMinutes != 60 {
		t.Fatalf("session metadata = %#v", resp.Data.Session)
	}
	if resp.Data.Services["breach"].RecommendedQuota != 50 {
		t.Fatalf("services = %#v", resp.Data.Services)
	}
	if resp.Data.Summary == nil || resp.Data.Summary.AvailableServices != 1 {
		t.Fatalf("summary = %#v", resp.Data.Summary)
	}
}

func TestSearchService_InitSessionOmitsEmptySearchType(t *testing.T) {
	var gotBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("ReadAll() error = %v", err)
		}
		if err := json.Unmarshal(body, &gotBody); err != nil {
			t.Fatalf("Unmarshal request body error = %v; body=%s", err, string(body))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"session":{"id":"sess_1","query":"example.com"}}}`))
	}))
	defer server.Close()

	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.Search.InitSession("example.com")
	if err != nil {
		t.Fatalf("InitSession() error = %v", err)
	}

	wantBody := map[string]interface{}{"query": "example.com"}
	if !reflect.DeepEqual(gotBody, wantBody) {
		t.Fatalf("body = %#v, want %#v", gotBody, wantBody)
	}
}

func TestSearchService_LegacyBreachRequestAndResponseParity(t *testing.T) {
	var gotMethod string
	var gotPath string
	var gotAPIKey string
	var gotQuery url.Values

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.EscapedPath()
		gotAPIKey = r.Header.Get("x-api-key")
		gotQuery = r.URL.Query()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"success": true,
			"message": "Request completed successfully",
			"data": {
				"results": [{
					"dbname": "deezer.com",
					"country": ["US"],
					"date_birth": ["2004-01-01 00:00:00.0000000"],
					"email": "user@gmail.com",
					"gender": "F",
					"language": "us",
					"created_at": "2017-04-16 00:00:00.0000000",
					"username": ["WinterFox"],
					"id": "09c1df9ddb7b657fa7f2d6778e40ac0a",
					"_version_": 1845592909587415000,
					"source_specific_field": "preserved"
				}],
				"results_found": 150,
				"results_shown": 25,
				"nextCursorMark": "AoE/next"
			}
		}`))
	}))
	defer server.Close()

	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.Search.Breach("winter fox@example.com", &SearchOptions{
		Cursor:   "AoE/prev",
		DBNames:  "linkedin,adobe",
		SearchID: "sess_123",
	})
	if err != nil {
		t.Fatalf("Breach() error = %v", err)
	}

	if gotMethod != http.MethodGet {
		t.Fatalf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/service/search-breach" {
		t.Fatalf("path = %s, want /service/search-breach", gotPath)
	}
	if gotAPIKey != "test-key" {
		t.Fatalf("x-api-key = %q, want test-key", gotAPIKey)
	}
	wantQuery := url.Values{
		"q":         []string{"winter fox@example.com"},
		"cursor":    []string{"AoE/prev"},
		"dbnames":   []string{"linkedin,adobe"},
		"search_id": []string{"sess_123"},
	}
	if !reflect.DeepEqual(gotQuery, wantQuery) {
		t.Fatalf("query = %#v, want %#v", gotQuery, wantQuery)
	}

	if resp == nil || !resp.Success || resp.Data == nil {
		t.Fatalf("unexpected response: %#v", resp)
	}
	if resp.Data.NextCursorMark != "AoE/next" || resp.Data.Cursor != "AoE/next" {
		t.Fatalf("cursor fields = nextCursorMark:%q cursor:%q, want AoE/next", resp.Data.NextCursorMark, resp.Data.Cursor)
	}
	if resp.Data.ResultsFound != 150 || resp.Data.ResultsShown != 25 || len(resp.Data.Results) != 1 {
		t.Fatalf("result metadata = %#v", resp.Data)
	}
	first := resp.Data.Results[0]
	if first.DBName != "deezer.com" || first.Gender != "F" || first.Language != "us" || first.Version != 1845592909587415000 {
		t.Fatalf("breach result = %#v", first)
	}
	if got := first.Extra["source_specific_field"]; got != "preserved" {
		t.Fatalf("extra source field = %#v, want preserved", got)
	}
}

func TestSearchService_LegacyStealerRequestAndResponseParity(t *testing.T) {
	var gotMethod string
	var gotPath string
	var gotAPIKey string
	var gotQuery url.Values

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.EscapedPath()
		gotAPIKey = r.Header.Get("x-api-key")
		gotQuery = r.URL.Query()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"success": true,
			"message": "Request completed successfully",
			"data": {
				"results": [{
					"id": "03b4f7e0ae1e2fac",
					"LOG": "https://aurorafn.dev/signup:Diddy:diddyisdaddy1",
					"domain": ["aurorafn.dev"],
					"subdomain": ["app.aurorafn.dev"],
					"path": ["/signup"],
					"email": ["t-diddy@hotmail.com"],
					"_version_": 1835054175050793000,
					"malware_family": "redline"
				}],
				"results_found": 500,
				"results_shown": 25,
				"nextCursorMark": "AoE/stealer"
			}
		}`))
	}))
	defer server.Close()

	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.Search.Stealer("diddy", &SearchOptions{
		Cursor:   "AoE/prev",
		DBNames:  "stealer-db",
		SearchID: "sess_456",
	})
	if err != nil {
		t.Fatalf("Stealer() error = %v", err)
	}

	if gotMethod != http.MethodGet {
		t.Fatalf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/service/search-stealer" {
		t.Fatalf("path = %s, want /service/search-stealer", gotPath)
	}
	if gotAPIKey != "test-key" {
		t.Fatalf("x-api-key = %q, want test-key", gotAPIKey)
	}
	wantQuery := url.Values{
		"q":         []string{"diddy"},
		"cursor":    []string{"AoE/prev"},
		"dbnames":   []string{"stealer-db"},
		"search_id": []string{"sess_456"},
	}
	if !reflect.DeepEqual(gotQuery, wantQuery) {
		t.Fatalf("query = %#v, want %#v", gotQuery, wantQuery)
	}

	if resp == nil || !resp.Success || resp.Data == nil {
		t.Fatalf("unexpected response: %#v", resp)
	}
	if resp.Data.NextCursorMark != "AoE/stealer" || resp.Data.Cursor != "AoE/stealer" {
		t.Fatalf("cursor fields = nextCursorMark:%q cursor:%q, want AoE/stealer", resp.Data.NextCursorMark, resp.Data.Cursor)
	}
	if resp.Data.ResultsFound != 500 || resp.Data.ResultsShown != 25 || len(resp.Data.Results) != 1 {
		t.Fatalf("result metadata = %#v", resp.Data)
	}
	first := resp.Data.Results[0]
	if first.ID != "03b4f7e0ae1e2fac" || first.LOG == "" || len(first.Subdomain) != 1 || len(first.Path) != 1 || first.Version != 1835054175050793000 {
		t.Fatalf("stealer result = %#v", first)
	}
	if got := first.Extra["malware_family"]; got != "redline" {
		t.Fatalf("extra source field = %#v, want redline", got)
	}
}

func TestSearchService_LegacyCursorAliasDecoding(t *testing.T) {
	var breach BreachSearchResponse
	if err := json.Unmarshal([]byte(`{"success":true,"data":{"results":[],"next_cursor_mark":"snake-cursor"}}`), &breach); err != nil {
		t.Fatalf("Unmarshal breach response error = %v", err)
	}
	if breach.Data == nil || breach.Data.NextCursorMark != "snake-cursor" || breach.Data.Cursor != "snake-cursor" {
		t.Fatalf("breach cursor aliases = %#v", breach.Data)
	}

	var stealer StealerSearchResponse
	if err := json.Unmarshal([]byte(`{"success":true,"data":{"results":[],"cursor":"legacy-cursor"}}`), &stealer); err != nil {
		t.Fatalf("Unmarshal stealer response error = %v", err)
	}
	if stealer.Data == nil || stealer.Data.NextCursorMark != "legacy-cursor" || stealer.Data.Cursor != "legacy-cursor" {
		t.Fatalf("stealer cursor aliases = %#v", stealer.Data)
	}
}
