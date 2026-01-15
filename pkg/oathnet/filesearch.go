package oathnet

import (
	"encoding/json"
	"fmt"
	"time"
)

// FileSearchService handles V2 file search operations.
type FileSearchService struct {
	client *Client
}

// FileSearchCreateOptions contains options for creating a file search job.
type FileSearchCreateOptions struct {
	SearchMode     string   // literal, regex, wildcard
	LogIDs         []string
	IncludeMatches bool
	CaseSensitive  bool
	ContextLines   int
	FilePattern    string
	MaxMatches     int
}

// Create creates a file search job.
func (s *FileSearchService) Create(expression string, opts *FileSearchCreateOptions) (*FileSearchJobResponse, error) {
	body := map[string]interface{}{
		"expression":      expression,
		"search_mode":     "literal",
		"include_matches": true,
		"case_sensitive":  false,
		"context_lines":   2,
		"max_matches":     100,
	}

	if opts != nil {
		if opts.SearchMode != "" {
			body["search_mode"] = opts.SearchMode
		}
		if opts.LogIDs != nil {
			body["log_ids"] = opts.LogIDs
		}
		body["include_matches"] = opts.IncludeMatches
		body["case_sensitive"] = opts.CaseSensitive
		if opts.ContextLines > 0 {
			body["context_lines"] = opts.ContextLines
		}
		if opts.FilePattern != "" {
			body["file_pattern"] = opts.FilePattern
		}
		if opts.MaxMatches > 0 {
			body["max_matches"] = opts.MaxMatches
		}
	}

	var rawResp map[string]interface{}
	err := s.client.post("/service/v2/file-search", body, &rawResp)
	if err != nil {
		return nil, err
	}

	// Handle wrapped or unwrapped response
	resp := &FileSearchJobResponse{Success: true}
	if _, ok := rawResp["success"]; ok {
		jsonData, _ := json.Marshal(rawResp)
		json.Unmarshal(jsonData, resp)
	} else {
		jsonData, _ := json.Marshal(rawResp)
		resp.Data = &FileSearchJobData{}
		json.Unmarshal(jsonData, resp.Data)
	}

	return resp, nil
}

// GetStatus gets file search job status.
func (s *FileSearchService) GetStatus(jobID string) (*FileSearchJobResponse, error) {
	var rawResp map[string]interface{}
	err := s.client.get(fmt.Sprintf("/service/v2/file-search/%s", jobID), nil, &rawResp)
	if err != nil {
		return nil, err
	}

	resp := &FileSearchJobResponse{Success: true}
	if _, ok := rawResp["success"]; ok {
		jsonData, _ := json.Marshal(rawResp)
		json.Unmarshal(jsonData, resp)
	} else {
		jsonData, _ := json.Marshal(rawResp)
		resp.Data = &FileSearchJobData{}
		json.Unmarshal(jsonData, resp.Data)
	}

	return resp, nil
}

// WaitForCompletion waits for a file search job to complete.
func (s *FileSearchService) WaitForCompletion(jobID string, pollInterval, timeout time.Duration) (*FileSearchJobResponse, error) {
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
			return nil, fmt.Errorf("file search job %s did not complete within %v", jobID, timeout)
		}

		// Use server-suggested poll interval if available
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

// Search creates a file search and waits for results.
func (s *FileSearchService) Search(expression string, opts *FileSearchCreateOptions, timeout time.Duration) (*FileSearchJobResponse, error) {
	job, err := s.Create(expression, opts)
	if err != nil {
		return nil, err
	}

	if job.Data == nil || job.Data.JobID == "" {
		return nil, fmt.Errorf("failed to create file search job")
	}

	return s.WaitForCompletion(job.Data.JobID, 2*time.Second, timeout)
}
