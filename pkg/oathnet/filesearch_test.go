package oathnet

import (
	"testing"
	"time"
)

func TestFileSearchService_Create(t *testing.T) {
	client := createTestClient(t)

	t.Run("create file search job with log IDs", func(t *testing.T) {
		// First get a log ID from victims search
		searchResult, err := client.Victims.Search(TestVictimsQuery, &VictimsSearchOptions{
			PageSize: 5,
		})
		if err != nil {
			t.Skipf("Skipping: Victims search failed: %v", err)
		}

		var logIDs []string
		if searchResult.Data != nil && len(searchResult.Data.Items) > 0 {
			for _, r := range searchResult.Data.Items {
				if r.LogID != "" {
					logIDs = append(logIDs, r.LogID)
					if len(logIDs) >= 3 {
						break
					}
				}
			}
		}

		if len(logIDs) == 0 {
			t.Skip("No log IDs available for file search")
		}

		result, err := client.FileSearch.Create("password", &FileSearchCreateOptions{
			SearchMode:     "literal",
			LogIDs:         logIDs,
			MaxMatches:     5,
			IncludeMatches: true,
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		if result == nil {
			t.Error("Expected result to be non-nil")
			return
		}
		if result.Data == nil || result.Data.JobID == "" {
			t.Error("Expected job ID to be set")
		}
	})
}

func TestFileSearchService_GetStatus(t *testing.T) {
	client := createTestClient(t)

	t.Run("get job status", func(t *testing.T) {
		// First get log IDs
		searchResult, err := client.Victims.Search(TestVictimsQuery, &VictimsSearchOptions{
			PageSize: 5,
		})
		if err != nil {
			t.Skipf("Skipping: Victims search failed: %v", err)
		}

		var logIDs []string
		if searchResult.Data != nil && len(searchResult.Data.Items) > 0 {
			for _, r := range searchResult.Data.Items {
				if r.LogID != "" {
					logIDs = append(logIDs, r.LogID)
					if len(logIDs) >= 3 {
						break
					}
				}
			}
		}

		if len(logIDs) == 0 {
			t.Skip("No log IDs available")
		}

		job, err := client.FileSearch.Create("password", &FileSearchCreateOptions{
			LogIDs:         logIDs,
			MaxMatches:     5,
			IncludeMatches: true,
		})
		if err != nil {
			t.Skipf("Skipping: Create failed: %v", err)
		}
		if job == nil || job.Data == nil || job.Data.JobID == "" {
			t.Skip("Skipping: No job ID returned")
		}

		status, err := client.FileSearch.GetStatus(job.Data.JobID)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		if status == nil || status.Data == nil {
			t.Error("Expected data to be non-nil")
		}
	})
}

func TestFileSearchService_Search(t *testing.T) {
	client := createTestClient(t)

	t.Run("search convenience method", func(t *testing.T) {
		// First get log IDs
		searchResult, err := client.Victims.Search(TestVictimsQuery, &VictimsSearchOptions{
			PageSize: 5,
		})
		if err != nil {
			t.Skipf("Skipping: Victims search failed: %v", err)
		}

		var logIDs []string
		if searchResult.Data != nil && len(searchResult.Data.Items) > 0 {
			for _, r := range searchResult.Data.Items {
				if r.LogID != "" {
					logIDs = append(logIDs, r.LogID)
					if len(logIDs) >= 3 {
						break
					}
				}
			}
		}

		if len(logIDs) == 0 {
			t.Skip("No log IDs available")
		}

		result, err := client.FileSearch.Search("password", &FileSearchCreateOptions{
			SearchMode:     "literal",
			LogIDs:         logIDs,
			MaxMatches:     3,
			IncludeMatches: true,
		}, 60*time.Second)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		if result == nil || result.Data == nil {
			t.Error("Expected result and data to be non-nil")
			return
		}
		if result.Data.Status != "completed" && result.Data.Status != "canceled" {
			t.Errorf("Expected status completed or canceled, got %s", result.Data.Status)
		}
	})
}
