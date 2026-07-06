package oathnet

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const bulkSearchBasePath = "/service/v2/bulk-search"

// BulkService handles bulk search operations.
type BulkService struct {
	client *Client
}

// BulkSearchQueryConfig is the shared V2 search/filter configuration applied
// to every term in a bulk search job.
type BulkSearchQueryConfig map[string]interface{}

// BulkSearchCreateRequest is the canonical V2 bulk-search request body.
type BulkSearchCreateRequest struct {
	Terms       []string              `json:"terms,omitempty"`
	Service     string                `json:"service"`
	Format      string                `json:"format,omitempty"`
	DBNames     []string              `json:"dbnames,omitempty"`
	QueryConfig BulkSearchQueryConfig `json:"query_config,omitempty"`
	Limit       int                   `json:"limit,omitempty"`
	Fields      []string              `json:"fields,omitempty"`
}

// BulkCreateOptions contains options for creating a bulk job.
type BulkCreateOptions struct {
	Format      string // csv, json, jsonl, txt, html
	DBNames     string // Deprecated: comma-separated legacy database filter. Use DBNameList.
	DBNameList  []string
	QueryConfig BulkSearchQueryConfig
	Limit       int
	Fields      []string
}

// Create creates a bulk search job.
func (s *BulkService) Create(terms []string, service string, opts *BulkCreateOptions) (*BulkJobResponse, error) {
	body := BulkSearchCreateRequest{
		Terms:   terms,
		Service: service,
	}

	if opts != nil {
		if opts.Format != "" {
			body.Format = opts.Format
		}
		if len(opts.DBNameList) > 0 {
			body.DBNames = append([]string(nil), opts.DBNameList...)
		} else if opts.DBNames != "" {
			body.DBNames = splitCommaList(opts.DBNames)
		}
		if opts.QueryConfig != nil {
			body.QueryConfig = opts.QueryConfig
		}
		if opts.Limit > 0 {
			body.Limit = opts.Limit
		}
		if opts.Fields != nil {
			body.Fields = opts.Fields
		}
	}

	var rawResp map[string]interface{}
	err := s.client.post(bulkSearchBasePath, body, &rawResp)
	if err != nil {
		return nil, err
	}

	return decodeBulkJobResponse(rawResp)
}

// GetStatus gets bulk job status.
func (s *BulkService) GetStatus(jobID string) (*BulkJobResponse, error) {
	var rawResp map[string]interface{}
	err := s.client.get(bulkSearchJobPath(jobID, ""), nil, &rawResp)
	if err != nil {
		return nil, err
	}

	return decodeBulkJobResponse(rawResp)
}

// List lists bulk search jobs.
func (s *BulkService) List(page, pageSize int) (*BulkJobListResponse, error) {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if pageSize > 0 {
		params.Set("page_size", strconv.Itoa(pageSize))
	}

	var resp BulkJobListResponse
	err := s.client.get(bulkSearchBasePath+"/list", params, &resp)
	return &resp, err
}

// Download downloads bulk search results.
func (s *BulkService) Download(jobID, outputPath string) error {
	data, err := s.client.getRaw(bulkSearchJobPath(jobID, "/download"))
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, data, 0644)
}

// WaitForCompletion waits for a bulk job to complete.
func (s *BulkService) WaitForCompletion(jobID string, pollInterval, timeout time.Duration) (*BulkJobResponse, error) {
	startTime := time.Now()

	for {
		resp, err := s.GetStatus(jobID)
		if err != nil {
			return nil, err
		}

		if resp.Data != nil {
			status := resp.Data.Status
			if status == "completed" || status == "canceled" {
				return resp, nil
			}
		}

		elapsed := time.Since(startTime)
		if elapsed >= timeout {
			return nil, fmt.Errorf("bulk job %s did not complete within %v", jobID, timeout)
		}

		sleepTime := pollInterval
		if resp.Data != nil && resp.Data.NextPollAfterMs > 0 {
			suggested := time.Duration(resp.Data.NextPollAfterMs) * time.Millisecond
			if suggested < sleepTime {
				sleepTime = suggested
			}
		}
		time.Sleep(sleepTime)
	}
}

// Search creates a bulk search, waits for completion, and downloads.
func (s *BulkService) Search(terms []string, service, outputPath string, opts *BulkCreateOptions, timeout time.Duration) error {
	job, err := s.Create(terms, service, opts)
	if err != nil {
		return err
	}

	if job.Data == nil || job.Data.JobID == "" {
		return fmt.Errorf("failed to create bulk job")
	}

	_, err = s.WaitForCompletion(job.Data.JobID, 5*time.Second, timeout)
	if err != nil {
		return err
	}

	return s.Download(job.Data.JobID, outputPath)
}

func bulkSearchJobPath(jobID, suffix string) string {
	return fmt.Sprintf("%s/%s%s", bulkSearchBasePath, url.PathEscape(jobID), suffix)
}

func decodeBulkJobResponse(rawResp map[string]interface{}) (*BulkJobResponse, error) {
	resp := &BulkJobResponse{Success: true}
	jsonData, err := json.Marshal(rawResp)
	if err != nil {
		return nil, err
	}

	if _, ok := rawResp["success"]; ok {
		if err := json.Unmarshal(jsonData, resp); err != nil {
			return nil, err
		}
		return resp, nil
	}

	resp.Data = &BulkJobData{}
	if err := json.Unmarshal(jsonData, resp.Data); err != nil {
		return nil, err
	}
	return resp, nil
}

func splitCommaList(value string) []string {
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			items = append(items, item)
		}
	}
	return items
}
