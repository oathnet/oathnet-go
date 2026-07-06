package oathnet

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
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

func TestFileSearchService_OpenAPIRequestConstruction(t *testing.T) {
	tests := []struct {
		name         string
		call         func(*Client) error
		wantMethod   string
		wantPath     string
		wantQuery    url.Values
		wantBody     map[string]interface{}
		responseBody string
	}{
		{
			name: "CreateFileSearchV2 sends current body",
			call: func(client *Client) error {
				resp, err := client.FileSearch.CreateFileSearchV2("api[_-]?key", &FileSearchCreateOptions{
					SearchMode:     "regex",
					LogIDs:         []string{"log-1"},
					IncludeMatches: true,
					CaseSensitive:  true,
					ContextLines:   3,
					FilePattern:    "*.txt",
					MaxMatches:     20,
				})
				if err != nil {
					return err
				}
				if resp == nil || resp.Data == nil || resp.Data.JobID != "file-job-1" {
					t.Fatalf("unexpected create response: %#v", resp)
				}
				return nil
			},
			wantMethod: http.MethodPost,
			wantPath:   "/service/v2/file-search",
			wantQuery:  url.Values{},
			wantBody: map[string]interface{}{
				"expression":      "api[_-]?key",
				"search_mode":     "regex",
				"log_ids":         []interface{}{"log-1"},
				"include_matches": true,
				"case_sensitive":  true,
				"context_lines":   float64(3),
				"file_pattern":    "*.txt",
				"max_matches":     float64(20),
			},
			responseBody: `{"job_id":"file-job-1","status":"queued"}`,
		},
		{
			name: "GetFileSearchV2 escapes job ID",
			call: func(client *Client) error {
				resp, err := client.FileSearch.GetFileSearchV2("job/with space")
				if err != nil {
					return err
				}
				if resp == nil || resp.Data == nil || resp.Data.JobID != "job/with space" {
					t.Fatalf("unexpected status response: %#v", resp)
				}
				return nil
			},
			wantMethod:   http.MethodGet,
			wantPath:     "/service/v2/file-search/job%2Fwith%20space",
			wantQuery:    url.Values{},
			responseBody: `{"job_id":"job/with space","status":"completed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotMethod string
			var gotPath string
			var gotQuery url.Values
			var gotBody map[string]interface{}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotMethod = r.Method
				gotPath = r.URL.EscapedPath()
				gotQuery = r.URL.Query()

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

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client, err := NewClient("test-key", WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			if err := tt.call(client); err != nil {
				t.Fatalf("file-search call error = %v", err)
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
