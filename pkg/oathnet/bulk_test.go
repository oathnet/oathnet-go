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
)

func TestBulkService_RequestConstruction(t *testing.T) {
	tests := []struct {
		name           string
		call           func(*Client) error
		responseStatus int
		responseBody   string
		contentType    string
		wantMethod     string
		wantPath       string
		wantQuery      url.Values
		wantBody       map[string]interface{}
	}{
		{
			name: "Create sends V2 body",
			call: func(client *Client) error {
				resp, err := client.Bulk.Create([]string{"alice@example.com", "bob@example.com"}, "breach", &BulkCreateOptions{
					Format:     "jsonl",
					DBNameList: []string{"linkedin.com", "twitter.com"},
					QueryConfig: BulkSearchQueryConfig{
						"wildcard":     true,
						"email_domain": []string{"gmail.com"},
					},
					Limit:  100,
					Fields: []string{"email", "password"},
				})
				if err != nil {
					return err
				}
				if resp == nil || !resp.Success || resp.Data == nil || resp.Data.JobID != "bulk-job-1" {
					t.Fatalf("unexpected create response: %#v", resp)
				}
				return nil
			},
			responseStatus: http.StatusAccepted,
			responseBody:   `{"success":true,"message":"created","data":{"job_id":"bulk-job-1","status":"queued","created_at":"2026-07-06T00:00:00Z","progress":{"records_done":0,"records_total":2,"bytes_done":0,"percent":0,"updated_at":"2026-07-06T00:00:00Z"},"request":{"type":"breach","service":"breach","format":"jsonl","limit":100,"fields":["email","password"],"request_count":2},"next_poll_after_ms":1000}}`,
			wantMethod:     http.MethodPost,
			wantPath:       "/service/v2/bulk-search",
			wantQuery:      url.Values{},
			wantBody: map[string]interface{}{
				"terms":   []interface{}{"alice@example.com", "bob@example.com"},
				"service": "breach",
				"format":  "jsonl",
				"dbnames": []interface{}{"linkedin.com", "twitter.com"},
				"query_config": map[string]interface{}{
					"wildcard":     true,
					"email_domain": []interface{}{"gmail.com"},
				},
				"limit":  float64(100),
				"fields": []interface{}{"email", "password"},
			},
		},
		{
			name: "Create converts legacy comma-separated DBNames to V2 array",
			call: func(client *Client) error {
				_, err := client.Bulk.Create([]string{"alice@example.com"}, "breach", &BulkCreateOptions{
					DBNames: "linkedin.com, twitter.com",
				})
				return err
			},
			responseStatus: http.StatusAccepted,
			responseBody:   `{"success":true,"data":{"job_id":"bulk-job-2","status":"queued","created_at":"2026-07-06T00:00:00Z"}}`,
			wantMethod:     http.MethodPost,
			wantPath:       "/service/v2/bulk-search",
			wantQuery:      url.Values{},
			wantBody: map[string]interface{}{
				"terms":   []interface{}{"alice@example.com"},
				"service": "breach",
				"dbnames": []interface{}{"linkedin.com", "twitter.com"},
			},
		},
		{
			name: "List sends pagination query",
			call: func(client *Client) error {
				resp, err := client.Bulk.List(2, 50)
				if err != nil {
					return err
				}
				if resp == nil || resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].JobID != "bulk-job-1" {
					t.Fatalf("unexpected list response: %#v", resp)
				}
				return nil
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"count":1,"next":null,"previous":"https://api.example.test/service/v2/bulk-search/list/?page=1&page_size=50","results":[{"id":"local-1","job_id":"bulk-job-1","status":"running","created_at":"2026-07-06T00:00:00Z","updated_at":"2026-07-06T00:00:01Z","search_service":"breach","output_format":"jsonl","results_expired":false,"query":"{\"service\":\"breach\"}","results_count":10,"lookups_deducted":1}]}`,
			wantMethod:     http.MethodGet,
			wantPath:       "/service/v2/bulk-search/list",
			wantQuery: url.Values{
				"page":      []string{"2"},
				"page_size": []string{"50"},
			},
		},
		{
			name: "GetStatus escapes job ID path segment",
			call: func(client *Client) error {
				resp, err := client.Bulk.GetStatus("job/with space")
				if err != nil {
					return err
				}
				if resp == nil || !resp.Success || resp.Data == nil || resp.Data.JobID != "job/with space" {
					t.Fatalf("unexpected status response: %#v", resp)
				}
				return nil
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"job_id":"job/with space","status":"completed","created_at":"2026-07-06T00:00:00Z","completed_at":"2026-07-06T00:00:02Z","progress":{"records_done":2,"records_total":2,"bytes_done":12,"percent":100,"updated_at":"2026-07-06T00:00:02Z"},"result":{"file_name":"bulk.jsonl","file_size":12,"records":2,"format":"jsonl","ready_at":"2026-07-06T00:00:02Z","expires_at":"2026-07-07T00:00:02Z"},"request":{"type":"breach","service":"breach","format":"jsonl","request_count":1}}`,
			wantMethod:     http.MethodGet,
			wantPath:       "/service/v2/bulk-search/job%2Fwith%20space",
			wantQuery:      url.Values{},
		},
		{
			name: "Download escapes job ID path segment and writes raw body",
			call: func(client *Client) error {
				outputPath := filepath.Join(t.TempDir(), "bulk.csv")
				err := client.Bulk.Download("job/with space", outputPath)
				if err != nil {
					return err
				}
				data, err := os.ReadFile(outputPath)
				if err != nil {
					return err
				}
				if string(data) != "email,password\nalice@example.com,secret\n" {
					t.Fatalf("downloaded data = %q", string(data))
				}
				return nil
			},
			responseStatus: http.StatusOK,
			responseBody:   "email,password\nalice@example.com,secret\n",
			contentType:    "text/csv",
			wantMethod:     http.MethodGet,
			wantPath:       "/service/v2/bulk-search/job%2Fwith%20space/download",
			wantQuery:      url.Values{},
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
					if len(body) == 0 {
						t.Fatal("request body is empty")
					}
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
				w.WriteHeader(tt.responseStatus)
				if tt.responseBody != "" {
					_, _ = w.Write([]byte(tt.responseBody))
				}
			}))
			defer server.Close()

			client, err := NewClient("test-key", WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			if err := tt.call(client); err != nil {
				t.Fatalf("bulk call error = %v", err)
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
