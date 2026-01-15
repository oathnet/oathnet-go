package oathnet

import (
	"testing"
	"time"
)

func TestExportsService_Create(t *testing.T) {
	client := createTestClient(t)

	t.Run("create docs export", func(t *testing.T) {
		result, err := client.Exports.Create("docs", &ExportCreateOptions{
			Format: "jsonl",
			Limit:  100,
			Search: map[string]string{
				"query": "gmail.com",
			},
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
		if result.Data == nil || result.Data.JobID == "" {
			t.Error("Expected job ID to be set")
		}
	})

	t.Run("create victims export", func(t *testing.T) {
		result, err := client.Exports.Create("victims", &ExportCreateOptions{
			Format: "jsonl",
			Limit:  100,
			Search: map[string]string{
				"query": "gmail",
			},
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		if result.Data == nil || result.Data.JobID == "" {
			t.Error("Expected job ID to be set")
		}
	})

	t.Run("create CSV export", func(t *testing.T) {
		result, err := client.Exports.Create("docs", &ExportCreateOptions{
			Format: "csv",
			Limit:  100,
			Search: map[string]string{
				"query": "gmail.com",
			},
		})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
	})

	t.Run("create export with fields", func(t *testing.T) {
		result, err := client.Exports.Create("docs", &ExportCreateOptions{
			Format: "jsonl",
			Limit:  100,
			Fields: []string{"email", "password", "domain"},
			Search: map[string]string{
				"query": "gmail.com",
			},
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

func TestExportsService_GetStatus(t *testing.T) {
	client := createTestClient(t)

	t.Run("get export status", func(t *testing.T) {
		job, err := client.Exports.Create("docs", &ExportCreateOptions{
			Format: "jsonl",
			Limit:  100,
			Search: map[string]string{
				"query": "gmail.com",
			},
		})
		if err != nil {
			t.Fatalf("Failed to create export: %v", err)
		}

		if job.Data == nil || job.Data.JobID == "" {
			t.Fatal("No job ID returned")
		}

		status, err := client.Exports.GetStatus(job.Data.JobID)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		if !status.Success {
			t.Error("Expected success to be true")
		}
	})
}

func TestExportsService_WaitForCompletion(t *testing.T) {
	client := createTestClient(t)

	t.Run("wait for export completion", func(t *testing.T) {
		job, err := client.Exports.Create("docs", &ExportCreateOptions{
			Format: "jsonl",
			Limit:  100,
			Search: map[string]string{
				"query": "gmail.com",
			},
		})
		if err != nil {
			t.Fatalf("Failed to create export: %v", err)
		}

		if job.Data == nil || job.Data.JobID == "" {
			t.Fatal("No job ID returned")
		}

		result, err := client.Exports.WaitForCompletion(
			job.Data.JobID,
			time.Second,
			120*time.Second,
		)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		if !result.Success {
			t.Error("Expected success to be true")
		}
		if result.Data.Status != "completed" && result.Data.Status != "canceled" {
			t.Errorf("Expected status completed or canceled, got %s", result.Data.Status)
		}
	})
}

func TestExportsService_Download(t *testing.T) {
	t.Skip("Skipping: Export download API has known issues (500 error)")
}
