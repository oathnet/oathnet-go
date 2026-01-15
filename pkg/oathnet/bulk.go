package oathnet

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

// BulkService handles bulk search operations.
type BulkService struct {
	client *Client
}

// BulkCreateOptions contains options for creating a bulk job.
type BulkCreateOptions struct {
	Format  string // jsonl, csv
	DBNames string
}

// Create creates a bulk search job.
func (s *BulkService) Create(terms []string, service string, opts *BulkCreateOptions) (*BulkJobResponse, error) {
	body := map[string]interface{}{
		"terms":   terms,
		"service": service,
		"format":  "jsonl",
	}

	if opts != nil {
		if opts.Format != "" {
			body["format"] = opts.Format
		}
		if opts.DBNames != "" {
			body["dbnames"] = opts.DBNames
		}
	}

	var rawResp map[string]interface{}
	err := s.client.post("/service/bulk-search", body, &rawResp)
	if err != nil {
		return nil, err
	}

	resp := &BulkJobResponse{Success: true}
	if _, ok := rawResp["success"]; ok {
		jsonData, _ := json.Marshal(rawResp)
		json.Unmarshal(jsonData, resp)
	} else {
		jsonData, _ := json.Marshal(rawResp)
		resp.Data = &BulkJobData{}
		json.Unmarshal(jsonData, resp.Data)
	}

	return resp, nil
}

// GetStatus gets bulk job status.
func (s *BulkService) GetStatus(jobID string) (*BulkJobResponse, error) {
	var rawResp map[string]interface{}
	err := s.client.get(fmt.Sprintf("/service/bulk-search/%s", jobID), nil, &rawResp)
	if err != nil {
		return nil, err
	}

	resp := &BulkJobResponse{Success: true}
	if _, ok := rawResp["success"]; ok {
		jsonData, _ := json.Marshal(rawResp)
		json.Unmarshal(jsonData, resp)
	} else {
		jsonData, _ := json.Marshal(rawResp)
		resp.Data = &BulkJobData{}
		json.Unmarshal(jsonData, resp.Data)
	}

	return resp, nil
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
	err := s.client.get("/service/bulk-search", params, &resp)
	return &resp, err
}

// Download downloads bulk search results.
func (s *BulkService) Download(jobID, outputPath string) error {
	data, err := s.client.getRaw(fmt.Sprintf("/service/bulk-search/%s/download", jobID))
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
