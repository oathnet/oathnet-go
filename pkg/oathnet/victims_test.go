package oathnet

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
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

func TestVictimsService_RawDetailRequestConstruction(t *testing.T) {
	tests := []struct {
		name         string
		call         func(*Client) error
		responseBody string
		contentType  string
		wantPath     string
		wantQuery    url.Values
	}{
		{
			name: "GetVictimManifestV2 escapes log ID and sends search_id",
			call: func(client *Client) error {
				resp, err := client.Victims.GetVictimManifestV2("log/with space", &VictimRawOptions{
					SearchID: "session-123",
				})
				if err != nil {
					return err
				}
				if resp == nil || resp.LogID != "log/with space" {
					t.Fatalf("unexpected manifest response: %#v", resp)
				}
				return nil
			},
			responseBody: `{"log_id":"log/with space","victim_tree":{"id":"root","name":"root","type":"directory"}}`,
			contentType:  "application/json",
			wantPath:     "/service/v2/victims/log%2Fwith%20space",
			wantQuery:    url.Values{"search_id": []string{"session-123"}},
		},
		{
			name: "GetVictimFileV2 escapes log and file IDs",
			call: func(client *Client) error {
				data, err := client.Victims.GetVictimFileV2(
					"log/with space",
					"Browser Passwords/file 1.txt",
					&VictimRawOptions{SearchID: "session-123"},
				)
				if err != nil {
					return err
				}
				if string(data) != "raw file content" {
					t.Fatalf("file content = %q", string(data))
				}
				return nil
			},
			responseBody: "raw file content",
			contentType:  "text/plain",
			wantPath:     "/service/v2/victims/log%2Fwith%20space/files/Browser%20Passwords%2Ffile%201.txt",
			wantQuery:    url.Values{"search_id": []string{"session-123"}},
		},
		{
			name: "DownloadVictimArchiveV2 returns bytes and DownloadArchive writes file",
			call: func(client *Client) error {
				data, err := client.Victims.DownloadVictimArchiveV2(
					"log/with space",
					&VictimRawOptions{SearchID: "session-123"},
				)
				if err != nil {
					return err
				}
				if string(data) != "zip bytes" {
					t.Fatalf("archive bytes = %q", string(data))
				}
				outputPath := filepath.Join(t.TempDir(), "victim.zip")
				err = client.Victims.DownloadArchive(
					"log/with space",
					outputPath,
					&VictimRawOptions{SearchID: "session-123"},
				)
				if err != nil {
					return err
				}
				saved, err := os.ReadFile(outputPath)
				if err != nil {
					return err
				}
				if string(saved) != string(data) {
					t.Fatalf("saved archive = %q, want %q", string(saved), string(data))
				}
				return nil
			},
			responseBody: "zip bytes",
			contentType:  "application/zip",
			wantPath:     "/service/v2/victims/log%2Fwith%20space/archive",
			wantQuery:    url.Values{"search_id": []string{"session-123"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath string
			var gotQuery url.Values
			var gotAPIKey string

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.EscapedPath()
				gotQuery = r.URL.Query()
				gotAPIKey = r.Header.Get("x-api-key")
				_, _ = io.ReadAll(r.Body)
				w.Header().Set("Content-Type", tt.contentType)
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client, err := NewClient("test-key", WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			if err := tt.call(client); err != nil {
				t.Fatalf("victims call error = %v", err)
			}
			if gotAPIKey != "test-key" {
				t.Fatalf("x-api-key = %q, want test-key", gotAPIKey)
			}
			if gotPath != tt.wantPath {
				t.Fatalf("path = %s, want %s", gotPath, tt.wantPath)
			}
			if !reflect.DeepEqual(gotQuery, tt.wantQuery) {
				t.Fatalf("query = %v, want %v", gotQuery, tt.wantQuery)
			}
		})
	}
}
