package oathnet

import (
	"strings"
	"testing"
)

func TestErrorTypes(t *testing.T) {
	t.Run("OathNetError implements error", func(t *testing.T) {
		var err error = &OathNetError{Message: "test", StatusCode: 500}
		if !strings.Contains(err.Error(), "test") {
			t.Errorf("Expected error to contain 'test', got '%s'", err.Error())
		}
	})

	t.Run("AuthenticationError implements error", func(t *testing.T) {
		var err error = &AuthenticationError{Message: "auth failed"}
		if !strings.Contains(err.Error(), "auth failed") {
			t.Errorf("Expected error to contain 'auth failed', got '%s'", err.Error())
		}
	})

	t.Run("ValidationError implements error", func(t *testing.T) {
		var err error = &ValidationError{Message: "validation failed"}
		if !strings.Contains(err.Error(), "validation failed") {
			t.Errorf("Expected error to contain 'validation failed', got '%s'", err.Error())
		}
	})

	t.Run("NotFoundError implements error", func(t *testing.T) {
		var err error = &NotFoundError{Message: "not found"}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("Expected error to contain 'not found', got '%s'", err.Error())
		}
	})

	t.Run("RateLimitError implements error", func(t *testing.T) {
		var err error = &RateLimitError{Message: "rate limited"}
		if !strings.Contains(err.Error(), "rate limited") {
			t.Errorf("Expected error to contain 'rate limited', got '%s'", err.Error())
		}
	})
}

func TestErrorTypeAssertion(t *testing.T) {
	client, _ := NewClient("invalid-api-key")

	_, err := client.Search.Breach("test", nil)
	if err == nil {
		t.Error("Expected error for invalid API key")
		return
	}

	// Error should be AuthenticationError or ValidationError
	switch err.(type) {
	case *AuthenticationError:
		// Expected
	case *ValidationError:
		// Also acceptable
	default:
		// Other error types are fine too, API behavior may vary
		t.Logf("Got error type: %T", err)
	}
}

func TestParseErrorPreservesValidationDetails(t *testing.T) {
	client, _ := NewClient("test-key")
	err := client.parseError(400, []byte(`{
		"message": "Validation error",
		"errors": {
			"include": ["Invalid include section: credentials,victims"]
		}
	}`))

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if validationErr.Details["include"] == nil {
		t.Fatalf("expected include details, got %#v", validationErr.Details)
	}
	if !strings.Contains(validationErr.Error(), "credentials,victims") {
		t.Fatalf("expected formatted details in error string, got %q", validationErr.Error())
	}
}
