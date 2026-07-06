package oathnet

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
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
