package oathnet

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"
)

const exportsV2BasePath = "/service/v2/exports"

// ExportsService handles V2 export operations.
type ExportsService struct {
	client *Client
}

// ExportCreateOptions contains options for creating an export job.
type ExportCreateOptions struct {
	Format string // jsonl, csv
	Limit  int
	Fields []string
	Search map[string]string
}

// Create creates an export job.
func (s *ExportsService) Create(exportType string, opts *ExportCreateOptions) (*ExportJobResponse, error) {
	body := map[string]interface{}{
		"type":   exportType,
		"format": "jsonl",
	}

	if opts != nil {
		if opts.Format != "" {
			body["format"] = opts.Format
		}
		if opts.Limit > 0 {
			body["limit"] = opts.Limit
		}
		if opts.Fields != nil {
			body["fields"] = opts.Fields
		}
		if opts.Search != nil {
			body["search"] = opts.Search
		}
	}

	var rawResp map[string]interface{}
	err := s.client.post(exportsV2BasePath, body, &rawResp)
	if err != nil {
		return nil, err
	}

	resp := &ExportJobResponse{Success: true}
	if _, ok := rawResp["success"]; ok {
		jsonData, _ := json.Marshal(rawResp)
		json.Unmarshal(jsonData, resp)
	} else {
		jsonData, _ := json.Marshal(rawResp)
		resp.Data = &ExportJobData{}
		json.Unmarshal(jsonData, resp.Data)
	}

	return resp, nil
}

// GetStatus gets export job status.
func (s *ExportsService) GetStatus(jobID string) (*ExportJobResponse, error) {
	var rawResp map[string]interface{}
	err := s.client.get(exportJobPath(jobID, ""), nil, &rawResp)
	if err != nil {
		return nil, err
	}

	resp := &ExportJobResponse{Success: true}
	if _, ok := rawResp["success"]; ok {
		jsonData, _ := json.Marshal(rawResp)
		json.Unmarshal(jsonData, resp)
	} else {
		jsonData, _ := json.Marshal(rawResp)
		resp.Data = &ExportJobData{}
		json.Unmarshal(jsonData, resp.Data)
	}

	return resp, nil
}

// List lists export jobs.
func (s *ExportsService) List(page, pageSize int) (*V2ExportJobListResponse, error) {
	params := url.Values{}
	addIntParam(params, "page", page)
	addIntParam(params, "page_size", pageSize)

	var resp V2ExportJobListResponse
	err := s.client.get(exportsV2BasePath+"/list", params, &resp)
	return &resp, err
}

// ListExportsV2 is an operation-id alias for List.
func (s *ExportsService) ListExportsV2(page, pageSize int) (*V2ExportJobListResponse, error) {
	return s.List(page, pageSize)
}

// Download downloads the export file.
func (s *ExportsService) Download(jobID, outputPath string) error {
	data, err := s.client.getRaw(exportJobPath(jobID, "/download"))
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, data, 0644)
}

// WaitForCompletion waits for an export job to complete.
func (s *ExportsService) WaitForCompletion(jobID string, pollInterval, timeout time.Duration) (*ExportJobResponse, error) {
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
			return nil, fmt.Errorf("export job %s did not complete within %v", jobID, timeout)
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

// Export creates an export, waits for completion, and downloads.
func (s *ExportsService) Export(exportType, outputPath string, opts *ExportCreateOptions, timeout time.Duration) error {
	job, err := s.Create(exportType, opts)
	if err != nil {
		return err
	}

	if job.Data == nil || job.Data.JobID == "" {
		return fmt.Errorf("failed to create export job")
	}

	_, err = s.WaitForCompletion(job.Data.JobID, 2*time.Second, timeout)
	if err != nil {
		return err
	}

	return s.Download(job.Data.JobID, outputPath)
}

func exportJobPath(jobID, suffix string) string {
	return fmt.Sprintf("%s/%s%s", exportsV2BasePath, url.PathEscape(jobID), suffix)
}
