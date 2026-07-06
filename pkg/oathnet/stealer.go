package oathnet

import (
	"encoding/json"
	"net/url"
	"strconv"
)

const (
	stealerV2SearchPath              = "/service/v2/stealer/search"
	stealerV2InvestigationSearchPath = "/service/v2/stealer/investigation/search"
	investigationV2AliasPath         = "/service/v2/investigate/search"
	phonebookV2Path                  = "/service/v2/phonebook"
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

// InvestigationSearchOptions contains query options for GET investigation search.
type InvestigationSearchOptions struct {
	Scope                 string
	Include               []string
	PageSize              int
	SearchID              string
	Filter                interface{}
	FilterID              string
	FilterMode            string
	Compact               bool
	View                  string
	IncludeCookieEvidence bool
	ExcludeCookieEvidence bool
	ExtraQuery            map[string][]string
}

// V2InvestigationSectionFilters contains section-specific flat filters for
// POST investigation search.
type V2InvestigationSectionFilters map[string]map[string]interface{}

// V2InvestigationCursors contains per-section cursor tokens. A nil value is
// encoded as JSON null, matching the documented frontend-style payload.
type V2InvestigationCursors map[string]*string

// V2InvestigationSearchRequest is the JSON body for POST investigation search.
type V2InvestigationSearchRequest struct {
	Q          string                        `json:"q,omitempty"`
	Scope      string                        `json:"scope,omitempty"`
	Include    []string                      `json:"include,omitempty"`
	FilterMode string                        `json:"filter_mode,omitempty"`
	Compact    bool                          `json:"compact,omitempty"`
	PageSize   int                           `json:"page_size,omitempty"`
	View       string                        `json:"view,omitempty"`
	SearchID   string                        `json:"search_id,omitempty"`
	Wildcard   bool                          `json:"wildcard,omitempty"`
	From       string                        `json:"from,omitempty"`
	To         string                        `json:"to,omitempty"`
	DateField  string                        `json:"date_field,omitempty"`
	LogID      string                        `json:"log_id,omitempty"`
	HasLogID   bool                          `json:"has_log_id,omitempty"`
	Sort       string                        `json:"sort,omitempty"`
	Fields     []string                      `json:"fields,omitempty"`
	Filter     interface{}                   `json:"filter,omitempty"`
	FilterID   string                        `json:"filter_id,omitempty"`
	Filters    V2InvestigationSectionFilters `json:"filters,omitempty"`
	Cursors    V2InvestigationCursors        `json:"cursors,omitempty"`
	Extra      map[string]interface{}        `json:"-"`
}

// MarshalJSON merges documented request fields with additional OpenAPI
// properties. Extra values intentionally override omitted/zero typed fields so
// callers can send explicit false or experimental fields when needed.
func (r V2InvestigationSearchRequest) MarshalJSON() ([]byte, error) {
	type requestAlias V2InvestigationSearchRequest

	data, err := json.Marshal(requestAlias(r))
	if err != nil {
		return nil, err
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	for key, value := range r.Extra {
		payload[key] = value
	}

	return json.Marshal(payload)
}

// PhonebookOptions contains options for V2 phonebook domain intelligence.
type PhonebookOptions struct {
	Query    string
	Alive    bool
	IsAlive  bool
	SearchID string
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
	err := s.client.get(stealerV2SearchPath, params, &resp)
	return &resp, err
}

// InvestigationSearch searches the canonical V2 stealer investigation route.
func (s *StealerV2Service) InvestigationSearch(query string, opts *InvestigationSearchOptions) (*V2InvestigationSearchResponse, error) {
	return s.investigationSearch(stealerV2InvestigationSearchPath, query, opts)
}

// InvestigationSearchPost searches the canonical V2 stealer investigation route with a JSON body.
func (s *StealerV2Service) InvestigationSearchPost(body *V2InvestigationSearchRequest) (*V2InvestigationSearchResponse, error) {
	return s.investigationSearchPost(stealerV2InvestigationSearchPath, body)
}

// InvestigationSearchAlias searches the legacy V2 investigation alias route.
// Prefer InvestigationSearch for new integrations.
func (s *StealerV2Service) InvestigationSearchAlias(query string, opts *InvestigationSearchOptions) (*V2InvestigationSearchResponse, error) {
	return s.investigationSearch(investigationV2AliasPath, query, opts)
}

// InvestigationSearchAliasPost searches the legacy V2 investigation alias route with a JSON body.
// Prefer InvestigationSearchPost for new integrations.
func (s *StealerV2Service) InvestigationSearchAliasPost(body *V2InvestigationSearchRequest) (*V2InvestigationSearchResponse, error) {
	return s.investigationSearchPost(investigationV2AliasPath, body)
}

func (s *StealerV2Service) investigationSearch(path, query string, opts *InvestigationSearchOptions) (*V2InvestigationSearchResponse, error) {
	params, err := buildInvestigationSearchParams(query, opts)
	if err != nil {
		return nil, err
	}

	var resp V2InvestigationSearchResponse
	err = s.client.get(path, params, &resp)
	return &resp, err
}

func (s *StealerV2Service) investigationSearchPost(path string, body *V2InvestigationSearchRequest) (*V2InvestigationSearchResponse, error) {
	var requestBody interface{}
	if body != nil {
		requestBody = body
	}

	var resp V2InvestigationSearchResponse
	err := s.client.post(path, requestBody, &resp)
	return &resp, err
}

// GetPhonebook returns raw V2 phonebook domain intelligence.
func (s *StealerV2Service) GetPhonebook(domain string, opts *PhonebookOptions) (*V2PhonebookResponse, error) {
	params := buildPhonebookParams(domain, opts)

	var resp V2PhonebookResponse
	err := s.client.get(phonebookV2Path, params, &resp)
	return &resp, err
}

// Phonebook returns raw V2 phonebook domain intelligence.
func (s *StealerV2Service) Phonebook(domain string, opts *PhonebookOptions) (*V2PhonebookResponse, error) {
	return s.GetPhonebook(domain, opts)
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

func buildInvestigationSearchParams(query string, opts *InvestigationSearchOptions) (url.Values, error) {
	params := url.Values{}
	addStringParam(params, "q", query)
	if opts == nil {
		return params, nil
	}

	addStringParam(params, "scope", opts.Scope)
	addRepeatedParam(params, "include", opts.Include)
	addIntParam(params, "page_size", opts.PageSize)
	addStringParam(params, "search_id", opts.SearchID)
	if err := addFilterParam(params, opts.Filter); err != nil {
		return nil, err
	}
	addStringParam(params, "filter_id", opts.FilterID)
	addStringParam(params, "filter_mode", opts.FilterMode)
	addBoolParam(params, "compact", opts.Compact)
	addStringParam(params, "view", opts.View)
	addBoolParam(params, "include_cookie_evidence", opts.IncludeCookieEvidence)
	addBoolParam(params, "exclude_cookie_evidence", opts.ExcludeCookieEvidence)

	for key, values := range opts.ExtraQuery {
		addRepeatedParam(params, key, values)
	}

	return params, nil
}

func buildPhonebookParams(domain string, opts *PhonebookOptions) url.Values {
	params := url.Values{}
	addStringParam(params, "domain", domain)
	if opts == nil {
		return params
	}

	addStringParam(params, "q", opts.Query)
	addBoolParam(params, "alive", opts.Alive)
	addBoolParam(params, "is_alive", opts.IsAlive)
	addStringParam(params, "search_id", opts.SearchID)

	return params
}
