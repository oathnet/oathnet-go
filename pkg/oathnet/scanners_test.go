package oathnet

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestScannerService_RequestConstruction(t *testing.T) {
	notifyOnZero := false
	webhookURL := "https://alerts.example.com/oathnet"
	queryConfig := ScannerQueryConfig{
		"q":        "example.com",
		"wildcard": true,
	}

	tests := []struct {
		name           string
		call           func(*Client) error
		responseStatus int
		responseBody   string
		wantMethod     string
		wantPath       string
		wantQuery      url.Values
		wantBody       map[string]interface{}
	}{
		{
			name: "GetQuota",
			call: func(client *Client) error {
				_, err := client.Scanners.GetQuota()
				return err
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"max_scanners":5,"current_count":1,"remaining":4,"can_create":true}`,
			wantMethod:     http.MethodGet,
			wantPath:       "/scanners/quota",
			wantQuery:      url.Values{},
		},
		{
			name: "List",
			call: func(client *Client) error {
				_, err := client.Scanners.List(&ScannerListOptions{Status: "active", ScannerType: "stealer"})
				return err
			},
			responseStatus: http.StatusOK,
			responseBody:   `[]`,
			wantMethod:     http.MethodGet,
			wantPath:       "/scanners",
			wantQuery: url.Values{
				"scanner_type": []string{"stealer"},
				"status":       []string{"active"},
			},
		},
		{
			name: "Create",
			call: func(client *Client) error {
				_, err := client.Scanners.Create(ScannerCreateRequest{
					Name:                "Example scanner",
					ScannerType:         "stealer",
					QueryConfig:         queryConfig,
					NotificationType:    "webhook",
					WebhookURL:          webhookURL,
					WebhookSecurityMode: "signed_json",
					NotifyOnZeroResults: &notifyOnZero,
				})
				return err
			},
			responseStatus: http.StatusCreated,
			responseBody:   scannerResponseJSON,
			wantMethod:     http.MethodPost,
			wantPath:       "/scanners/create",
			wantQuery:      url.Values{},
			wantBody: map[string]interface{}{
				"name":                   "Example scanner",
				"scanner_type":           "stealer",
				"query_config":           map[string]interface{}{"q": "example.com", "wildcard": true},
				"notification_type":      "webhook",
				"webhook_url":            webhookURL,
				"webhook_security_mode":  "signed_json",
				"notify_on_zero_results": false,
			},
		},
		{
			name: "Get escapes scanner UID",
			call: func(client *Client) error {
				_, err := client.Scanners.Get("scanner/with space")
				return err
			},
			responseStatus: http.StatusOK,
			responseBody:   scannerResponseJSON,
			wantMethod:     http.MethodGet,
			wantPath:       "/scanners/scanner%2Fwith%20space",
			wantQuery:      url.Values{},
		},
		{
			name: "Update omits status and sends explicit false",
			call: func(client *Client) error {
				_, err := client.Scanners.Update("scanner-123", ScannerUpdateRequest{
					Name:                "Renamed scanner",
					QueryConfig:         queryConfig,
					NotificationType:    "webhook",
					WebhookURL:          &webhookURL,
					WebhookSecurityMode: "signed_json",
					NotifyOnZeroResults: &notifyOnZero,
				})
				return err
			},
			responseStatus: http.StatusOK,
			responseBody:   scannerResponseJSON,
			wantMethod:     http.MethodPatch,
			wantPath:       "/scanners/scanner-123/update",
			wantQuery:      url.Values{},
			wantBody: map[string]interface{}{
				"name":                   "Renamed scanner",
				"query_config":           map[string]interface{}{"q": "example.com", "wildcard": true},
				"notification_type":      "webhook",
				"webhook_url":            webhookURL,
				"webhook_security_mode":  "signed_json",
				"notify_on_zero_results": false,
			},
		},
		{
			name: "Delete accepts 204 no body",
			call: func(client *Client) error {
				return client.Scanners.Delete("scanner-123")
			},
			responseStatus: http.StatusNoContent,
			wantMethod:     http.MethodDelete,
			wantPath:       "/scanners/scanner-123/delete",
			wantQuery:      url.Values{},
		},
		{
			name: "TestDraftDelivery returns success false as data",
			call: func(client *Client) error {
				resp, err := client.Scanners.TestDraftDelivery(ScannerDraftTestRequest{
					ScannerType:         "stealer",
					QueryConfig:         queryConfig,
					NotificationType:    "webhook",
					WebhookURL:          webhookURL,
					WebhookSecurityMode: "signed_json",
					NotifyOnZeroResults: &notifyOnZero,
				})
				if err != nil {
					return err
				}
				if resp.Success {
					t.Fatal("Success = true, want false")
				}
				return nil
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"success":false,"notification_type":"webhook","target":"https://alerts.example.com/oathnet","status_code":500,"message":"delivery failed"}`,
			wantMethod:     http.MethodPost,
			wantPath:       "/scanners/test-delivery",
			wantQuery:      url.Values{},
			wantBody: map[string]interface{}{
				"scanner_type":           "stealer",
				"query_config":           map[string]interface{}{"q": "example.com", "wildcard": true},
				"notification_type":      "webhook",
				"webhook_url":            webhookURL,
				"webhook_security_mode":  "signed_json",
				"notify_on_zero_results": false,
			},
		},
		{
			name: "TestNotification",
			call: func(client *Client) error {
				_, err := client.Scanners.TestNotification("scanner-123")
				return err
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"success":true,"message":"sent"}`,
			wantMethod:     http.MethodPost,
			wantPath:       "/scanners/scanner-123/test",
			wantQuery:      url.Values{},
			wantBody:       map[string]interface{}{},
		},
		{
			name: "GetWebhookSecurity",
			call: func(client *Client) error {
				_, err := client.Scanners.GetWebhookSecurity("scanner-123")
				return err
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"success":true,"data":{"scanner_uid":"scanner-123","webhook_security_mode":"signed_json","secret_configured":true}}`,
			wantMethod:     http.MethodGet,
			wantPath:       "/scanners/scanner-123/webhook-security",
			wantQuery:      url.Values{},
		},
		{
			name: "RotateWebhookSecret",
			call: func(client *Client) error {
				_, err := client.Scanners.RotateWebhookSecret("scanner-123")
				return err
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"success":true,"data":{"scanner_uid":"scanner-123","secret":"secret-value"}}`,
			wantMethod:     http.MethodPost,
			wantPath:       "/scanners/scanner-123/webhook-security/rotate",
			wantQuery:      url.Values{},
			wantBody:       map[string]interface{}{},
		},
		{
			name: "Pause",
			call: func(client *Client) error {
				_, err := client.Scanners.Pause("scanner-123")
				return err
			},
			responseStatus: http.StatusOK,
			responseBody:   scannerResponseJSON,
			wantMethod:     http.MethodPost,
			wantPath:       "/scanners/scanner-123/pause",
			wantQuery:      url.Values{},
			wantBody:       map[string]interface{}{},
		},
		{
			name: "Resume",
			call: func(client *Client) error {
				_, err := client.Scanners.Resume("scanner-123")
				return err
			},
			responseStatus: http.StatusOK,
			responseBody:   scannerResponseJSON,
			wantMethod:     http.MethodPost,
			wantPath:       "/scanners/scanner-123/resume",
			wantQuery:      url.Values{},
			wantBody:       map[string]interface{}{},
		},
		{
			name: "Trigger",
			call: func(client *Client) error {
				_, err := client.Scanners.Trigger("scanner-123")
				return err
			},
			responseStatus: http.StatusAccepted,
			responseBody:   `{"message":"queued","scanner_uid":"scanner-123"}`,
			wantMethod:     http.MethodPost,
			wantPath:       "/scanners/scanner-123/trigger",
			wantQuery:      url.Values{},
			wantBody:       map[string]interface{}{},
		},
		{
			name: "ListRuns",
			call: func(client *Client) error {
				_, err := client.Scanners.ListRuns("scanner-123", &ScannerRunsOptions{Status: "completed", Limit: 50})
				return err
			},
			responseStatus: http.StatusOK,
			responseBody:   `[]`,
			wantMethod:     http.MethodGet,
			wantPath:       "/scanners/scanner-123/runs",
			wantQuery: url.Values{
				"limit":  []string{"50"},
				"status": []string{"completed"},
			},
		},
		{
			name: "GetRun",
			call: func(client *Client) error {
				_, err := client.Scanners.GetRun("scanner-123", "run-456")
				return err
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"uid":"run-456","status":"completed","results_count":1,"notification_sent":true,"previous_results_count":null,"results_delta":1,"notification_logs":[]}`,
			wantMethod:     http.MethodGet,
			wantPath:       "/scanners/scanner-123/runs/run-456",
			wantQuery:      url.Values{},
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
					if len(body) == 0 {
						t.Fatal("request body is empty")
					}
					if err := json.Unmarshal(body, &gotBody); err != nil {
						t.Fatalf("Unmarshal request body error = %v; body=%s", err, string(body))
					}
				} else if len(body) > 0 {
					t.Fatalf("request body = %s, want empty", string(body))
				}

				w.Header().Set("Content-Type", "application/json")
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
				t.Fatalf("scanner call error = %v", err)
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

const scannerResponseJSON = `{
	"uid":"scanner-123",
	"user":"user-123",
	"name":"Example scanner",
	"scanner_type":"stealer",
	"scanner_type_display":"Stealer",
	"status":"active",
	"status_display":"Active",
	"query_config":{"q":"example.com"},
	"notification_type":"webhook",
	"notification_type_display":"Webhook",
	"webhook_url":"https://alerts.example.com/oathnet",
	"webhook_security_mode":"signed_json",
	"notify_on_zero_results":false,
	"created_at":"2026-07-06T00:00:00Z",
	"updated_at":"2026-07-06T00:00:00Z",
	"last_run_at":null,
	"last_found_at":null,
	"next_run_at":null,
	"total_results_found":0,
	"total_runs":0,
	"consecutive_failures":0
}`
