package oathnet

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestAIService_CreateRequestConstruction(t *testing.T) {
	var gotMethod string
	var gotPath string
	var gotQuery url.Values
	var gotAPIKey string
	var gotBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.EscapedPath()
		gotQuery = r.URL.Query()
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
			"filter_id": "abcdefabcdefabcdefabcdef",
			"filter": {
				"and": [
					{"field": "country", "operator": "eq", "value": "fr"},
					{"field": "age", "operator": "gt", "value": "30"}
				]
			}
		}`))
	}))
	defer server.Close()

	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.AI.Create(V2AIFilterRequest{
		Query:    "French users over 30",
		Index:    AIFilterIndexBreach,
		FilterID: "0123456789abcdef01234567",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Fatalf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/service/v2/ai/filter" {
		t.Fatalf("path = %s, want /service/v2/ai/filter", gotPath)
	}
	if len(gotQuery) != 0 {
		t.Fatalf("query = %v, want empty", gotQuery)
	}
	if gotAPIKey != "test-key" {
		t.Fatalf("x-api-key = %q, want test-key", gotAPIKey)
	}

	wantBody := map[string]interface{}{
		"query":     "French users over 30",
		"index":     "breach",
		"filter_id": "0123456789abcdef01234567",
	}
	if !reflect.DeepEqual(gotBody, wantBody) {
		t.Fatalf("body = %#v, want %#v", gotBody, wantBody)
	}

	if resp == nil || resp.FilterID != "abcdefabcdefabcdefabcdef" {
		t.Fatalf("unexpected response: %#v", resp)
	}
	clauses, ok := resp.Filter["and"].([]interface{})
	if !ok || len(clauses) != 2 {
		t.Fatalf("filter and clauses = %#v, want 2 clauses", resp.Filter["and"])
	}
	firstClause, ok := clauses[0].(map[string]interface{})
	if !ok || firstClause["field"] != "country" || firstClause["operator"] != "eq" || firstClause["value"] != "fr" {
		t.Fatalf("unexpected first filter clause: %#v", clauses[0])
	}
}

func TestAIService_CreateOmitsOptionalFields(t *testing.T) {
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
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"filter_id":"abcdefabcdefabcdefabcdef","filter":{"field":"email","operator":"contains","value":"alice"}}`))
	}))
	defer server.Close()

	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if _, err := client.AI.Create(V2AIFilterRequest{Query: "alice"}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	wantBody := map[string]interface{}{
		"query": "alice",
	}
	if !reflect.DeepEqual(gotBody, wantBody) {
		t.Fatalf("body = %#v, want %#v", gotBody, wantBody)
	}
}

func TestAIService_GetContextRequestConstruction(t *testing.T) {
	var gotMethod string
	var gotPath string
	var gotQuery url.Values
	var gotAPIKey string
	var gotBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.EscapedPath()
		gotQuery = r.URL.Query()
		gotAPIKey = r.Header.Get("x-api-key")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("ReadAll() error = %v", err)
		}
		gotBody = body

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": "filter/with space",
			"index_type": "breach",
			"query": "French users",
			"filter": {"field": "country", "operator": "eq", "value": "fr"},
			"sample_data": {"email": "alice@example.fr"},
			"field_values": {"country": ["fr"], "city": ["Paris"]},
			"total_hits": 12345678901,
			"history": [
				{"query": "French users", "filter": {"field": "country", "operator": "eq", "value": "fr"}}
			],
			"source": "ai",
			"parent_id": "0123456789abcdef01234567",
			"created_at": "2026-07-06T00:00:00Z",
			"expires_at": "2026-07-06T01:00:00Z"
		}`))
	}))
	defer server.Close()

	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.AI.GetContext("filter/with space")
	if err != nil {
		t.Fatalf("GetContext() error = %v", err)
	}

	if gotMethod != http.MethodGet {
		t.Fatalf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/service/v2/ai/filter/filter%2Fwith%20space" {
		t.Fatalf("path = %s, want escaped filter context path", gotPath)
	}
	if len(gotQuery) != 0 {
		t.Fatalf("query = %v, want empty", gotQuery)
	}
	if gotAPIKey != "test-key" {
		t.Fatalf("x-api-key = %q, want test-key", gotAPIKey)
	}
	if len(gotBody) != 0 {
		t.Fatalf("GET body = %s, want empty", string(gotBody))
	}

	if resp == nil || resp.ID != "filter/with space" || resp.IndexType != "breach" {
		t.Fatalf("unexpected context response: %#v", resp)
	}
	if resp.TotalHits != 12345678901 {
		t.Fatalf("total_hits = %d, want 12345678901", resp.TotalHits)
	}
	if resp.FieldValues["city"][0] != "Paris" {
		t.Fatalf("field_values = %#v, want city Paris", resp.FieldValues)
	}
	if len(resp.History) != 1 || resp.History[0].Query != "French users" || resp.History[0].Filter["field"] != "country" {
		t.Fatalf("history = %#v, want decoded filter history", resp.History)
	}
	if resp.SampleData["email"] != "alice@example.fr" || resp.ParentID != "0123456789abcdef01234567" || resp.ExpiresAt == "" {
		t.Fatalf("unexpected context metadata: %#v", resp)
	}
}

func TestAIService_ErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"filter_id not found or expired"}`))
	}))
	defer server.Close()

	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	_, err = client.AI.GetContext("0123456789abcdef01234567")
	var notFound *NotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("error = %T %v, want *NotFoundError", err, err)
	}
	if notFound.Message != "filter_id not found or expired" {
		t.Fatalf("error message = %q, want filter_id not found or expired", notFound.Message)
	}
}
