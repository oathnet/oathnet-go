package oathnet

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
)

// VictimsService handles V2 victims operations.
type VictimsService struct {
	client *Client
}

// VictimsSearchOptions contains options for victims search.
type VictimsSearchOptions struct {
	Cursor        string
	PageSize      int
	Email         string
	IP            string
	DiscordID     string
	ComputerName  string
	Country       string
	MalwareFamily string
	Wildcard      bool
}

// Search searches victim profiles.
func (s *VictimsService) Search(query string, opts *VictimsSearchOptions) (*V2VictimsResponse, error) {
	params := url.Values{}
	if query != "" {
		params.Set("q", query)
	}

	if opts != nil {
		if opts.Cursor != "" {
			params.Set("cursor", opts.Cursor)
		}
		if opts.PageSize > 0 {
			params.Set("page_size", strconv.Itoa(opts.PageSize))
		}
		if opts.Email != "" {
			params.Add("email[]", opts.Email)
		}
		if opts.IP != "" {
			params.Add("ip[]", opts.IP)
		}
		if opts.DiscordID != "" {
			params.Add("discord_id[]", opts.DiscordID)
		}
		if opts.ComputerName != "" {
			params.Add("username[]", opts.ComputerName)
		}
		if opts.Country != "" {
			params.Set("country", opts.Country)
		}
		if opts.MalwareFamily != "" {
			params.Set("malware_family", opts.MalwareFamily)
		}
		if opts.Wildcard {
			params.Set("wildcard", "true")
		}
	}

	var rawResp map[string]interface{}
	err := s.client.get("/service/v2/victims/search", params, &rawResp)
	if err != nil {
		return nil, err
	}

	// Handle wrapped or unwrapped response
	resp := &V2VictimsResponse{Success: true}
	if _, ok := rawResp["success"]; ok {
		jsonData, _ := json.Marshal(rawResp)
		json.Unmarshal(jsonData, resp)
	} else {
		// API returns unwrapped response
		jsonData, _ := json.Marshal(rawResp)
		resp.Data = &V2VictimsData{}
		json.Unmarshal(jsonData, resp.Data)
	}

	return resp, nil
}

// GetManifest gets the victim file manifest (file tree).
// Note: This endpoint returns unwrapped response.
func (s *VictimsService) GetManifest(logID string) (*VictimManifestData, error) {
	var resp VictimManifestData
	err := s.client.get(fmt.Sprintf("/service/v2/victims/%s", logID), nil, &resp)
	return &resp, err
}

// GetFile gets victim file content.
func (s *VictimsService) GetFile(logID, fileID string) ([]byte, error) {
	return s.client.getRaw(fmt.Sprintf("/service/v2/victims/%s/files/%s", logID, fileID))
}

// DownloadArchive downloads victim archive as ZIP.
func (s *VictimsService) DownloadArchive(logID string, outputPath string) error {
	data, err := s.client.getRaw(fmt.Sprintf("/service/v2/victims/%s/archive", logID))
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, data, 0644)
}
