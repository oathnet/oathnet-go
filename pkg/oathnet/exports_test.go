package oathnet

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
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

func TestExportsService_OpenAPIRequestConstruction(t *testing.T) {
	tests := []struct {
		name         string
		call         func(*Client) error
		responseBody string
		contentType  string
		wantMethod   string
		wantPath     string
		wantQuery    url.Values
		wantBody     map[string]interface{}
	}{
		{
			name: "CreateExportV2 sends current export body",
			call: func(client *Client) error {
				resp, err := client.Exports.CreateExportV2("breach", &ExportCreateOptions{
					Format:      "html",
					Limit:       10,
					Fields:      []string{"email"},
					Service:     "breach",
					Search:      map[string]interface{}{"query": "user@example.com"},
					QueryConfig: map[string]interface{}{"filter_id": "flt-123"},
				})
				if err != nil {
					return err
				}
				if resp == nil || !resp.Success || resp.Data == nil || resp.Data.JobID != "export-1" {
					t.Fatalf("unexpected create response: %#v", resp)
				}
				return nil
			},
			responseBody: `{"job_id":"export-1","status":"queued","request":{"type":"breach","service":"breach","format":"html","limit":10,"fields":["email"],"request_count":1}}`,
			wantMethod:   http.MethodPost,
			wantPath:     "/service/v2/exports",
			wantQuery:    url.Values{},
			wantBody: map[string]interface{}{
				"type":    "breach",
				"format":  "html",
				"limit":   float64(10),
				"fields":  []interface{}{"email"},
				"service": "breach",
				"search": map[string]interface{}{
					"query": "user@example.com",
				},
				"query_config": map[string]interface{}{
					"filter_id": "flt-123",
				},
			},
		},
		{
			name: "GetExportV2 escapes job ID",
			call: func(client *Client) error {
				resp, err := client.Exports.GetExportV2("job/with space")
				if err != nil {
					return err
				}
				if resp == nil || resp.Data == nil || resp.Data.JobID != "job/with space" {
					t.Fatalf("unexpected status response: %#v", resp)
				}
				return nil
			},
			responseBody: `{"job_id":"job/with space","status":"completed"}`,
			wantMethod:   http.MethodGet,
			wantPath:     "/service/v2/exports/job%2Fwith%20space",
			wantQuery:    url.Values{},
		},
		{
			name: "DownloadExportV2 returns raw bytes from escaped path",
			call: func(client *Client) error {
				data, err := client.Exports.DownloadExportV2("job/with space")
				if err != nil {
					return err
				}
				if string(data) != "email,source\nuser@example.com,breach\n" {
					t.Fatalf("download bytes = %q", string(data))
				}
				outputPath := filepath.Join(t.TempDir(), "export.csv")
				if err := client.Exports.Download("job/with space", outputPath); err != nil {
					return err
				}
				saved, err := os.ReadFile(outputPath)
				if err != nil {
					return err
				}
				if string(saved) != string(data) {
					t.Fatalf("saved bytes = %q, want %q", string(saved), string(data))
				}
				return nil
			},
			responseBody: "email,source\nuser@example.com,breach\n",
			contentType:  "text/csv",
			wantMethod:   http.MethodGet,
			wantPath:     "/service/v2/exports/job%2Fwith%20space/download",
			wantQuery:    url.Values{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotMethod string
			var gotPath string
			var gotQuery url.Values
			var gotBody map[string]interface{}
			var gotAPIKey string

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotMethod = r.Method
				gotPath = r.URL.EscapedPath()
				gotQuery = r.URL.Query()
				gotAPIKey = r.Header.Get("x-api-key")

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("ReadAll() error = %v", err)
				}
				if tt.wantBody != nil {
					if err := json.Unmarshal(body, &gotBody); err != nil {
						t.Fatalf("Unmarshal request body error = %v; body=%s", err, string(body))
					}
				} else if len(body) > 0 {
					t.Fatalf("request body = %s, want empty", string(body))
				}

				contentType := tt.contentType
				if contentType == "" {
					contentType = "application/json"
				}
				w.Header().Set("Content-Type", contentType)
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client, err := NewClient("test-key", WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			if err := tt.call(client); err != nil {
				t.Fatalf("exports call error = %v", err)
			}
			if gotAPIKey != "test-key" {
				t.Fatalf("x-api-key = %q, want test-key", gotAPIKey)
			}
			if gotMethod != tt.wantMethod {
				t.Fatalf("method = %s, want %s", gotMethod, tt.wantMethod)
			}
			if gotPath != tt.wantPath {
				t.Fatalf("path = %s, want %s", gotPath, tt.wantPath)
			}
			if !reflect.DeepEqual(gotQuery, tt.wantQuery) {
				t.Fatalf("query = %v, want %v", gotQuery, tt.wantQuery)
			}
			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Fatalf("body = %#v, want %#v", gotBody, tt.wantBody)
			}
		})
	}
}
