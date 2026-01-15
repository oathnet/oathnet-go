package oathnet

import (
	"testing"
)

func TestVictimsService_Search(t *testing.T) {
	client := createTestClient(t)

	t.Run("basic victims search", func(t *testing.T) {
		result, err := client.Victims.Search(TestVictimsQuery, &VictimsSearchOptions{
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

	t.Run("cursor pagination", func(t *testing.T) {
		result1, err := client.Victims.Search(TestVictimsQuery, &VictimsSearchOptions{
			PageSize: 5,
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if result1.Data != nil && result1.Data.NextCursor != "" {
			result2, err := client.Victims.Search(TestVictimsQuery, &VictimsSearchOptions{
				PageSize: 5,
				Cursor:   result1.Data.NextCursor,
			})
			if err != nil {
				t.Errorf("Unexpected error on page 2: %v", err)
			}
			if !result2.Success {
				t.Error("Expected success to be true on page 2")
			}
		}
	})

	t.Run("search with email filter", func(t *testing.T) {
		result, err := client.Victims.Search("", &VictimsSearchOptions{
			Email:    "gmail.com",
			PageSize: 5,
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
	})

	t.Run("wildcard search", func(t *testing.T) {
		result, err := client.Victims.Search(TestVictimsQuery, &VictimsSearchOptions{
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

func TestVictimsService_GetManifest(t *testing.T) {
	client := createTestClient(t)

	t.Run("get victim manifest", func(t *testing.T) {
		// First get a log ID from search
		searchResult, err := client.Victims.Search(TestVictimsQuery, &VictimsSearchOptions{
			PageSize: 5,
		})
		if err != nil {
			t.Skipf("Search failed: %v", err)
		}

		var logID string
		if searchResult.Data != nil && len(searchResult.Data.Items) > 0 {
			for _, r := range searchResult.Data.Items {
				if r.LogID != "" {
					logID = r.LogID
					break
				}
			}
		}

		if logID == "" {
			t.Skip("No log ID available")
		}

		manifest, err := client.Victims.GetManifest(logID)
		if err != nil {
			// Manifest may not be available for all logs
			t.Logf("GetManifest error (may be expected): %v", err)
			return
		}
		if manifest == nil {
			t.Error("Expected manifest to be non-nil")
		}
	})
}

func TestVictimsService_GetFile(t *testing.T) {
	client := createTestClient(t)

	t.Run("get victim file", func(t *testing.T) {
		// First get a log ID and file ID
		searchResult, err := client.Victims.Search(TestVictimsQuery, &VictimsSearchOptions{
			PageSize: 5,
		})
		if err != nil {
			t.Skipf("Search failed: %v", err)
		}

		var logID string
		if searchResult.Data != nil && len(searchResult.Data.Items) > 0 {
			for _, r := range searchResult.Data.Items {
				if r.LogID != "" {
					logID = r.LogID
					break
				}
			}
		}

		if logID == "" {
			t.Skip("No log ID available")
		}

		manifest, err := client.Victims.GetManifest(logID)
		if err != nil || manifest == nil || manifest.VictimTree == nil {
			t.Skip("No manifest or files available")
		}

		// Find a file in the tree
		var fileID string
		findFile(manifest.VictimTree, &fileID)
		if fileID == "" {
			t.Skip("No file ID available")
		}

		content, err := client.Victims.GetFile(logID, fileID)
		if err != nil {
			t.Logf("GetFile error (may be expected): %v", err)
			return
		}
		if len(content) == 0 {
			t.Log("File content is empty (may be expected)")
		}
	})
}

// findFile recursively searches for a file ID in the manifest tree
func findFile(node *VictimManifestNode, fileID *string) {
	if node == nil || *fileID != "" {
		return
	}
	if node.Type == "file" && node.ID != "" {
		*fileID = node.ID
		return
	}
	for _, child := range node.Children {
		findFile(&child, fileID)
	}
}
