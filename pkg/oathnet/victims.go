package oathnet

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
)

const (
	victimsV2SearchPath           = "/service/v2/victims/search"
	filesMetadataV2SearchPath     = "/service/v2/files/search"
	victimPropertiesV2SearchPath  = "/service/v2/victims/properties/search"
	victimPropertiesPathSuffix    = "/properties"
	victimSummaryPathSuffix       = "/summary"
	victimCookiesPathSuffix       = "/cookies"
	victimCookiesDomainPathSuffix = "/cookies/domain"
)

// VictimsService handles V2 victims operations.
type VictimsService struct {
	client *Client
}

// VictimsSearchOptions contains options for victims search.
type VictimsSearchOptions struct {
	Cursor          string
	PageSize        int
	Sort            string
	From            string
	To              string
	DateField       string
	Wildcard        bool
	LogID           string
	Filter          interface{}
	FilterID        string
	TotalDocsMin    int
	TotalDocsMax    int
	ServiceCountMin int
	ServiceCountMax int
	Email           string
	Emails          []string
	EmailDomains    []string
	IP              string
	IPs             []string
	HWIDs           []string
	DiscordID       string
	DiscordIDs      []string
	ComputerName    string
	Username        string
	Usernames       []string
	Country         string
	Countries       []string
	City            string
	Cities          []string
	OS              string
	OSes            []string
	Services        []string
	SteamIDs        []string
	SteamNames      []string
	Phones          []string
	Domains         []string
	Subdomains      []string
	IdentityStates  []string
	VictimIPs       []string
	Antivirus       []string
	InfectionPaths  []string
	Fields          []string
	View            string
	SearchID        string
	ExtraQuery      map[string][]string

	// Deprecated compatibility filter retained for callers of older SDK versions.
	MalwareFamily string
}

// V2FileMetadataSearchRequest is the JSON body for POST /service/v2/files/search.
type V2FileMetadataSearchRequest struct {
	Query    string `json:"q,omitempty"`
	LogID    string `json:"log_id,omitempty"`
	Name     string `json:"name,omitempty"`
	Folder   string `json:"folder,omitempty"`
	Kind     string `json:"kind,omitempty"`
	Ext      string `json:"ext,omitempty"`
	SizeMin  int    `json:"size_min,omitempty"`
	SizeMax  int    `json:"size_max,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
	Cursor   string `json:"cursor,omitempty"`
	SearchID string `json:"search_id,omitempty"`
}

// FileMetadataSearchOptions contains query options for GET /service/v2/files/search.
type FileMetadataSearchOptions = V2FileMetadataSearchRequest

// V2VictimPropertiesSearchRequest is the JSON body for POST /service/v2/victims/properties/search.
type V2VictimPropertiesSearchRequest struct {
	Query                 string      `json:"q,omitempty"`
	LogID                 string      `json:"log_id,omitempty"`
	PropertyType          interface{} `json:"property_type,omitempty"`
	Service               string      `json:"service,omitempty"`
	IdentityKind          string      `json:"identity_kind,omitempty"`
	AccountID             string      `json:"account_id,omitempty"`
	Username              string      `json:"username,omitempty"`
	DisplayName           string      `json:"display_name,omitempty"`
	Value                 string      `json:"value,omitempty"`
	Domain                string      `json:"domain,omitempty"`
	Active                *bool       `json:"active,omitempty"`
	SourceType            string      `json:"source_type,omitempty"`
	SourcePath            string      `json:"source_path,omitempty"`
	SourceFileID          string      `json:"source_file_id,omitempty"`
	Confidence            interface{} `json:"confidence,omitempty"`
	ConfidenceMin         *float64    `json:"confidence_min,omitempty"`
	IncludeCookieEvidence *bool       `json:"include_cookie_evidence,omitempty"`
	ExcludeCookieEvidence *bool       `json:"exclude_cookie_evidence,omitempty"`
	PageSize              int         `json:"page_size,omitempty"`
	Cursor                string      `json:"cursor,omitempty"`
	Sort                  string      `json:"sort,omitempty"`
	SearchID              string      `json:"search_id,omitempty"`
}

// VictimPropertiesSearchOptions contains query options for victim property search endpoints.
type VictimPropertiesSearchOptions = V2VictimPropertiesSearchRequest

// VictimSummaryOptions contains query options for victim summary lookup.
type VictimSummaryOptions struct {
	SearchID string
}

// VictimCookieInventoryOptions contains query options for victim cookie inventory lookup.
type VictimCookieInventoryOptions struct {
	Domain       string
	Status       string
	Query        string
	IncludeItems *bool
	PageSize     int
	Cursor       string
	SearchID     string
}

// VictimCookieDomainOptions contains query options for selected-domain cookie inspection.
type VictimCookieDomainOptions struct {
	Domain   string
	FileID   string
	SearchID string
}

// Search searches victim profiles.
func (s *VictimsService) Search(query string, opts *VictimsSearchOptions) (*V2VictimsResponse, error) {
	params, err := buildVictimsV2SearchParams(query, opts, true)
	if err != nil {
		return nil, err
	}

	var rawResp map[string]interface{}
	err = s.client.get(victimsV2SearchPath, params, &rawResp)
	if err != nil {
		return nil, err
	}

	return decodeV2VictimsResponse(rawResp)
}

// SearchPost searches victim profiles using a JSON body for filter/filter_id.
func (s *VictimsService) SearchPost(body V2SearchPostBody, opts *VictimsSearchOptions) (*V2VictimsResponse, error) {
	params, err := buildVictimsV2SearchParams("", opts, false)
	if err != nil {
		return nil, err
	}
	path := pathWithQuery(victimsV2SearchPath, params)

	var requestBody interface{}
	if body != nil {
		requestBody = body
	}

	var rawResp map[string]interface{}
	err = s.client.post(path, requestBody, &rawResp)
	if err != nil {
		return nil, err
	}

	return decodeV2VictimsResponse(rawResp)
}

// SearchVictimsV2Post is an operation-id alias for SearchPost.
func (s *VictimsService) SearchVictimsV2Post(body V2SearchPostBody, opts *VictimsSearchOptions) (*V2VictimsResponse, error) {
	return s.SearchPost(body, opts)
}

// SearchFilesMetadata searches sanitized victim file metadata.
func (s *VictimsService) SearchFilesMetadata(opts *FileMetadataSearchOptions) (*V2FileMetadataSearchResponse, error) {
	params := buildFileMetadataSearchParams(opts)

	var resp V2FileMetadataSearchResponse
	err := s.client.get(filesMetadataV2SearchPath, params, &resp)
	return &resp, err
}

// SearchFilesMetadataPost searches sanitized victim file metadata with a JSON body.
func (s *VictimsService) SearchFilesMetadataPost(body *V2FileMetadataSearchRequest) (*V2FileMetadataSearchResponse, error) {
	var requestBody interface{}
	if body != nil {
		requestBody = body
	}

	var resp V2FileMetadataSearchResponse
	err := s.client.post(filesMetadataV2SearchPath, requestBody, &resp)
	return &resp, err
}

// SearchFilesMetadataV2 is an operation-id alias for SearchFilesMetadata.
func (s *VictimsService) SearchFilesMetadataV2(opts *FileMetadataSearchOptions) (*V2FileMetadataSearchResponse, error) {
	return s.SearchFilesMetadata(opts)
}

// SearchFilesMetadataV2Post is an operation-id alias for SearchFilesMetadataPost.
func (s *VictimsService) SearchFilesMetadataV2Post(body *V2FileMetadataSearchRequest) (*V2FileMetadataSearchResponse, error) {
	return s.SearchFilesMetadataPost(body)
}

// SearchVictimProperties searches sanitized victim properties globally.
func (s *VictimsService) SearchVictimProperties(opts *VictimPropertiesSearchOptions) (*V2VictimPropertiesSearchResponse, error) {
	params := buildVictimPropertiesSearchParams(opts, false)

	var resp V2VictimPropertiesSearchResponse
	err := s.client.get(victimPropertiesV2SearchPath, params, &resp)
	return &resp, err
}

// SearchVictimPropertiesPost searches sanitized victim properties with a JSON body.
func (s *VictimsService) SearchVictimPropertiesPost(body *V2VictimPropertiesSearchRequest) (*V2VictimPropertiesSearchResponse, error) {
	var requestBody interface{}
	if body != nil {
		requestBody = body
	}

	var resp V2VictimPropertiesSearchResponse
	err := s.client.post(victimPropertiesV2SearchPath, requestBody, &resp)
	return &resp, err
}

// SearchVictimPropertiesV2 is an operation-id alias for SearchVictimProperties.
func (s *VictimsService) SearchVictimPropertiesV2(opts *VictimPropertiesSearchOptions) (*V2VictimPropertiesSearchResponse, error) {
	return s.SearchVictimProperties(opts)
}

// SearchVictimPropertiesV2Post is an operation-id alias for SearchVictimPropertiesPost.
func (s *VictimsService) SearchVictimPropertiesV2Post(body *V2VictimPropertiesSearchRequest) (*V2VictimPropertiesSearchResponse, error) {
	return s.SearchVictimPropertiesPost(body)
}

// GetProperties returns sanitized properties for one victim/log.
func (s *VictimsService) GetProperties(logID string, opts *VictimPropertiesSearchOptions) (*V2VictimPropertiesSearchResponse, error) {
	params := buildVictimPropertiesSearchParams(opts, false)

	var resp V2VictimPropertiesSearchResponse
	err := s.client.get(victimV2Path(logID, victimPropertiesPathSuffix), params, &resp)
	return &resp, err
}

// GetVictimPropertiesV2 is an operation-id alias for GetProperties.
func (s *VictimsService) GetVictimPropertiesV2(logID string, opts *VictimPropertiesSearchOptions) (*V2VictimPropertiesSearchResponse, error) {
	return s.GetProperties(logID, opts)
}

// GetSummary returns a deterministic summary for one victim/log.
func (s *VictimsService) GetSummary(logID string, opts *VictimSummaryOptions) (*V2VictimSummaryResponse, error) {
	params := url.Values{}
	if opts != nil {
		addStringParam(params, "search_id", opts.SearchID)
	}

	var resp V2VictimSummaryResponse
	err := s.client.get(victimV2Path(logID, victimSummaryPathSuffix), params, &resp)
	return &resp, err
}

// GetVictimSummaryV2 is an operation-id alias for GetSummary.
func (s *VictimsService) GetVictimSummaryV2(logID string, opts *VictimSummaryOptions) (*V2VictimSummaryResponse, error) {
	return s.GetSummary(logID, opts)
}

// GetCookies returns value-redacted cookie metadata for one victim/log.
func (s *VictimsService) GetCookies(logID string, opts *VictimCookieInventoryOptions) (*V2VictimCookieInventoryResponse, error) {
	params := buildVictimCookieInventoryParams(opts)

	var resp V2VictimCookieInventoryResponse
	err := s.client.get(victimV2Path(logID, victimCookiesPathSuffix), params, &resp)
	return &resp, err
}

// GetVictimCookiesV2 is an operation-id alias for GetCookies.
func (s *VictimsService) GetVictimCookiesV2(logID string, opts *VictimCookieInventoryOptions) (*V2VictimCookieInventoryResponse, error) {
	return s.GetCookies(logID, opts)
}

// InspectCookieDomain returns raw cookie lines for one selected domain.
func (s *VictimsService) InspectCookieDomain(logID string, opts *VictimCookieDomainOptions) (string, error) {
	params := buildVictimCookieDomainParams(opts)
	path := pathWithQuery(victimV2Path(logID, victimCookiesDomainPathSuffix), params)

	data, err := s.client.getRawWithHeaders(path, map[string]string{"Accept": "text/plain"})
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// InspectVictimCookieDomainV2 is an operation-id alias for InspectCookieDomain.
func (s *VictimsService) InspectVictimCookieDomainV2(logID string, opts *VictimCookieDomainOptions) (string, error) {
	return s.InspectCookieDomain(logID, opts)
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

func buildVictimsV2SearchParams(query string, opts *VictimsSearchOptions, includeStructuredQuery bool) (url.Values, error) {
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
	addStringParam(params, "log_id", opts.LogID)
	if includeStructuredQuery {
		if err := addFilterParam(params, opts.Filter); err != nil {
			return nil, err
		}
		addStringParam(params, "filter_id", opts.FilterID)
	}
	addIntParam(params, "total_docs_min", opts.TotalDocsMin)
	addIntParam(params, "total_docs_max", opts.TotalDocsMax)
	addIntParam(params, "service_count_min", opts.ServiceCountMin)
	addIntParam(params, "service_count_max", opts.ServiceCountMax)
	addSingleAndRepeatedParam(params, "email[]", opts.Email, opts.Emails)
	addRepeatedParam(params, "email_domain[]", opts.EmailDomains)
	addSingleAndRepeatedParam(params, "ip[]", opts.IP, opts.IPs)
	addRepeatedParam(params, "hwid[]", opts.HWIDs)
	addSingleAndRepeatedParam(params, "discord_id[]", opts.DiscordID, opts.DiscordIDs)
	addSingleAndRepeatedParam(params, "username[]", opts.ComputerName, appendSingle(opts.Username, opts.Usernames))
	addSingleAndRepeatedParam(params, "country[]", opts.Country, opts.Countries)
	addSingleAndRepeatedParam(params, "city[]", opts.City, opts.Cities)
	addSingleAndRepeatedParam(params, "os[]", opts.OS, opts.OSes)
	addRepeatedParam(params, "service[]", opts.Services)
	addRepeatedParam(params, "steam_id[]", opts.SteamIDs)
	addRepeatedParam(params, "steam_name[]", opts.SteamNames)
	addRepeatedParam(params, "phone[]", opts.Phones)
	addRepeatedParam(params, "domain[]", opts.Domains)
	addRepeatedParam(params, "subdomain[]", opts.Subdomains)
	addRepeatedParam(params, "identity_state[]", opts.IdentityStates)
	addRepeatedParam(params, "victim_ip[]", opts.VictimIPs)
	addRepeatedParam(params, "antivirus[]", opts.Antivirus)
	addRepeatedParam(params, "infection_path[]", opts.InfectionPaths)
	addRepeatedParam(params, "fields[]", opts.Fields)
	addStringParam(params, "view", opts.View)
	addStringParam(params, "search_id", opts.SearchID)

	addStringParam(params, "malware_family", opts.MalwareFamily)

	for key, values := range opts.ExtraQuery {
		addRepeatedParam(params, key, values)
	}

	return params, nil
}

func decodeV2VictimsResponse(rawResp map[string]interface{}) (*V2VictimsResponse, error) {
	resp := &V2VictimsResponse{Success: true}
	jsonData, err := json.Marshal(rawResp)
	if err != nil {
		return nil, err
	}
	if _, ok := rawResp["success"]; ok {
		if err := json.Unmarshal(jsonData, resp); err != nil {
			return nil, err
		}
		return resp, nil
	}

	resp.Data = &V2VictimsData{}
	if err := json.Unmarshal(jsonData, resp.Data); err != nil {
		return nil, err
	}
	return resp, nil
}

func buildFileMetadataSearchParams(opts *FileMetadataSearchOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		return params
	}

	addStringParam(params, "q", opts.Query)
	addStringParam(params, "log_id", opts.LogID)
	addStringParam(params, "name", opts.Name)
	addStringParam(params, "folder", opts.Folder)
	addStringParam(params, "kind", opts.Kind)
	addStringParam(params, "ext", opts.Ext)
	addIntParam(params, "size_min", opts.SizeMin)
	addIntParam(params, "size_max", opts.SizeMax)
	addIntParam(params, "page_size", opts.PageSize)
	addStringParam(params, "cursor", opts.Cursor)
	addStringParam(params, "search_id", opts.SearchID)
	return params
}

func buildVictimPropertiesSearchParams(opts *VictimPropertiesSearchOptions, includeLogID bool) url.Values {
	params := url.Values{}
	if opts == nil {
		return params
	}

	addStringParam(params, "q", opts.Query)
	if includeLogID {
		addStringParam(params, "log_id", opts.LogID)
	}
	addStringOrStringSliceParam(params, "property_type", opts.PropertyType)
	addStringParam(params, "service", opts.Service)
	addStringParam(params, "identity_kind", opts.IdentityKind)
	addStringParam(params, "account_id", opts.AccountID)
	addStringParam(params, "username", opts.Username)
	addStringParam(params, "display_name", opts.DisplayName)
	addStringParam(params, "value", opts.Value)
	addStringParam(params, "domain", opts.Domain)
	addBoolPtrParam(params, "active", opts.Active)
	addStringParam(params, "source_type", opts.SourceType)
	addStringOrStringSliceParam(params, "confidence", opts.Confidence)
	addFloatPtrParam(params, "confidence_min", opts.ConfidenceMin)
	addBoolPtrParam(params, "include_cookie_evidence", opts.IncludeCookieEvidence)
	addBoolPtrParam(params, "exclude_cookie_evidence", opts.ExcludeCookieEvidence)
	addIntParam(params, "page_size", opts.PageSize)
	addStringParam(params, "cursor", opts.Cursor)
	addStringParam(params, "sort", opts.Sort)
	addStringParam(params, "search_id", opts.SearchID)
	return params
}

func buildVictimCookieInventoryParams(opts *VictimCookieInventoryOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		return params
	}

	addStringParam(params, "domain", opts.Domain)
	addStringParam(params, "status", opts.Status)
	addStringParam(params, "q", opts.Query)
	addBoolPtrParam(params, "include_items", opts.IncludeItems)
	addIntParam(params, "page_size", opts.PageSize)
	addStringParam(params, "cursor", opts.Cursor)
	addStringParam(params, "search_id", opts.SearchID)
	return params
}

func buildVictimCookieDomainParams(opts *VictimCookieDomainOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		return params
	}

	addStringParam(params, "domain", opts.Domain)
	addStringParam(params, "file_id", opts.FileID)
	addStringParam(params, "search_id", opts.SearchID)
	return params
}

func victimV2Path(logID, suffix string) string {
	return fmt.Sprintf("/service/v2/victims/%s%s", url.PathEscape(logID), suffix)
}

func appendSingle(single string, values []string) []string {
	if single == "" {
		return values
	}
	out := make([]string, 0, len(values)+1)
	out = append(out, single)
	out = append(out, values...)
	return out
}

func addStringOrStringSliceParam(params url.Values, key string, value interface{}) {
	switch typed := value.(type) {
	case string:
		addStringParam(params, key, typed)
	case []string:
		addRepeatedParam(params, key, typed)
	}
}

func addBoolPtrParam(params url.Values, key string, value *bool) {
	if value != nil {
		params.Set(key, strconv.FormatBool(*value))
	}
}

func addFloatPtrParam(params url.Values, key string, value *float64) {
	if value != nil {
		params.Set(key, strconv.FormatFloat(*value, 'f', -1, 64))
	}
}
