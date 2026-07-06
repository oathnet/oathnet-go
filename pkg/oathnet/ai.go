package oathnet

import (
	"fmt"
	"net/url"
)

const aiFilterBasePath = "/service/v2/ai/filter"

// AIService handles V2 AI filter operations.
type AIService struct {
	client *Client
}

// AIFilterIndex is a supported AI filter target index.
type AIFilterIndex string

const (
	AIFilterIndexBreach               AIFilterIndex = "breach"
	AIFilterIndexDocs                 AIFilterIndex = "docs"
	AIFilterIndexVictims              AIFilterIndex = "victims"
	AIFilterIndexInvestigation        AIFilterIndex = "investigation"
	AIFilterIndexStealerInvestigation AIFilterIndex = "stealer_investigation"
)

// StructuredFilterNode is the recursive V2 structured filter tree returned by
// AI filter endpoints and accepted by V2 search/filter-aware APIs.
type StructuredFilterNode = StructuredFilter

// V2AIFilterRequest is the request body for POST /service/v2/ai/filter.
type V2AIFilterRequest struct {
	Query    string        `json:"query"`
	Index    AIFilterIndex `json:"index,omitempty"`
	FilterID string        `json:"filter_id,omitempty"`
}

// V2AIResponse is the raw AI filter creation response.
type V2AIResponse struct {
	FilterID string               `json:"filter_id,omitempty"`
	Filter   StructuredFilterNode `json:"filter,omitempty"`
}

// V2AIHistoryItem is one prior natural-language filter turn in a filter context.
type V2AIHistoryItem struct {
	Query  string               `json:"query,omitempty"`
	Filter StructuredFilterNode `json:"filter,omitempty"`
}

// V2AIContextResponse is the raw AI filter context response.
type V2AIContextResponse struct {
	ID          string                 `json:"id,omitempty"`
	IndexType   string                 `json:"index_type,omitempty"`
	Query       string                 `json:"query,omitempty"`
	Filter      StructuredFilterNode   `json:"filter,omitempty"`
	SampleData  map[string]interface{} `json:"sample_data,omitempty"`
	FieldValues map[string][]string    `json:"field_values,omitempty"`
	TotalHits   int64                  `json:"total_hits,omitempty"`
	History     []V2AIHistoryItem      `json:"history,omitempty"`
	Source      string                 `json:"source,omitempty"`
	ParentID    string                 `json:"parent_id,omitempty"`
	CreatedAt   string                 `json:"created_at,omitempty"`
	ExpiresAt   string                 `json:"expires_at,omitempty"`
}

// Create translates natural language into a structured filter context.
func (s *AIService) Create(req V2AIFilterRequest) (*V2AIResponse, error) {
	var resp V2AIResponse
	err := s.client.post(aiFilterBasePath, req, &resp)
	return &resp, err
}

// GetContext gets a transient AI/search filter context by ID.
func (s *AIService) GetContext(filterID string) (*V2AIContextResponse, error) {
	var resp V2AIContextResponse
	err := s.client.get(aiFilterContextPath(filterID), nil, &resp)
	return &resp, err
}

func aiFilterContextPath(filterID string) string {
	return fmt.Sprintf("%s/%s", aiFilterBasePath, url.PathEscape(filterID))
}
