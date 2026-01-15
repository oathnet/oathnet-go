package oathnet

import (
	"net/url"
	"strconv"
)

// StealerV2Service handles V2 stealer operations.
type StealerV2Service struct {
	client *Client
}

// StealerSearchOptions contains options for V2 stealer search.
type StealerSearchOptions struct {
	Cursor        string
	PageSize      int
	Domain        string
	Subdomain     string
	Email         string
	Username      string
	Password      string
	Wildcard      bool
	HasLogID      bool
	Country       string
	MalwareFamily string
}

// Search searches the V2 stealer database.
func (s *StealerV2Service) Search(query string, opts *StealerSearchOptions) (*V2StealerResponse, error) {
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
		if opts.Domain != "" {
			params.Add("domain[]", opts.Domain)
		}
		if opts.Subdomain != "" {
			params.Add("subdomain[]", opts.Subdomain)
		}
		if opts.Email != "" {
			params.Set("email", opts.Email)
		}
		if opts.Username != "" {
			params.Add("username[]", opts.Username)
		}
		if opts.Password != "" {
			params.Add("password[]", opts.Password)
		}
		if opts.Wildcard {
			params.Set("wildcard", "true")
		}
		if opts.HasLogID {
			params.Set("has_log_id", "true")
		}
		if opts.Country != "" {
			params.Set("country", opts.Country)
		}
		if opts.MalwareFamily != "" {
			params.Set("malware_family", opts.MalwareFamily)
		}
	}

	var resp V2StealerResponse
	err := s.client.get("/service/v2/stealer/search", params, &resp)
	return &resp, err
}

// Subdomain extracts subdomains from stealer data.
func (s *StealerV2Service) Subdomain(domain string, query string) (*SubdomainResponse, error) {
	params := url.Values{}
	params.Set("domain", domain)
	if query != "" {
		params.Set("q", query)
	}

	var resp SubdomainResponse
	err := s.client.get("/service/v2/stealer/subdomain", params, &resp)
	return &resp, err
}
