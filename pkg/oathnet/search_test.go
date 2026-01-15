package oathnet

import (
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
