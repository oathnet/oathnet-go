// Package oathnet provides a Go SDK for the OathNet API.
package oathnet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is the main OathNet API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client

	// Services
	Search     *SearchService
	OSINT      *OSINTService
	Stealer    *StealerV2Service
	Victims    *VictimsService
	FileSearch *FileSearchService
	Exports    *ExportsService
	Bulk       *BulkService
	Utility    *UtilityService
}

// ClientOption is a function that configures the client.
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithTimeout sets a custom HTTP timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// NewClient creates a new OathNet API client.
func NewClient(apiKey string, opts ...ClientOption) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	c := &Client{
		apiKey:  apiKey,
		baseURL: "https://oathnet.org/api",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	// Initialize services
	c.Search = &SearchService{client: c}
	c.OSINT = &OSINTService{client: c}
	c.Stealer = &StealerV2Service{client: c}
	c.Victims = &VictimsService{client: c}
	c.FileSearch = &FileSearchService{client: c}
	c.Exports = &ExportsService{client: c}
	c.Bulk = &BulkService{client: c}
	c.Utility = &UtilityService{client: c}

	return c, nil
}

// get performs a GET request.
func (c *Client) get(path string, params url.Values, result interface{}) error {
	reqURL := c.baseURL + path
	if len(params) > 0 {
		reqURL += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.handleResponse(resp, result)
}

// post performs a POST request.
func (c *Client) post(path string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest("POST", c.baseURL+path, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.handleResponse(resp, result)
}

// getRaw performs a GET request and returns raw bytes.
func (c *Client) getRaw(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, c.parseError(resp.StatusCode, body)
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) handleResponse(resp *http.Response, result interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return c.parseError(resp.StatusCode, body)
	}

	if result != nil {
		return json.Unmarshal(body, result)
	}

	return nil
}

func (c *Client) parseError(statusCode int, body []byte) error {
	var errResp struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	json.Unmarshal(body, &errResp)

	message := errResp.Message
	if message == "" {
		message = errResp.Error
	}
	if message == "" {
		message = string(body)
	}

	switch statusCode {
	case 401:
		return &AuthenticationError{Message: message}
	case 400:
		// Check for auth-related messages
		if containsAuthMessage(message) {
			return &AuthenticationError{Message: message}
		}
		return &ValidationError{Message: message}
	case 404:
		return &NotFoundError{Message: message}
	case 429:
		return &RateLimitError{Message: message}
	default:
		return &OathNetError{Message: message, StatusCode: statusCode}
	}
}

func containsAuthMessage(msg string) bool {
	return contains(msg, "credentials") || contains(msg, "api key") || contains(msg, "invalid api key")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsLower(s, substr))
}

func containsLower(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if toLower(s[i:i+len(substr)]) == toLower(substr) {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}
