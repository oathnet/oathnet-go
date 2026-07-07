package oathnet

import (
	"encoding/json"
	"testing"
)

func TestSearchService_Breach(t *testing.T) {
	client := createTestClient(t)

	t.Run("basic breach search", func(t *testing.T) {
		result, err := client.Search.Breach(TestBreachQuery, nil)
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

	t.Run("breach search with cursor", func(t *testing.T) {
		// First request
		result, err := client.Search.Breach(TestBreachQuery, nil)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}

		// Second request with cursor if available
		if result.Data != nil && result.Data.Cursor != "" {
			result2, err := client.Search.Breach(TestBreachQuery, &SearchOptions{
				Cursor: result.Data.Cursor,
			})
			if err != nil {
				t.Errorf("Unexpected error on paginated request: %v", err)
			}
			if !result2.Success {
				t.Error("Expected success to be true on paginated request")
			}
		}
	})

	t.Run("breach search with database filter", func(t *testing.T) {
		result, err := client.Search.Breach("ahmed", &SearchOptions{
			DBNames: "free.fr",
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
	})
}

func TestSearchService_Stealer(t *testing.T) {
	client := createTestClient(t)

	t.Run("basic stealer search", func(t *testing.T) {
		result, err := client.Search.Stealer("diddy", nil)
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

	t.Run("stealer results have LOG field", func(t *testing.T) {
		result, err := client.Search.Stealer("diddy", nil)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(result.Data.Results) > 0 {
			// Check first result has LOG field
			first := result.Data.Results[0]
			if first.LOG == "" {
				t.Error("Expected LOG field in stealer result")
			}
		}
	})
}

func TestSearchService_InitSession(t *testing.T) {
	client := createTestClient(t)

	t.Run("initialize search session", func(t *testing.T) {
		result, err := client.Search.InitSession("test@example.com")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
		if result.Data == nil || result.Data.Session.ID == "" {
			t.Error("Expected session ID to be set")
		}
	})
}

func TestLegacyBreachResultAcceptsArrayScalars(t *testing.T) {
	payload := []byte(`{
		"success": true,
		"data": {
			"results": [{
				"email": ["alice@example.com", "a@example.com"],
				"ip": ["198.51.100.10", "198.51.100.11"],
				"dbname": ["example", "example.org"]
			}],
			"results_found": 1,
			"results_shown": 1
		}
	}`)

	var resp BreachSearchResponse
	if err := json.Unmarshal(payload, &resp); err != nil {
		t.Fatalf("legacy breach response should tolerate list-valued scalars: %v", err)
	}
	if got := string(resp.Data.Results[0].Email); got != "alice@example.com,a@example.com" {
		t.Fatalf("unexpected email value: %q", got)
	}
	if got := string(resp.Data.Results[0].IP); got != "198.51.100.10,198.51.100.11" {
		t.Fatalf("unexpected ip value: %q", got)
	}
	if got := string(resp.Data.Results[0].DBName); got != "example,example.org" {
		t.Fatalf("unexpected dbname value: %q", got)
	}
}

func TestLegacyStealerResultAcceptsScalarLists(t *testing.T) {
	payload := []byte(`{
		"success": true,
		"data": {
			"results": [{
				"url": "https://example.com/login",
				"domain": "example.com",
				"email": "alice@example.com",
				"LOG": "log_123"
			}],
			"results_found": 1,
			"results_shown": 1
		}
	}`)

	var resp StealerSearchResponse
	if err := json.Unmarshal(payload, &resp); err != nil {
		t.Fatalf("legacy stealer response should tolerate scalar list fields: %v", err)
	}
	result := resp.Data.Results[0]
	if len(result.URL) != 1 || result.URL[0] != "https://example.com/login" {
		t.Fatalf("unexpected url values: %#v", result.URL)
	}
	if len(result.Domain) != 1 || result.Domain[0] != "example.com" {
		t.Fatalf("unexpected domain values: %#v", result.Domain)
	}
}
