package oathnet

import "net/url"

// UtilityService handles utility operations.
type UtilityService struct {
	client *Client
}

// DBNameAutocomplete autocompletes database names.
func (s *UtilityService) DBNameAutocomplete(query string) ([]string, error) {
	params := url.Values{}
	params.Set("q", query)

	var resp []string
	err := s.client.get("/service/dbname-autocomplete", params, &resp)
	return resp, err
}
