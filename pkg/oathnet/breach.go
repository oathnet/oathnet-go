package oathnet

import (
	"encoding/json"
	"net/url"
	"strconv"
)

const breachV2SearchPath = "/service/v2/breach/search"

// BreachV2Service handles V2 breach search and autocomplete operations.
type BreachV2Service struct {
	client *Client
}

// StructuredFilter is a JSON-compatible V2 structured filter tree.
type StructuredFilter map[string]interface{}

// V2SearchPostBody is the JSON body used by POST variants of V2 search endpoints.
type V2SearchPostBody map[string]interface{}

// BreachV2SearchOptions contains query options for V2 breach search.
type BreachV2SearchOptions struct {
	Cursor     string
	PageSize   int
	Sort       string
	From       string
	To         string
	DateField  string
	Wildcard   bool
	Logic      string
	Filter     interface{}
	FilterID   string
	Fields     []string
	SearchID   string
	ExtraQuery map[string][]string

	Email          []string
	EmailDomains   []string
	Domains        []string
	Usernames      []string
	Passwords      []string
	PasswordHashes []string
	IPs            []string
	Phones         []string
	FirstNames     []string
	LastNames      []string
	FullNames      []string
	Cities         []string
	Countries      []string
	States         []string
	PostalCodes    []string
	DBNames        []string
	DiscordIDs     []string
	IBANs          []string
	SSNs           []string
	DateBirthFrom  string
	DateBirthTo    string
	Names          []string
	Genders        []string
	Addresses      []string
	Discord        []string
	Social         []string
	Financial      []string
	Gaming         []string
}

// BreachV2AutocompleteOptions contains options for value autocomplete.
type BreachV2AutocompleteOptions struct {
	Field       string
	Query       string
	Limit       int
	IncludeInfo bool
}

// BreachV2AutocompleteDBNamesOptions contains options for DB-name autocomplete.
type BreachV2AutocompleteDBNamesOptions struct {
	Query string
	Limit int
}

// BreachV2AutocompleteFieldsOptions contains options for field coverage lookup.
type BreachV2AutocompleteFieldsOptions struct {
	Field string
	Limit int
}

// Search searches V2 breach records with query parameters.
func (s *BreachV2Service) Search(query string, opts *BreachV2SearchOptions) (*V2BreachSearchResponse, error) {
	params, err := buildBreachV2SearchParams(query, opts, true)
	if err != nil {
		return nil, err
	}

	var resp V2BreachSearchResponse
	err = s.client.get(breachV2SearchPath, params, &resp)
	return &resp, err
}

// SearchPost searches V2 breach records using a JSON body for filter/filter_id.
func (s *BreachV2Service) SearchPost(body V2SearchPostBody, opts *BreachV2SearchOptions) (*V2BreachSearchResponse, error) {
	params, err := buildBreachV2SearchParams("", opts, false)
	if err != nil {
		return nil, err
	}
	path := pathWithQuery(breachV2SearchPath, params)

	var requestBody interface{}
	if body != nil {
		requestBody = body
	}

	var resp V2BreachSearchResponse
	err = s.client.post(path, requestBody, &resp)
	return &resp, err
}

// Autocomplete returns raw autocomplete values for supported breach fields.
func (s *BreachV2Service) Autocomplete(opts *BreachV2AutocompleteOptions) (*V2AutocompleteValueResponse, error) {
	params := url.Values{}
	if opts != nil {
		addStringParam(params, "field", opts.Field)
		addStringParam(params, "q", opts.Query)
		addIntParam(params, "limit", opts.Limit)
		addBoolParam(params, "include_info", opts.IncludeInfo)
	}

	var resp V2AutocompleteValueResponse
	err := s.client.get("/service/v2/breach/autocomplete", params, &resp)
	return &resp, err
}

// AutocompleteDBNames returns raw DB-name autocomplete results.
func (s *BreachV2Service) AutocompleteDBNames(opts *BreachV2AutocompleteDBNamesOptions) (*V2AutocompleteDBNamesResponse, error) {
	params := url.Values{}
	if opts != nil {
		addStringParam(params, "q", opts.Query)
		addIntParam(params, "limit", opts.Limit)
	}

	var resp V2AutocompleteDBNamesResponse
	err := s.client.get("/service/v2/breach/autocomplete/dbnames", params, &resp)
	return &resp, err
}

// AutocompleteFields returns DB names where a breach field is populated.
func (s *BreachV2Service) AutocompleteFields(opts *BreachV2AutocompleteFieldsOptions) (*V2AutocompleteFieldsResponse, error) {
	params := url.Values{}
	if opts != nil {
		addStringParam(params, "field", opts.Field)
		addIntParam(params, "limit", opts.Limit)
	}

	var resp V2AutocompleteFieldsResponse
	err := s.client.get("/service/v2/breach/autocomplete/fields", params, &resp)
	return &resp, err
}

func buildBreachV2SearchParams(query string, opts *BreachV2SearchOptions, includeStructuredQuery bool) (url.Values, error) {
	params := url.Values{}
	addStringParam(params, "q", query)
	if opts == nil {
		return params, nil
	}

	addStringParam(params, "cursor", opts.Cursor)
	addIntParam(params, "page_size", opts.PageSize)
	addStringParam(params, "sort", opts.Sort)
	addStringParam(params, "from", opts.From)
	addStringParam(params, "to", opts.To)
	addStringParam(params, "date_field", opts.DateField)
	addBoolParam(params, "wildcard", opts.Wildcard)
	addStringParam(params, "logic", opts.Logic)
	if includeStructuredQuery {
		if err := addFilterParam(params, opts.Filter); err != nil {
			return nil, err
		}
		addStringParam(params, "filter_id", opts.FilterID)
	}

	addRepeatedParam(params, "email[]", opts.Email)
	addRepeatedParam(params, "email_domain[]", opts.EmailDomains)
	addRepeatedParam(params, "domain[]", opts.Domains)
	addRepeatedParam(params, "username[]", opts.Usernames)
	addRepeatedParam(params, "password[]", opts.Passwords)
	addRepeatedParam(params, "password_hash[]", opts.PasswordHashes)
	addRepeatedParam(params, "ip[]", opts.IPs)
	addRepeatedParam(params, "phone[]", opts.Phones)
	addRepeatedParam(params, "first_name[]", opts.FirstNames)
	addRepeatedParam(params, "last_name[]", opts.LastNames)
	addRepeatedParam(params, "full_name[]", opts.FullNames)
	addRepeatedParam(params, "city[]", opts.Cities)
	addRepeatedParam(params, "country[]", opts.Countries)
	addRepeatedParam(params, "state[]", opts.States)
	addRepeatedParam(params, "postal_code[]", opts.PostalCodes)
	addRepeatedParam(params, "dbname[]", opts.DBNames)
	addRepeatedParam(params, "discord_id[]", opts.DiscordIDs)
	addRepeatedParam(params, "iban[]", opts.IBANs)
	addRepeatedParam(params, "ssn[]", opts.SSNs)
	addStringParam(params, "date_birth_from", opts.DateBirthFrom)
	addStringParam(params, "date_birth_to", opts.DateBirthTo)
	addRepeatedParam(params, "name[]", opts.Names)
	addRepeatedParam(params, "gender[]", opts.Genders)
	addRepeatedParam(params, "address[]", opts.Addresses)
	addRepeatedParam(params, "discord[]", opts.Discord)
	addRepeatedParam(params, "social[]", opts.Social)
	addRepeatedParam(params, "financial[]", opts.Financial)
	addRepeatedParam(params, "gaming[]", opts.Gaming)
	addRepeatedParam(params, "fields[]", opts.Fields)
	addStringParam(params, "search_id", opts.SearchID)

	for key, values := range opts.ExtraQuery {
		addRepeatedParam(params, key, values)
	}

	return params, nil
}

func addStringParam(params url.Values, key, value string) {
	if value != "" {
		params.Set(key, value)
	}
}

func addIntParam(params url.Values, key string, value int) {
	if value > 0 {
		params.Set(key, strconv.Itoa(value))
	}
}

func addBoolParam(params url.Values, key string, value bool) {
	if value {
		params.Set(key, "true")
	}
}

func addRepeatedParam(params url.Values, key string, values []string) {
	for _, value := range values {
		if value != "" {
			params.Add(key, value)
		}
	}
}

func addFilterParam(params url.Values, filter interface{}) error {
	if filter == nil {
		return nil
	}
	if encoded, ok := filter.(string); ok {
		addStringParam(params, "filter", encoded)
		return nil
	}

	data, err := json.Marshal(filter)
	if err != nil {
		return err
	}
	addStringParam(params, "filter", string(data))
	return nil
}

func pathWithQuery(path string, params url.Values) string {
	if len(params) == 0 {
		return path
	}
	return path + "?" + params.Encode()
}
