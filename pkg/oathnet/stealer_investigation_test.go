package oathnet

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestStealerV2Service_InvestigationSearchRequestConstruction(t *testing.T) {
	var gotMethod string
	var gotPath string
	var gotRawQuery string
	var gotQuery url.Values
	var gotAPIKey string
	var gotBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.EscapedPath()
		gotRawQuery = r.URL.RawQuery
		gotQuery = r.URL.Query()
		gotAPIKey = r.Header.Get("x-api-key")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("ReadAll() error = %v", err)
		}
		gotBody = body

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(investigationSearchResponseJSON))
	}))
	defer server.Close()

	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.Stealer.InvestigationSearch("alice+tag@example.com", &InvestigationSearchOptions{
		Scope:                 "all",
		Include:               []string{"credentials", "victims", "evidence"},
		PageSize:              25,
		SearchID:              "sess/123",
		Filter:                StructuredFilter{"field": "domain", "operator": "eq", "value": "example.com"},
		FilterID:              "0123456789abcdef01234567",
		FilterMode:            "intersect",
		Compact:               true,
		View:                  "enriched",
		IncludeCookieEvidence: true,
		ExcludeCookieEvidence: true,
		ExtraQuery: map[string][]string{
			"custom[]": {"value/with space"},
		},
	})
	if err != nil {
		t.Fatalf("InvestigationSearch() error = %v", err)
	}

	if gotMethod != http.MethodGet {
		t.Fatalf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/service/v2/stealer/investigation/search" {
		t.Fatalf("path = %s, want canonical investigation path", gotPath)
	}
	if gotAPIKey != "test-key" {
		t.Fatalf("x-api-key = %q, want test-key", gotAPIKey)
	}
	if len(gotBody) != 0 {
		t.Fatalf("GET body = %s, want empty", string(gotBody))
	}

	wantQuery := url.Values{
		"q":                       {"alice+tag@example.com"},
		"scope":                   {"all"},
		"include":                 {"credentials", "victims", "evidence"},
		"page_size":               {"25"},
		"search_id":               {"sess/123"},
		"filter":                  {`{"field":"domain","operator":"eq","value":"example.com"}`},
		"filter_id":               {"0123456789abcdef01234567"},
		"filter_mode":             {"intersect"},
		"compact":                 {"true"},
		"view":                    {"enriched"},
		"include_cookie_evidence": {"true"},
		"exclude_cookie_evidence": {"true"},
		"custom[]":                {"value/with space"},
	}
	if !reflect.DeepEqual(gotQuery, wantQuery) {
		t.Fatalf("query = %#v, want %#v", gotQuery, wantQuery)
	}
	for _, encoded := range []string{
		"q=alice%2Btag%40example.com",
		"search_id=sess%2F123",
		"include=credentials",
		"include=evidence",
		"filter=%7B%22field%22%3A%22domain%22",
		"custom%5B%5D=value%2Fwith+space",
	} {
		if !strings.Contains(gotRawQuery, encoded) {
			t.Fatalf("raw query %q does not contain %q", gotRawQuery, encoded)
		}
	}

	assertInvestigationResponseDecoded(t, resp)
}

func TestStealerV2Service_InvestigationSearchPostRequestConstruction(t *testing.T) {
	victimCursor := "victims cursor"
	var gotMethod string
	var gotPath string
	var gotQuery url.Values
	var gotAPIKey string
	var gotBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.EscapedPath()
		gotQuery = r.URL.Query()
		gotAPIKey = r.Header.Get("x-api-key")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("ReadAll() error = %v", err)
		}
		if err := json.Unmarshal(body, &gotBody); err != nil {
			t.Fatalf("Unmarshal request body error = %v; body=%s", err, string(body))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(investigationSearchResponseJSON))
	}))
	defer server.Close()

	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.Stealer.InvestigationSearchPost(&V2InvestigationSearchRequest{
		Q:          "example.com",
		Scope:      "all",
		Include:    []string{"credentials", "victims", "evidence", "files", "related_credentials"},
		FilterMode: "fanout",
		Compact:    true,
		PageSize:   25,
		View:       "enriched",
		SearchID:   "sess_0123456789abcdef",
		Wildcard:   true,
		From:       "2026-01-01T00:00:00Z",
		To:         "2026-02-01T00:00:00Z",
		DateField:  "indexed_at",
		LogID:      "log-1",
		HasLogID:   true,
		Sort:       "-indexed_at",
		Fields:     []string{"id", "log_id"},
		Filter:     StructuredFilter{"field": "domain", "operator": "eq", "value": "example.com"},
		FilterID:   "0123456789abcdef01234567",
		Filters: V2InvestigationSectionFilters{
			"credentials": {
				"domain":     []string{"example.com"},
				"has_log_id": true,
			},
			"evidence": {
				"confidence": []string{"high"},
				"service":    "discord",
			},
		},
		Cursors: V2InvestigationCursors{
			"credentials": nil,
			"victims":     &victimCursor,
		},
		Extra: map[string]interface{}{
			"debug": true,
		},
	})
	if err != nil {
		t.Fatalf("InvestigationSearchPost() error = %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Fatalf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/service/v2/stealer/investigation/search" {
		t.Fatalf("path = %s, want canonical investigation path", gotPath)
	}
	if len(gotQuery) != 0 {
		t.Fatalf("query = %v, want empty", gotQuery)
	}
	if gotAPIKey != "test-key" {
		t.Fatalf("x-api-key = %q, want test-key", gotAPIKey)
	}

	wantBody := map[string]interface{}{
		"q":           "example.com",
		"scope":       "all",
		"include":     []interface{}{"credentials", "victims", "evidence", "files", "related_credentials"},
		"filter_mode": "fanout",
		"compact":     true,
		"page_size":   float64(25),
		"view":        "enriched",
		"search_id":   "sess_0123456789abcdef",
		"wildcard":    true,
		"from":        "2026-01-01T00:00:00Z",
		"to":          "2026-02-01T00:00:00Z",
		"date_field":  "indexed_at",
		"log_id":      "log-1",
		"has_log_id":  true,
		"sort":        "-indexed_at",
		"fields":      []interface{}{"id", "log_id"},
		"filter":      map[string]interface{}{"field": "domain", "operator": "eq", "value": "example.com"},
		"filter_id":   "0123456789abcdef01234567",
		"filters": map[string]interface{}{
			"credentials": map[string]interface{}{
				"domain":     []interface{}{"example.com"},
				"has_log_id": true,
			},
			"evidence": map[string]interface{}{
				"confidence": []interface{}{"high"},
				"service":    "discord",
			},
		},
		"cursors": map[string]interface{}{
			"credentials": nil,
			"victims":     "victims cursor",
		},
		"debug": true,
	}
	if !reflect.DeepEqual(gotBody, wantBody) {
		t.Fatalf("body = %#v, want %#v", gotBody, wantBody)
	}

	assertInvestigationResponseDecoded(t, resp)
}

func TestStealerV2Service_InvestigationAliasPaths(t *testing.T) {
	tests := []struct {
		name       string
		call       func(*Client) error
		wantMethod string
		wantPath   string
		wantQuery  url.Values
		wantBody   map[string]interface{}
	}{
		{
			name: "GET legacy alias",
			call: func(client *Client) error {
				resp, err := client.Stealer.InvestigationSearchAlias("example.com", &InvestigationSearchOptions{Scope: "all"})
				if err != nil {
					return err
				}
				if resp == nil || !resp.Success {
					t.Fatalf("unexpected alias response: %#v", resp)
				}
				return nil
			},
			wantMethod: http.MethodGet,
			wantPath:   "/service/v2/investigate/search",
			wantQuery:  url.Values{"q": {"example.com"}, "scope": {"all"}},
		},
		{
			name: "POST legacy alias",
			call: func(client *Client) error {
				resp, err := client.Stealer.InvestigationSearchAliasPost(&V2InvestigationSearchRequest{Q: "example.com", Scope: "all"})
				if err != nil {
					return err
				}
				if resp == nil || !resp.Success {
					t.Fatalf("unexpected alias post response: %#v", resp)
				}
				return nil
			},
			wantMethod: http.MethodPost,
			wantPath:   "/service/v2/investigate/search",
			wantQuery:  url.Values{},
			wantBody:   map[string]interface{}{"q": "example.com", "scope": "all"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotMethod string
			var gotPath string
			var gotQuery url.Values
			var gotBody map[string]interface{}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotMethod = r.Method
				gotPath = r.URL.EscapedPath()
				gotQuery = r.URL.Query()

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("ReadAll() error = %v", err)
				}
				if tt.wantBody != nil {
					if err := json.Unmarshal(body, &gotBody); err != nil {
						t.Fatalf("Unmarshal request body error = %v; body=%s", err, string(body))
					}
				} else if len(body) > 0 {
					t.Fatalf("request body = %s, want empty", string(body))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"success":true,"data":{"query":"example.com","scope":"all"}}`))
			}))
			defer server.Close()

			client, err := NewClient("test-key", WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			if err := tt.call(client); err != nil {
				t.Fatalf("alias call error = %v", err)
			}
			if gotMethod != tt.wantMethod {
				t.Fatalf("method = %s, want %s", gotMethod, tt.wantMethod)
			}
			if gotPath != tt.wantPath {
				t.Fatalf("path = %s, want %s", gotPath, tt.wantPath)
			}
			if !reflect.DeepEqual(gotQuery, tt.wantQuery) {
				t.Fatalf("query = %v, want %v", gotQuery, tt.wantQuery)
			}
			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Fatalf("body = %#v, want %#v", gotBody, tt.wantBody)
			}
		})
	}
}

func TestStealerV2Service_GetPhonebookRequestConstruction(t *testing.T) {
	var gotMethod string
	var gotPath string
	var gotRawQuery string
	var gotQuery url.Values
	var gotAPIKey string
	var gotBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.EscapedPath()
		gotRawQuery = r.URL.RawQuery
		gotQuery = r.URL.Query()
		gotAPIKey = r.Header.Get("x-api-key")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("ReadAll() error = %v", err)
		}
		gotBody = body

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"domain": "example.com",
			"subdomains": ["accounts.example.com", "mail.example.com"],
			"subdomain_results": [{
				"domain": "accounts.example.com",
				"count": 4,
				"latest_pwned_at": "2026-03-20T12:00:00Z",
				"latest_indexed_at": "2026-03-21T12:00:00Z",
				"redacted": false
			}],
			"emails": [{
				"email": "alice@example.com",
				"count": 3,
				"stealer_count": 2,
				"breach_result_count": 1,
				"breach_count": 1,
				"latest_pwned_at": "2026-03-20T12:00:00Z",
				"latest_indexed_at": "2026-03-21T12:00:00Z",
				"redacted": false
			}],
			"count": 12,
			"email_count": 7,
			"policy_redacted": true,
			"upgrade_required": true,
			"redaction_marker": "upgrade",
			"message": "limited",
			"visible_subdomain_limit": 2,
			"visible_email_limit": 1,
			"redacted_subdomain_count": 10,
			"redacted_email_count": 6
		}`))
	}))
	defer server.Close()

	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.Stealer.GetPhonebook("example.com", &PhonebookOptions{
		Query:    "example.com",
		Alive:    true,
		IsAlive:  true,
		SearchID: "sess/123",
	})
	if err != nil {
		t.Fatalf("GetPhonebook() error = %v", err)
	}

	if gotMethod != http.MethodGet {
		t.Fatalf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/service/v2/phonebook" {
		t.Fatalf("path = %s, want /service/v2/phonebook", gotPath)
	}
	if gotAPIKey != "test-key" {
		t.Fatalf("x-api-key = %q, want test-key", gotAPIKey)
	}
	if len(gotBody) != 0 {
		t.Fatalf("GET body = %s, want empty", string(gotBody))
	}
	wantQuery := url.Values{
		"domain":    {"example.com"},
		"q":         {"example.com"},
		"alive":     {"true"},
		"is_alive":  {"true"},
		"search_id": {"sess/123"},
	}
	if !reflect.DeepEqual(gotQuery, wantQuery) {
		t.Fatalf("query = %#v, want %#v", gotQuery, wantQuery)
	}
	for _, encoded := range []string{"domain=example.com", "q=example.com", "search_id=sess%2F123"} {
		if !strings.Contains(gotRawQuery, encoded) {
			t.Fatalf("raw query %q does not contain %q", gotRawQuery, encoded)
		}
	}

	if resp == nil || resp.Domain != "example.com" || resp.Count != 12 || resp.EmailCount != 7 {
		t.Fatalf("unexpected phonebook response: %#v", resp)
	}
	if len(resp.SubdomainResults) != 1 || resp.SubdomainResults[0].Domain != "accounts.example.com" || resp.SubdomainResults[0].Count != 4 {
		t.Fatalf("unexpected subdomain results: %#v", resp.SubdomainResults)
	}
	if len(resp.Emails) != 1 || resp.Emails[0].Email != "alice@example.com" || resp.Emails[0].StealerCount != 2 || resp.Emails[0].BreachResultCount != 1 {
		t.Fatalf("unexpected emails: %#v", resp.Emails)
	}
	if resp.VisibleSubdomainLimit == nil || *resp.VisibleSubdomainLimit != 2 || resp.VisibleEmailLimit == nil || *resp.VisibleEmailLimit != 1 {
		t.Fatalf("unexpected visible limits: %#v %#v", resp.VisibleSubdomainLimit, resp.VisibleEmailLimit)
	}
	if !resp.PolicyRedacted || !resp.UpgradeRequired || resp.RedactionMarker != "upgrade" || resp.RedactedEmailCount != 6 {
		t.Fatalf("unexpected redaction metadata: %#v", resp)
	}
}

func assertInvestigationResponseDecoded(t *testing.T, resp *V2InvestigationSearchResponse) {
	t.Helper()

	if resp == nil || !resp.Success || resp.Data == nil {
		t.Fatalf("unexpected response: %#v", resp)
	}
	if resp.Data.Query != "example.com" || resp.Data.Scope != "all" {
		t.Fatalf("unexpected query/scope: %#v", resp.Data)
	}
	if resp.Data.Sections == nil || resp.Data.Sections.Credentials == nil || len(resp.Data.Sections.Credentials.Items) != 1 {
		t.Fatalf("unexpected credential section: %#v", resp.Data.Sections)
	}
	credential := resp.Data.Sections.Credentials.Items[0]
	if credential.URLStr != "https://accounts.example.com/login" || credential.SourceType != "stealer" || credential.CanonicalCredentialID != "canonical-1" {
		t.Fatalf("credential did not decode investigation fields: %#v", credential)
	}
	if !credential.IsRelatedCredential || len(credential.InvestigationLabels) != 1 || credential.InvestigationLabels[0] != "direct" {
		t.Fatalf("credential investigation labels = %#v", credential)
	}
	if resp.Data.Sections.Credentials.CredentialStats == nil || resp.Data.Sections.Credentials.CredentialStats.Total != 3 {
		t.Fatalf("credential stats = %#v", resp.Data.Sections.Credentials.CredentialStats)
	}
	if len(resp.Data.Sections.Credentials.Victims) != 1 || resp.Data.Sections.Credentials.Victims[0].LogID != "log-1" {
		t.Fatalf("linked victims = %#v", resp.Data.Sections.Credentials.Victims)
	}
	if resp.Data.Sections.Victims == nil || len(resp.Data.Sections.Victims.Items) != 1 {
		t.Fatalf("unexpected victims section: %#v", resp.Data.Sections.Victims)
	}
	victim := resp.Data.Sections.Victims.Items[0]
	if victim.ServiceCount != 1 || victim.IdentityState != "linked" || victim.DeviceOS != "Windows" || victim.GeoCity != "Austin" {
		t.Fatalf("victim investigation fields did not decode: %#v", victim)
	}
	if resp.Data.Evidence == nil || len(resp.Data.Evidence.Items) != 1 || resp.Data.Evidence.Items[0].PropertyID != "prop-1" {
		t.Fatalf("unexpected evidence: %#v", resp.Data.Evidence)
	}
	if resp.Data.Files == nil || len(resp.Data.Files.Items) != 1 || resp.Data.Files.Items[0].FileID != "file-1" {
		t.Fatalf("unexpected files: %#v", resp.Data.Files)
	}
	if len(resp.Data.Links) != 1 || resp.Data.Links[0].Source == nil || resp.Data.Links[0].Source.Section != "credentials" {
		t.Fatalf("unexpected links: %#v", resp.Data.Links)
	}
	if len(resp.Data.Relations) != 1 || resp.Data.Relations[0].Type != "same_log" {
		t.Fatalf("unexpected relations: %#v", resp.Data.Relations)
	}
	if resp.Data.SectionErrors["files"].Code != "plan" {
		t.Fatalf("unexpected section errors: %#v", resp.Data.SectionErrors)
	}
	if resp.Data.Intersection == nil || !resp.Data.Intersection.Applied || resp.Data.Intersection.Constraints["credentials"] != 1 {
		t.Fatalf("unexpected intersection: %#v", resp.Data.Intersection)
	}
	if !resp.Data.PolicyRedacted || !resp.Data.UpgradeRequired || resp.Data.RedactionMarker != "upgrade" {
		t.Fatalf("unexpected policy metadata: %#v", resp.Data)
	}
}

const investigationSearchResponseJSON = `{
	"success": true,
	"message": "Investigation completed",
	"data": {
		"query": "example.com",
		"scope": "all",
		"sections": {
			"credentials": {
				"items": [{
					"id": "cred-1",
					"log_id": "log-1",
					"url_str": "https://accounts.example.com/login",
					"domain": ["example.com"],
					"source_type": "stealer",
					"password_hash": "hash",
					"archive_hash": "archive-1",
					"canonical_credential_id": "canonical-1",
					"investigation_labels": ["direct"],
					"investigation_link_types": ["same_log"],
					"is_related_credential": true,
					"indexed_at": "2026-03-21T12:00:00Z"
				}],
				"meta": {"total": 1, "count": 1, "has_more": false, "filter_id": "filter-1"},
				"next_cursor": "cred-cursor",
				"victims": [{"log_id": "log-1", "device_users": ["alice"]}],
				"victims_meta": {"total": 1, "count": 1},
				"victims_next_cursor": "victim-cursor",
				"credential_stats": {
					"direct_total": 1,
					"direct_loaded": 1,
					"linked_total": 2,
					"linked_loaded": 1,
					"total": 3,
					"loaded": 2
				}
			},
			"victims": {
				"items": [{
					"log_id": "log-1",
					"device_users": ["alice"],
					"phone_numbers": ["+15551234567"],
					"steam_ids": ["steam-1"],
					"steam_names": ["alice"],
					"services": ["discord"],
					"service_count": 1,
					"identity_state": "linked",
					"device_os": "Windows",
					"device_country": "US",
					"device_city": "Austin",
					"infection_path": "C:\\Users\\alice",
					"antivirus": ["Defender"],
					"domains": ["example.com"],
					"subdomains": ["accounts.example.com"],
					"email_domains": ["example.com"],
					"victim_ip": "203.0.113.5",
					"geo_country": "US",
					"geo_city": "Austin"
				}],
				"meta": {"total": 1, "count": 1}
			}
		},
		"evidence": {
			"items": [{
				"log_id": "log-1",
				"property_id": "prop-1",
				"property_type": "identity",
				"service": "discord",
				"identity_kind": "account",
				"account_id": "123",
				"username": "alice",
				"display_name": "Alice",
				"value": "alice#0001",
				"domain": "discord.com",
				"active": true,
				"source_type": "cookie",
				"source_path": "Cookies.txt",
				"source_file_id": "file-1",
				"confidence": 0.98,
				"confidence_label": "high",
				"confidence_score": 98,
				"indexed_at": "2026-03-21T12:00:00Z"
			}],
			"meta": {"total": 1, "count": 1},
			"policy_redacted": true,
			"redaction_marker": "upgrade"
		},
		"files": {
			"items": [{
				"log_id": "log-1",
				"file_id": "file-1",
				"name": "Cookies.txt",
				"folder": "Browser",
				"path": "Browser/Cookies.txt",
				"ext": "txt",
				"kind": "cookies",
				"size_bytes": 1234
			}],
			"meta": {"count": 1}
		},
		"related_credentials": {
			"items": [],
			"meta": {"count": 0}
		},
		"links": [{
			"relation_type": "same_log",
			"display_label": "Credential belongs to victim log",
			"log_id": "log-1",
			"source": {"section": "credentials", "id": "cred-1"},
			"target": {"section": "victims", "log_id": "log-1"},
			"matched_field": "log_id",
			"matched_value_label": "log-1",
			"credential_id": "cred-1",
			"property_id": "prop-1",
			"file_id": "file-1",
			"confidence": 1,
			"reason": "same log"
		}],
		"relations": [{
			"type": "same_log",
			"log_id": "log-1",
			"credential_id": "cred-1",
			"property_id": "prop-1",
			"file_id": "file-1",
			"reason": "same log",
			"evidence": "log_id"
		}],
		"section_errors": {
			"files": {"section": "files", "error": "denied", "code": "plan", "reason": "upgrade"}
		},
		"intersection": {
			"mode": "intersect",
			"applied": true,
			"constraints": {"credentials": 1},
			"candidate_cap": 100,
			"truncated": false
		},
		"policy_redacted": true,
		"upgrade_required": true,
		"redaction_marker": "upgrade"
	}
}`
