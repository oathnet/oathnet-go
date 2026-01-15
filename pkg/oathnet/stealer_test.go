package oathnet

import (
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
