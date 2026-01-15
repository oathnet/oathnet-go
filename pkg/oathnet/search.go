package oathnet

import "net/url"

// SearchService handles search operations.
type SearchService struct {
	client *Client
}

// SearchOptions contains options for breach/stealer search.
type SearchOptions struct {
	Cursor   string
	DBNames  string
	SearchID string
}

// InitSession initializes a search session.
func (s *SearchService) InitSession(query string) (*SearchSessionResponse, error) {
	body := map[string]interface{}{
		"query": query,
	}

	var resp SearchSessionResponse
	err := s.client.post("/service/search/init", body, &resp)
	return &resp, err
}

// Breach searches the breach database.
func (s *SearchService) Breach(query string, opts *SearchOptions) (*BreachSearchResponse, error) {
	params := url.Values{}
	params.Set("q", query)

	if opts != nil {
		if opts.Cursor != "" {
			params.Set("cursor", opts.Cursor)
		}
		if opts.DBNames != "" {
			params.Set("dbnames", opts.DBNames)
		}
		if opts.SearchID != "" {
			params.Set("search_id", opts.SearchID)
		}
	}

	var resp BreachSearchResponse
	err := s.client.get("/service/search-breach", params, &resp)
	return &resp, err
}

// Stealer searches the stealer database (legacy).
func (s *SearchService) Stealer(query string, opts *SearchOptions) (*StealerSearchResponse, error) {
	params := url.Values{}
	params.Set("q", query)

	if opts != nil {
		if opts.Cursor != "" {
			params.Set("cursor", opts.Cursor)
		}
		if opts.DBNames != "" {
			params.Set("dbnames", opts.DBNames)
		}
		if opts.SearchID != "" {
			params.Set("search_id", opts.SearchID)
		}
	}

	var resp StealerSearchResponse
	err := s.client.get("/service/search-stealer", params, &resp)
	return &resp, err
}
