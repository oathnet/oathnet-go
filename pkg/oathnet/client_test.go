package oathnet

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	t.Run("requires API key", func(t *testing.T) {
		_, err := NewClient("")
		if err == nil {
			t.Error("Expected error for empty API key")
		}
	})

	t.Run("creates client with API key", func(t *testing.T) {
		client, err := NewClient("test-api-key")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if client == nil {
			t.Error("Expected client to be non-nil")
		}
	})

	t.Run("accepts custom options", func(t *testing.T) {
		client, err := NewClient("test-api-key",
			WithBaseURL("https://custom.api.com"),
			WithTimeout(60*time.Second),
		)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if client == nil {
			t.Error("Expected client to be non-nil")
		}
		if client.baseURL != "https://custom.api.com" {
			t.Errorf("Expected custom base URL, got %s", client.baseURL)
		}
	})
}

func TestClientServices(t *testing.T) {
	client, err := NewClient("test-api-key")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	t.Run("has Search service", func(t *testing.T) {
		if client.Search == nil {
			t.Error("Expected Search service to be non-nil")
		}
	})

	t.Run("has OSINT service", func(t *testing.T) {
		if client.OSINT == nil {
			t.Error("Expected OSINT service to be non-nil")
		}
	})

	t.Run("has Stealer service", func(t *testing.T) {
		if client.Stealer == nil {
			t.Error("Expected Stealer service to be non-nil")
		}
	})

	t.Run("has Victims service", func(t *testing.T) {
		if client.Victims == nil {
			t.Error("Expected Victims service to be non-nil")
		}
	})

	t.Run("has FileSearch service", func(t *testing.T) {
		if client.FileSearch == nil {
			t.Error("Expected FileSearch service to be non-nil")
		}
	})

	t.Run("has Exports service", func(t *testing.T) {
		if client.Exports == nil {
			t.Error("Expected Exports service to be non-nil")
		}
	})

	t.Run("has Bulk service", func(t *testing.T) {
		if client.Bulk == nil {
			t.Error("Expected Bulk service to be non-nil")
		}
	})

	t.Run("has Utility service", func(t *testing.T) {
		if client.Utility == nil {
			t.Error("Expected Utility service to be non-nil")
		}
	})
}

func TestClientIntegration(t *testing.T) {
	client := createTestClient(t)

	t.Run("successful API call", func(t *testing.T) {
		result, err := client.Search.Breach(TestBreachQuery, nil)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
	})
}

func TestInvalidAPIKey(t *testing.T) {
	client, _ := NewClient("invalid-api-key")

	_, err := client.Search.Breach("test", nil)
	if err == nil {
		t.Error("Expected error for invalid API key")
	}
}
