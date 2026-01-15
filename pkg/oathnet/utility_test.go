package oathnet

import (
	"testing"
)

func TestUtilityService_DBNameAutocomplete(t *testing.T) {
	client := createTestClient(t)

	t.Run("autocomplete database names", func(t *testing.T) {
		result, err := client.Utility.DBNameAutocomplete("linkedin")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result == nil {
			t.Error("Expected result to be non-nil")
		}
	})

	t.Run("autocomplete with common prefix", func(t *testing.T) {
		result, err := client.Utility.DBNameAutocomplete("face")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result == nil {
			t.Error("Expected result to be non-nil")
		}
	})
}
