package oathnet

import (
	"fmt"
	"net/url"
	"strconv"
)

// ScannerService handles scanner and notification management operations.
type ScannerService struct {
	client *Client
}

// GetQuota gets the current scanner quota.
func (s *ScannerService) GetQuota() (*ScannerQuotaResponse, error) {
	var resp ScannerQuotaResponse
	err := s.client.get("/scanners/quota", nil, &resp)
	return &resp, err
}

// List lists scanners, optionally filtered by status or scanner type.
func (s *ScannerService) List(opts *ScannerListOptions) ([]Scanner, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Status != "" {
			params.Set("status", opts.Status)
		}
		if opts.ScannerType != "" {
			params.Set("scanner_type", opts.ScannerType)
		}
	}

	var resp []Scanner
	err := s.client.get("/scanners", params, &resp)
	return resp, err
}

// Create creates a scanner.
func (s *ScannerService) Create(req ScannerCreateRequest) (*Scanner, error) {
	var resp Scanner
	err := s.client.post("/scanners/create", req, &resp)
	return &resp, err
}

// Get gets a scanner by UID.
func (s *ScannerService) Get(scannerUID string) (*Scanner, error) {
	var resp Scanner
	err := s.client.get(scannerPath(scannerUID, ""), nil, &resp)
	return &resp, err
}

// Update partially updates a scanner. Use Pause or Resume for status changes.
func (s *ScannerService) Update(scannerUID string, req ScannerUpdateRequest) (*Scanner, error) {
	var resp Scanner
	err := s.client.patch(scannerPath(scannerUID, "/update"), req, &resp)
	return &resp, err
}

// Delete deletes a scanner. A successful delete returns HTTP 204 with no body.
func (s *ScannerService) Delete(scannerUID string) error {
	return s.client.delete(scannerPath(scannerUID, "/delete"), nil)
}

// TestDraftDelivery tests delivery settings before saving a scanner.
func (s *ScannerService) TestDraftDelivery(req ScannerDraftTestRequest) (*ScannerTestResponse, error) {
	var resp ScannerTestResponse
	err := s.client.post("/scanners/test-delivery", req, &resp)
	return &resp, err
}

// TestNotification sends a test notification for an existing scanner.
func (s *ScannerService) TestNotification(scannerUID string) (*ScannerTestResponse, error) {
	var resp ScannerTestResponse
	err := s.client.post(scannerPath(scannerUID, "/test"), map[string]interface{}{}, &resp)
	return &resp, err
}

// GetWebhookSecurity gets webhook verification metadata for a scanner.
func (s *ScannerService) GetWebhookSecurity(scannerUID string) (*ScannerWebhookSecurityResponse, error) {
	var resp ScannerWebhookSecurityResponse
	err := s.client.get(scannerPath(scannerUID, "/webhook-security"), nil, &resp)
	return &resp, err
}

// RotateWebhookSecret rotates a scanner webhook secret.
func (s *ScannerService) RotateWebhookSecret(scannerUID string) (*ScannerWebhookRotateResponse, error) {
	var resp ScannerWebhookRotateResponse
	err := s.client.post(scannerPath(scannerUID, "/webhook-security/rotate"), map[string]interface{}{}, &resp)
	return &resp, err
}

// Pause pauses scheduling for a scanner.
func (s *ScannerService) Pause(scannerUID string) (*Scanner, error) {
	var resp Scanner
	err := s.client.post(scannerPath(scannerUID, "/pause"), map[string]interface{}{}, &resp)
	return &resp, err
}

// Resume resumes a paused or disabled scanner.
func (s *ScannerService) Resume(scannerUID string) (*Scanner, error) {
	var resp Scanner
	err := s.client.post(scannerPath(scannerUID, "/resume"), map[string]interface{}{}, &resp)
	return &resp, err
}

// Trigger queues an immediate run for a scanner.
func (s *ScannerService) Trigger(scannerUID string) (*ScannerTriggerResponse, error) {
	var resp ScannerTriggerResponse
	err := s.client.post(scannerPath(scannerUID, "/trigger"), map[string]interface{}{}, &resp)
	return &resp, err
}

// ListRuns lists scanner runs, optionally filtered by status and limit.
func (s *ScannerService) ListRuns(scannerUID string, opts *ScannerRunsOptions) ([]ScannerRun, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Status != "" {
			params.Set("status", opts.Status)
		}
		if opts.Limit > 0 {
			params.Set("limit", strconv.Itoa(opts.Limit))
		}
	}

	var resp []ScannerRun
	err := s.client.get(scannerPath(scannerUID, "/runs"), params, &resp)
	return resp, err
}

// GetRun gets a scanner run and its notification attempts.
func (s *ScannerService) GetRun(scannerUID, runUID string) (*ScannerRunDetail, error) {
	var resp ScannerRunDetail
	path := fmt.Sprintf("%s/%s", scannerPath(scannerUID, "/runs"), url.PathEscape(runUID))
	err := s.client.get(path, nil, &resp)
	return &resp, err
}

func scannerPath(scannerUID, suffix string) string {
	return fmt.Sprintf("/scanners/%s%s", url.PathEscape(scannerUID), suffix)
}

// ScannerQueryConfig is the persistent search definition stored on a scanner.
// It accepts the same flexible JSON fields documented for scanner query_config.
type ScannerQueryConfig map[string]interface{}

// ScannerListOptions contains filters for listing scanners.
type ScannerListOptions struct {
	Status      string
	ScannerType string
}

// ScannerRunsOptions contains filters for listing scanner runs.
type ScannerRunsOptions struct {
	Status string
	Limit  int
}

// ScannerCreateRequest contains scanner creation fields.
type ScannerCreateRequest struct {
	Name                string             `json:"name"`
	ScannerType         string             `json:"scanner_type"`
	QueryConfig         ScannerQueryConfig `json:"query_config"`
	NotificationType    string             `json:"notification_type"`
	WebhookURL          string             `json:"webhook_url,omitempty"`
	WebhookSecret       string             `json:"webhook_secret,omitempty"`
	WebhookSecurityMode string             `json:"webhook_security_mode,omitempty"`
	NotifyOnZeroResults *bool              `json:"notify_on_zero_results,omitempty"`
}

// ScannerUpdateRequest contains partial scanner update fields.
// It intentionally has no status field; use Pause or Resume for state changes.
type ScannerUpdateRequest struct {
	Name                string             `json:"name,omitempty"`
	ScannerType         string             `json:"scanner_type,omitempty"`
	QueryConfig         ScannerQueryConfig `json:"query_config,omitempty"`
	NotificationType    string             `json:"notification_type,omitempty"`
	WebhookURL          *string            `json:"webhook_url,omitempty"`
	WebhookSecret       string             `json:"webhook_secret,omitempty"`
	WebhookSecurityMode string             `json:"webhook_security_mode,omitempty"`
	NotifyOnZeroResults *bool              `json:"notify_on_zero_results,omitempty"`
}

// ScannerDraftTestRequest contains delivery-test fields for an unsaved scanner.
type ScannerDraftTestRequest struct {
	Name                string             `json:"name,omitempty"`
	ScannerType         string             `json:"scanner_type"`
	QueryConfig         ScannerQueryConfig `json:"query_config"`
	NotificationType    string             `json:"notification_type"`
	WebhookURL          string             `json:"webhook_url,omitempty"`
	WebhookSecret       string             `json:"webhook_secret,omitempty"`
	WebhookSecurityMode string             `json:"webhook_security_mode,omitempty"`
	NotifyOnZeroResults *bool              `json:"notify_on_zero_results,omitempty"`
}

// ScannerQuotaResponse is the raw scanner quota payload.
type ScannerQuotaResponse struct {
	MaxScanners  int  `json:"max_scanners"`
	CurrentCount int  `json:"current_count"`
	Remaining    int  `json:"remaining"`
	CanCreate    bool `json:"can_create"`
}

// Scanner is a saved scanner configuration and status snapshot.
type Scanner struct {
	UID                        string             `json:"uid"`
	User                       string             `json:"user,omitempty"`
	Name                       string             `json:"name"`
	ScannerType                string             `json:"scanner_type"`
	ScannerTypeDisplay         string             `json:"scanner_type_display,omitempty"`
	Status                     string             `json:"status"`
	StatusDisplay              string             `json:"status_display,omitempty"`
	QueryConfig                ScannerQueryConfig `json:"query_config"`
	NotificationType           string             `json:"notification_type"`
	NotificationTypeDisplay    string             `json:"notification_type_display,omitempty"`
	WebhookURL                 *string            `json:"webhook_url,omitempty"`
	WebhookSecurityMode        string             `json:"webhook_security_mode,omitempty"`
	WebhookSecretLastRotatedAt *string            `json:"webhook_secret_last_rotated_at,omitempty"`
	NotifyOnZeroResults        bool               `json:"notify_on_zero_results"`
	CreatedAt                  string             `json:"created_at,omitempty"`
	UpdatedAt                  string             `json:"updated_at,omitempty"`
	LastRunAt                  *string            `json:"last_run_at,omitempty"`
	LastFoundAt                *string            `json:"last_found_at,omitempty"`
	NextRunAt                  *string            `json:"next_run_at,omitempty"`
	TotalResultsFound          int                `json:"total_results_found"`
	TotalRuns                  int                `json:"total_runs"`
	ConsecutiveFailures        int                `json:"consecutive_failures"`
}

// ScannerTestPreview describes the payload that would be delivered.
type ScannerTestPreview struct {
	EventName           string                 `json:"event_name,omitempty"`
	SearchURL           string                 `json:"search_url,omitempty"`
	SecurityMode        string                 `json:"security_mode,omitempty"`
	VerificationMethod  string                 `json:"verification_method,omitempty"`
	AuthorizationScheme string                 `json:"authorization_scheme,omitempty"`
	EncryptionEnabled   bool                   `json:"encryption_enabled,omitempty"`
	HeaderNames         []string               `json:"header_names,omitempty"`
	BodyKind            string                 `json:"body_kind,omitempty"`
	BodyExample         map[string]interface{} `json:"body_example,omitempty"`
}

// ScannerTestResponse is returned for draft and saved scanner delivery tests.
// A 200 response can contain Success=false when delivery failed; that is data,
// not a transport error.
type ScannerTestResponse struct {
	Success                bool                `json:"success"`
	NotificationType       string              `json:"notification_type,omitempty"`
	Target                 string              `json:"target,omitempty"`
	StatusCode             *int                `json:"status_code,omitempty"`
	Message                string              `json:"message,omitempty"`
	Preview                *ScannerTestPreview `json:"preview,omitempty"`
	GeneratedWebhookSecret string              `json:"generated_webhook_secret,omitempty"`
}

// ScannerTriggerResponse acknowledges an immediate scanner run request.
type ScannerTriggerResponse struct {
	Message    string `json:"message,omitempty"`
	ScannerUID string `json:"scanner_uid,omitempty"`
}

// ScannerRun is a scanner run summary.
type ScannerRun struct {
	UID                string                 `json:"uid"`
	Status             string                 `json:"status"`
	StatusDisplay      string                 `json:"status_display,omitempty"`
	StartedAt          *string                `json:"started_at,omitempty"`
	CompletedAt        *string                `json:"completed_at,omitempty"`
	DurationSeconds    *float64               `json:"duration_seconds,omitempty"`
	SearchFrom         string                 `json:"search_from,omitempty"`
	SearchTo           string                 `json:"search_to,omitempty"`
	ResultsCount       int                    `json:"results_count"`
	ResultsSample      []interface{}          `json:"results_sample,omitempty"`
	RequestParamsUsed  map[string]interface{} `json:"request_params_used,omitempty"`
	ExecutionMeta      map[string]interface{} `json:"execution_meta,omitempty"`
	NotificationSent   bool                   `json:"notification_sent"`
	NotificationSentAt *string                `json:"notification_sent_at,omitempty"`
	NotificationError  *string                `json:"notification_error,omitempty"`
	ErrorMessage       *string                `json:"error_message,omitempty"`
}

// ScannerNotificationAttempt is a delivery attempt for a scanner run.
type ScannerNotificationAttempt struct {
	UID                string  `json:"uid"`
	AttemptNumber      int     `json:"attempt_number"`
	AttemptedAt        string  `json:"attempted_at,omitempty"`
	Success            bool    `json:"success"`
	ResponseStatusCode *int    `json:"response_status_code,omitempty"`
	ErrorMessage       string  `json:"error_message,omitempty"`
	ResponseExcerpt    *string `json:"response_excerpt,omitempty"`
	DeliveryID         string  `json:"delivery_id,omitempty"`
	SecurityMode       string  `json:"security_mode,omitempty"`
	Target             string  `json:"target,omitempty"`
}

// ScannerRunDetail is a scanner run plus detail-only fields.
type ScannerRunDetail struct {
	ScannerRun
	SearchURL            string                       `json:"search_url,omitempty"`
	PreviousResultsCount *int                         `json:"previous_results_count,omitempty"`
	ResultsDelta         *int                         `json:"results_delta,omitempty"`
	NotificationLogs     []ScannerNotificationAttempt `json:"notification_logs,omitempty"`
}

// ScannerWebhookSecurityResponse wraps webhook security metadata.
type ScannerWebhookSecurityResponse struct {
	Success bool                        `json:"success"`
	Message string                      `json:"message,omitempty"`
	Data    *ScannerWebhookSecurityData `json:"data,omitempty"`
}

type ScannerWebhookSecurityData struct {
	ScannerUID          string   `json:"scanner_uid,omitempty"`
	NotificationType    string   `json:"notification_type,omitempty"`
	WebhookSecurityMode string   `json:"webhook_security_mode,omitempty"`
	VerificationMethod  string   `json:"verification_method,omitempty"`
	EncryptionEnabled   bool     `json:"encryption_enabled,omitempty"`
	SecretConfigured    bool     `json:"secret_configured,omitempty"`
	SecretPreview       *string  `json:"secret_preview,omitempty"`
	KeyID               *string  `json:"key_id,omitempty"`
	LastRotatedAt       *string  `json:"last_rotated_at,omitempty"`
	Headers             []string `json:"headers,omitempty"`
}

// ScannerWebhookRotateResponse wraps rotated webhook secret data.
type ScannerWebhookRotateResponse struct {
	Success bool                      `json:"success"`
	Message string                    `json:"message,omitempty"`
	Data    *ScannerWebhookRotateData `json:"data,omitempty"`
}

type ScannerWebhookRotateData struct {
	ScannerUID          string  `json:"scanner_uid,omitempty"`
	WebhookSecurityMode string  `json:"webhook_security_mode,omitempty"`
	Secret              string  `json:"secret,omitempty"`
	SecretPreview       *string `json:"secret_preview,omitempty"`
	KeyID               *string `json:"key_id,omitempty"`
	LastRotatedAt       *string `json:"last_rotated_at,omitempty"`
}
