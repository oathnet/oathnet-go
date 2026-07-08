package oathnet

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestStealerV2Service_SearchPostRequestConstruction(t *testing.T) {
	var gotMethod string
	var gotPath string
	var gotQuery url.Values
	var gotBody map[string]interface{}
	var gotAPIKey string

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
		_, _ = w.Write([]byte(`{
			"success": true,
			"data": {
				"items": [{
					"id": "credential-1",
					"log_id": "log-1",
					"url_str": "https://accounts.example.com/login",
					"domain": ["example.com"],
					"subdomain": ["accounts.example.com"],
					"username": "alice",
					"password": "secret",
					"archive_hash": "archive-hash"
				}],
				"meta": {"total": 1, "count": 1, "took_ms": 5, "has_more": true},
				"next_cursor": "next-stealer"
			}
		}`))
	}))
	defer server.Close()

	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.Stealer.SearchStealerV2Post(V2SearchPostBody{
		"filter_id": "0123456789abcdef01234567",
		"filter": map[string]interface{}{
			"field":    "domain",
			"operator": "eq",
			"value":    "example.com",
		},
	}, &StealerSearchOptions{
		Cursor:    "cursor/with space",
		PageSize:  50,
		Sort:      "-indexed_at",
		From:      "2026-01-01T00:00:00Z",
		To:        "2026-02-01T00:00:00Z",
		DateField: "indexed_at",
		Fields:    []string{"id", "log_id"},
		SearchID:  "sess/123",
		View:      "enriched",
	})
	if err != nil {
		t.Fatalf("SearchStealerV2Post() error = %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Fatalf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/service/v2/stealer/search" {
		t.Fatalf("path = %s, want /service/v2/stealer/search", gotPath)
	}
	if gotAPIKey != "test-key" {
		t.Fatalf("x-api-key = %q, want test-key", gotAPIKey)
	}
	wantQuery := url.Values{
		"cursor":     {"cursor/with space"},
		"page_size":  {"50"},
		"sort":       {"-indexed_at"},
		"from":       {"2026-01-01T00:00:00Z"},
		"to":         {"2026-02-01T00:00:00Z"},
		"date_field": {"indexed_at"},
		"fields[]":   {"id", "log_id"},
		"search_id":  {"sess/123"},
		"view":       {"enriched"},
	}
	if !reflect.DeepEqual(gotQuery, wantQuery) {
		t.Fatalf("query = %#v, want %#v", gotQuery, wantQuery)
	}
	wantBody := map[string]interface{}{
		"filter_id": "0123456789abcdef01234567",
		"filter": map[string]interface{}{
			"field":    "domain",
			"operator": "eq",
			"value":    "example.com",
		},
	}
	if !reflect.DeepEqual(gotBody, wantBody) {
		t.Fatalf("body = %#v, want %#v", gotBody, wantBody)
	}
	if resp == nil || !resp.Success || resp.Data == nil || resp.Data.NextCursor != "next-stealer" {
		t.Fatalf("unexpected response: %#v", resp)
	}
	if len(resp.Data.Items) != 1 || resp.Data.Items[0].LogID != "log-1" || resp.Data.Items[0].ArchiveHash != "archive-hash" {
		t.Fatalf("unexpected items: %#v", resp.Data.Items)
	}
}

func TestVictimsService_SearchRequestConstruction(t *testing.T) {
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
		if err := json.Unmarshal(body, &gotBody); err != nil {
			t.Fatalf("Unmarshal request body error = %v; body=%s", err, string(body))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"success": true,
			"data": {
				"items": [{
					"log_id": "victim-1",
					"device_users": ["alice"],
					"device_ips": ["192.0.2.10"],
					"discord_ids": ["1234567890"],
					"total_docs": 12,
					"services": ["discord"],
					"service_count": 1,
					"device_country": "US"
				}],
				"meta": {"total": 1, "count": 1, "took_ms": 8, "filter_id": "0123456789abcdef01234567"},
				"next_cursor": "next-victim"
			}
		}`))
	}))
	defer server.Close()

	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.Victims.Search("example.com", &VictimsSearchOptions{
		Cursor:          "victim cursor",
		PageSize:        50,
		Sort:            "-indexed_at",
		From:            "2026-01-01T00:00:00Z",
		To:              "2026-01-31T00:00:00Z",
		DateField:       "indexed_at",
		Wildcard:        true,
		LogID:           "log-123",
		Filter:          map[string]interface{}{"field": "service", "operator": "eq", "value": "discord"},
		FilterID:        "0123456789abcdef01234567",
		TotalDocsMin:    2,
		TotalDocsMax:    20,
		ServiceCountMin: 1,
		ServiceCountMax: 8,
		Emails:          []string{"alice@example.com"},
		EmailDomains:    []string{"example.com"},
		IPs:             []string{"203.0.113.11"},
		HWIDs:           []string{"hwid-1"},
		DiscordIDs:      []string{"1234567890"},
		Usernames:       []string{"alice"},
		Countries:       []string{"US"},
		Cities:          []string{"New York"},
		OSes:            []string{"Windows"},
		Services:        []string{"discord"},
		SteamIDs:        []string{"steam-1"},
		SteamNames:      []string{"AliceSteam"},
		Phones:          []string{"+15551234567"},
		Domains:         []string{"example.com"},
		Subdomains:      []string{"accounts.example.com"},
		IdentityStates:  []string{"active"},
		VictimIPs:       []string{"198.51.100.4"},
		Antivirus:       []string{"defender"},
		InfectionPaths:  []string{"C:/Users/Alice"},
		Fields:          []string{"log_id", "device_users"},
		View:            "enriched",
		SearchID:        "session-123",
	})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Fatalf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/service/v2/victims/search" {
		t.Fatalf("path = %s, want /service/v2/victims/search", gotPath)
	}
	wantQuery := url.Values{
		"cursor":            {"victim cursor"},
		"page_size":         {"50"},
		"sort":              {"-indexed_at"},
		"from":              {"2026-01-01T00:00:00Z"},
		"to":                {"2026-01-31T00:00:00Z"},
		"date_field":        {"indexed_at"},
		"wildcard":          {"true"},
		"log_id":            {"log-123"},
		"total_docs_min":    {"2"},
		"total_docs_max":    {"20"},
		"service_count_min": {"1"},
		"service_count_max": {"8"},
		"email[]":           {"alice@example.com"},
		"email_domain[]":    {"example.com"},
		"ip[]":              {"203.0.113.11"},
		"hwid[]":            {"hwid-1"},
		"discord_id[]":      {"1234567890"},
		"username[]":        {"alice"},
		"country[]":         {"US"},
		"city[]":            {"New York"},
		"os[]":              {"Windows"},
		"service[]":         {"discord"},
		"steam_id[]":        {"steam-1"},
		"steam_name[]":      {"AliceSteam"},
		"phone[]":           {"+15551234567"},
		"domain[]":          {"example.com"},
		"subdomain[]":       {"accounts.example.com"},
		"identity_state[]":  {"active"},
		"victim_ip[]":       {"198.51.100.4"},
		"antivirus[]":       {"defender"},
		"infection_path[]":  {"C:/Users/Alice"},
		"fields[]":          {"log_id", "device_users"},
		"view":              {"enriched"},
		"search_id":         {"session-123"},
	}
	if !reflect.DeepEqual(gotQuery, wantQuery) {
		t.Fatalf("query = %#v, want %#v", gotQuery, wantQuery)
	}
	wantBody := map[string]interface{}{
		"q":         "example.com",
		"filter":    map[string]interface{}{"field": "service", "operator": "eq", "value": "discord"},
		"filter_id": "0123456789abcdef01234567",
	}
	if !reflect.DeepEqual(gotBody, wantBody) {
		t.Fatalf("body = %#v, want %#v", gotBody, wantBody)
	}
	if resp == nil || !resp.Success || resp.Data == nil || resp.Data.NextCursor != "next-victim" {
		t.Fatalf("unexpected response: %#v", resp)
	}
	if len(resp.Data.Items) != 1 || resp.Data.Items[0].LogID != "victim-1" || resp.Data.Items[0].DeviceCountry != "US" {
		t.Fatalf("unexpected items: %#v", resp.Data.Items)
	}
}

func TestVictimsService_SearchPostRequestConstruction(t *testing.T) {
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
		if err := json.Unmarshal(body, &gotBody); err != nil {
			t.Fatalf("Unmarshal request body error = %v; body=%s", err, string(body))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"success": true,
			"data": {
				"items": [{
					"log_id": "victim-1",
					"device_users": ["alice"],
					"device_ips": ["192.0.2.10"],
					"discord_ids": ["1234567890"],
					"total_docs": 12,
					"services": ["discord"],
					"service_count": 1,
					"device_country": "US"
				}],
				"meta": {"total": 1, "count": 1, "took_ms": 8},
				"next_cursor": "next-victim"
			}
		}`))
	}))
	defer server.Close()

	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.Victims.SearchVictimsV2Post(V2SearchPostBody{
		"filter": map[string]interface{}{
			"and": []interface{}{
				map[string]interface{}{"field": "service", "operator": "eq", "value": "discord"},
			},
		},
	}, &VictimsSearchOptions{
		Cursor:    "victim cursor",
		PageSize:  25,
		Sort:      "-pwned_at",
		From:      "2026-03-01T00:00:00Z",
		To:        "2026-04-01T00:00:00Z",
		DateField: "pwned_at",
		Fields:    []string{"log_id", "services"},
		SearchID:  "sess-victims",
		View:      "enriched",
	})
	if err != nil {
		t.Fatalf("SearchVictimsV2Post() error = %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Fatalf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/service/v2/victims/search" {
		t.Fatalf("path = %s, want /service/v2/victims/search", gotPath)
	}
	wantQuery := url.Values{
		"cursor":     {"victim cursor"},
		"page_size":  {"25"},
		"sort":       {"-pwned_at"},
		"from":       {"2026-03-01T00:00:00Z"},
		"to":         {"2026-04-01T00:00:00Z"},
		"date_field": {"pwned_at"},
		"fields[]":   {"log_id", "services"},
		"search_id":  {"sess-victims"},
		"view":       {"enriched"},
	}
	if !reflect.DeepEqual(gotQuery, wantQuery) {
		t.Fatalf("query = %#v, want %#v", gotQuery, wantQuery)
	}
	if gotBody["filter"] == nil {
		t.Fatalf("body = %#v, want filter", gotBody)
	}
	if resp == nil || !resp.Success || resp.Data == nil || resp.Data.NextCursor != "next-victim" {
		t.Fatalf("unexpected response: %#v", resp)
	}
	if len(resp.Data.Items) != 1 || resp.Data.Items[0].LogID != "victim-1" || resp.Data.Items[0].DeviceCountry != "US" {
		t.Fatalf("unexpected items: %#v", resp.Data.Items)
	}
}

func TestVictimsService_FileMetadataSearchRequestConstruction(t *testing.T) {
	tests := []struct {
		name       string
		call       func(*Client) (*V2FileMetadataSearchResponse, error)
		wantMethod string
		wantQuery  url.Values
		wantBody   map[string]interface{}
	}{
		{
			name: "GET",
			call: func(client *Client) (*V2FileMetadataSearchResponse, error) {
				return client.Victims.SearchFilesMetadataV2(&FileMetadataSearchOptions{
					Query:    "cookie",
					LogID:    "log-1",
					Name:     "Cookies.txt",
					Folder:   "Browser",
					Kind:     "cookies",
					Ext:      "txt",
					SizeMin:  10,
					SizeMax:  2048,
					PageSize: 25,
					Cursor:   "file cursor",
					SearchID: "sess-files",
				})
			},
			wantMethod: http.MethodGet,
			wantQuery: url.Values{
				"q":         {"cookie"},
				"log_id":    {"log-1"},
				"name":      {"Cookies.txt"},
				"folder":    {"Browser"},
				"kind":      {"cookies"},
				"ext":       {"txt"},
				"size_min":  {"10"},
				"size_max":  {"2048"},
				"page_size": {"25"},
				"cursor":    {"file cursor"},
				"search_id": {"sess-files"},
			},
		},
		{
			name: "POST",
			call: func(client *Client) (*V2FileMetadataSearchResponse, error) {
				return client.Victims.SearchFilesMetadataV2Post(&V2FileMetadataSearchRequest{
					Query:    "cookie",
					Kind:     "cookies",
					PageSize: 10,
					SearchID: "sess-files",
				})
			},
			wantMethod: http.MethodPost,
			wantQuery:  url.Values{},
			wantBody: map[string]interface{}{
				"q":         "cookie",
				"kind":      "cookies",
				"page_size": float64(10),
				"search_id": "sess-files",
			},
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
				_, _ = w.Write([]byte(`{
					"success": true,
					"message": "V2-File-Metadata-Search completed successfully",
					"data": {
						"items": [{"log_id":"log-1","file_id":"file-1","name":"Cookies.txt","folder":"Browser","path":"Browser/Cookies.txt","ext":"txt","kind":"cookies","size_bytes":1234}],
						"meta": {"total": 1, "count": 1, "took_ms": 3},
						"next_cursor": "next-file",
						"policy_redacted": true,
						"upgrade_required": true,
						"redaction_marker": "upgrade"
					}
				}`))
			}))
			defer server.Close()

			client, err := NewClient("test-key", WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			resp, err := tt.call(client)
			if err != nil {
				t.Fatalf("file metadata call error = %v", err)
			}
			if gotMethod != tt.wantMethod {
				t.Fatalf("method = %s, want %s", gotMethod, tt.wantMethod)
			}
			if gotPath != "/service/v2/files/search" {
				t.Fatalf("path = %s, want /service/v2/files/search", gotPath)
			}
			if !reflect.DeepEqual(gotQuery, tt.wantQuery) {
				t.Fatalf("query = %#v, want %#v", gotQuery, tt.wantQuery)
			}
			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Fatalf("body = %#v, want %#v", gotBody, tt.wantBody)
			}
			if resp == nil || !resp.Success || resp.Data == nil || resp.Data.NextCursor != "next-file" || !resp.Data.PolicyRedacted || !resp.Data.UpgradeRequired {
				t.Fatalf("unexpected response: %#v", resp)
			}
			if len(resp.Data.Items) != 1 || resp.Data.Items[0].FileID != "file-1" || resp.Data.Items[0].SizeBytes != 1234 {
				t.Fatalf("unexpected items: %#v", resp.Data.Items)
			}
		})
	}
}

func TestVictimsService_PropertiesSearchAndDetailRequestConstruction(t *testing.T) {
	active := false
	confidenceMin := 0.75
	includeCookieEvidence := true
	excludeCookieEvidence := true

	tests := []struct {
		name       string
		call       func(*Client) (*V2VictimPropertiesSearchResponse, error)
		wantMethod string
		wantPath   string
		wantQuery  url.Values
		wantBody   map[string]interface{}
	}{
		{
			name: "global GET",
			call: func(client *Client) (*V2VictimPropertiesSearchResponse, error) {
				return client.Victims.SearchVictimPropertiesV2(&VictimPropertiesSearchOptions{
					Query:                 "example.com",
					PropertyType:          []string{"account", "cookie_domain"},
					Service:               "discord",
					IdentityKind:          "user",
					AccountID:             "account-1",
					Username:              "alice*",
					DisplayName:           "Alice",
					Value:                 "alice@example.com",
					Domain:                "example.com",
					Active:                &active,
					SourceType:            "cookies",
					SourcePath:            "not-a-query-param",
					SourceFileID:          "not-a-query-param",
					Confidence:            []string{"high", "medium"},
					ConfidenceMin:         &confidenceMin,
					IncludeCookieEvidence: &includeCookieEvidence,
					ExcludeCookieEvidence: &excludeCookieEvidence,
					PageSize:              25,
					Cursor:                "properties cursor",
					Sort:                  "-indexed_at",
					SearchID:              "sess-props",
				})
			},
			wantMethod: http.MethodGet,
			wantPath:   "/service/v2/victims/properties/search",
			wantQuery: url.Values{
				"q":                       {"example.com"},
				"property_type":           {"account", "cookie_domain"},
				"service":                 {"discord"},
				"identity_kind":           {"user"},
				"account_id":              {"account-1"},
				"username":                {"alice*"},
				"display_name":            {"Alice"},
				"value":                   {"alice@example.com"},
				"domain":                  {"example.com"},
				"active":                  {"false"},
				"source_type":             {"cookies"},
				"confidence":              {"high", "medium"},
				"confidence_min":          {"0.75"},
				"include_cookie_evidence": {"true"},
				"exclude_cookie_evidence": {"true"},
				"page_size":               {"25"},
				"cursor":                  {"properties cursor"},
				"sort":                    {"-indexed_at"},
				"search_id":               {"sess-props"},
			},
		},
		{
			name: "POST",
			call: func(client *Client) (*V2VictimPropertiesSearchResponse, error) {
				return client.Victims.SearchVictimPropertiesV2Post(&V2VictimPropertiesSearchRequest{
					Query:                 "example.com",
					LogID:                 "log-1",
					PropertyType:          "account",
					Service:               "discord",
					Confidence:            []string{"high"},
					ExcludeCookieEvidence: &excludeCookieEvidence,
					PageSize:              10,
					SearchID:              "sess-props",
				})
			},
			wantMethod: http.MethodPost,
			wantPath:   "/service/v2/victims/properties/search",
			wantQuery:  url.Values{},
			wantBody: map[string]interface{}{
				"q":                       "example.com",
				"log_id":                  "log-1",
				"property_type":           "account",
				"service":                 "discord",
				"confidence":              []interface{}{"high"},
				"exclude_cookie_evidence": true,
				"page_size":               float64(10),
				"search_id":               "sess-props",
			},
		},
		{
			name: "victim GET escapes log ID",
			call: func(client *Client) (*V2VictimPropertiesSearchResponse, error) {
				return client.Victims.GetVictimPropertiesV2("log/with space", &VictimPropertiesSearchOptions{
					Query:        "discord",
					PropertyType: "account",
					Service:      "discord",
					PageSize:     5,
					SearchID:     "sess-one",
				})
			},
			wantMethod: http.MethodGet,
			wantPath:   "/service/v2/victims/log%2Fwith%20space/properties",
			wantQuery: url.Values{
				"q":             {"discord"},
				"property_type": {"account"},
				"service":       {"discord"},
				"page_size":     {"5"},
				"search_id":     {"sess-one"},
			},
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
				_, _ = w.Write([]byte(`{
					"success": true,
					"message": "V2-Victim-Properties-Search completed successfully",
					"data": {
						"items": [{"log_id":"log-1","property_id":"prop-1","property_type":"account","service":"discord","identity_kind":"user","account_id":"account-1","username":"alice","display_name":"Alice","value":"alice@example.com","domain":"example.com","active":false,"source_type":"cookies","source_path":"Cookies.txt","source_file_id":"file-1","confidence":0.99,"confidence_label":"high","confidence_score":99,"indexed_at":"2026-07-06T00:00:00Z"}],
						"meta": {"total": 1, "count": 1, "took_ms": 7},
						"next_cursor": "next-props",
						"policy_redacted": true,
						"upgrade_required": true,
						"redaction_marker": "upgrade"
					}
				}`))
			}))
			defer server.Close()

			client, err := NewClient("test-key", WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			resp, err := tt.call(client)
			if err != nil {
				t.Fatalf("properties call error = %v", err)
			}
			if gotMethod != tt.wantMethod {
				t.Fatalf("method = %s, want %s", gotMethod, tt.wantMethod)
			}
			if gotPath != tt.wantPath {
				t.Fatalf("path = %s, want %s", gotPath, tt.wantPath)
			}
			if !reflect.DeepEqual(gotQuery, tt.wantQuery) {
				t.Fatalf("query = %#v, want %#v", gotQuery, tt.wantQuery)
			}
			if gotQuery.Get("source_path") != "" || gotQuery.Get("source_file_id") != "" {
				t.Fatalf("query unexpectedly included POST-only source fields: %#v", gotQuery)
			}
			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Fatalf("body = %#v, want %#v", gotBody, tt.wantBody)
			}
			if resp == nil || !resp.Success || resp.Data == nil || resp.Data.NextCursor != "next-props" || !resp.Data.PolicyRedacted || !resp.Data.UpgradeRequired {
				t.Fatalf("unexpected response: %#v", resp)
			}
			if len(resp.Data.Items) != 1 || resp.Data.Items[0].PropertyID != "prop-1" || resp.Data.Items[0].ConfidenceLabel != "high" {
				t.Fatalf("unexpected items: %#v", resp.Data.Items)
			}
		})
	}
}

func TestVictimsService_GetPropertiesAcceptsUnwrappedData(t *testing.T) {
	var gotPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.EscapedPath()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"items": [{"log_id":"log-1","property_id":"prop-1","property_type":"account","service":"discord"}],
			"meta": {"total": 1, "count": 1, "took_ms": 7},
			"next_cursor": "next-props",
			"policy_redacted": true
		}`))
	}))
	defer server.Close()

	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.Victims.GetVictimPropertiesV2("log/with space", &VictimPropertiesSearchOptions{PageSize: 5})
	if err != nil {
		t.Fatalf("GetVictimPropertiesV2() error = %v", err)
	}
	if gotPath != "/service/v2/victims/log%2Fwith%20space/properties" {
		t.Fatalf("path = %s, want escaped victim properties path", gotPath)
	}
	if resp == nil || !resp.Success || resp.Data == nil {
		t.Fatalf("unexpected response: %#v", resp)
	}
	if resp.Data.NextCursor != "next-props" || !resp.Data.PolicyRedacted {
		t.Fatalf("unexpected response data: %#v", resp.Data)
	}
	if len(resp.Data.Items) != 1 || resp.Data.Items[0].PropertyID != "prop-1" {
		t.Fatalf("unexpected items: %#v", resp.Data.Items)
	}
}

func TestVictimsService_SummaryCookiesAndDomainRequestConstruction(t *testing.T) {
	includeItems := false

	tests := []struct {
		name        string
		call        func(*Client) error
		response    string
		contentType string
		wantMethod  string
		wantPath    string
		wantQuery   url.Values
		wantAccept  string
	}{
		{
			name: "summary",
			call: func(client *Client) error {
				resp, err := client.Victims.GetVictimSummaryV2("log/with space", &VictimSummaryOptions{SearchID: "sess-summary"})
				if err != nil {
					return err
				}
				if resp == nil || resp.LogID != "log/with space" || resp.Files["total"].(float64) != 7 || len(resp.Warnings) != 1 {
					t.Fatalf("unexpected summary response: %#v", resp)
				}
				return nil
			},
			response:   `{"log_id":"log/with space","generated_at":"2026-07-06T00:00:00Z","stale":true,"victim":{"country":"US"},"assessment":{"risk":"high"},"access":{"plan":"pro"},"targets":{"domains":2},"files":{"total":7},"cookies":{"active":3},"cookie_investigation":{"domains":1},"history":{"seen":1},"cards":{"identity":true},"domains":{"example.com":1},"artifacts":{"cookies":true},"warnings":["cached"],"policy_redacted":true,"upgrade_required":true,"redaction_marker":"upgrade"}`,
			wantMethod: http.MethodGet,
			wantPath:   "/service/v2/victims/log%2Fwith%20space/summary",
			wantQuery:  url.Values{"search_id": {"sess-summary"}},
		},
		{
			name: "cookie inventory",
			call: func(client *Client) error {
				resp, err := client.Victims.GetVictimCookiesV2("log/with space", &VictimCookieInventoryOptions{
					Domain:       "example.com",
					Status:       "active",
					Query:        "session",
					IncludeItems: &includeItems,
					PageSize:     25,
					Cursor:       "cookie cursor",
					SearchID:     "sess-cookies",
				})
				if err != nil {
					return err
				}
				if resp == nil || resp.LogID != "log/with space" || resp.NextCursor != "next-cookie" || len(resp.Items) != 1 || len(resp.Domains) != 1 {
					t.Fatalf("unexpected cookie response: %#v", resp)
				}
				return nil
			},
			response:   `{"log_id":"log/with space","items":[{"domain":"example.com","cookie_domain":".example.com","name":"sid","path":"/","expires_at":"2026-08-01T00:00:00Z","expires_unix":1785542400,"status":"active","session":false,"secure":true,"http_only":true,"source_file_id":"file-1","source_path":"Cookies.txt","line_number":12}],"domains":[{"domain":"example.com","count":1,"active":1,"expired":0,"session":0,"secure":1,"http_only":1,"source_files":1}],"meta":{"total":1,"count":1,"took_ms":4},"next_cursor":"next-cookie","files_scanned":2,"files_matched":1,"truncated":false,"values_redacted":true,"warnings":["values redacted"],"policy_redacted":true,"upgrade_required":true,"redaction_marker":"upgrade"}`,
			wantMethod: http.MethodGet,
			wantPath:   "/service/v2/victims/log%2Fwith%20space/cookies",
			wantQuery: url.Values{
				"domain":        {"example.com"},
				"status":        {"active"},
				"q":             {"session"},
				"include_items": {"false"},
				"page_size":     {"25"},
				"cursor":        {"cookie cursor"},
				"search_id":     {"sess-cookies"},
			},
		},
		{
			name: "inspect cookie domain",
			call: func(client *Client) error {
				text, err := client.Victims.InspectVictimCookieDomainV2("log/with space", &VictimCookieDomainOptions{
					Domain:   "example.com",
					FileID:   "file/with space",
					SearchID: "sess-domain",
				})
				if err != nil {
					return err
				}
				if text != ".example.com\tTRUE\t/\tTRUE\t1785542400\tsid\tredacted\n" {
					t.Fatalf("inspect text = %q", text)
				}
				return nil
			},
			response:    ".example.com\tTRUE\t/\tTRUE\t1785542400\tsid\tredacted\n",
			contentType: "text/plain",
			wantMethod:  http.MethodGet,
			wantPath:    "/service/v2/victims/log%2Fwith%20space/cookies/domain",
			wantQuery: url.Values{
				"domain":    {"example.com"},
				"file_id":   {"file/with space"},
				"search_id": {"sess-domain"},
			},
			wantAccept: "text/plain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotMethod string
			var gotPath string
			var gotQuery url.Values
			var gotAccept string
			var gotBody []byte

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotMethod = r.Method
				gotPath = r.URL.EscapedPath()
				gotQuery = r.URL.Query()
				gotAccept = r.Header.Get("Accept")

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("ReadAll() error = %v", err)
				}
				gotBody = body

				contentType := tt.contentType
				if contentType == "" {
					contentType = "application/json"
				}
				w.Header().Set("Content-Type", contentType)
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client, err := NewClient("test-key", WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			if err := tt.call(client); err != nil {
				t.Fatalf("victim detail call error = %v", err)
			}
			if gotMethod != tt.wantMethod {
				t.Fatalf("method = %s, want %s", gotMethod, tt.wantMethod)
			}
			if gotPath != tt.wantPath {
				t.Fatalf("path = %s, want %s", gotPath, tt.wantPath)
			}
			if !reflect.DeepEqual(gotQuery, tt.wantQuery) {
				t.Fatalf("query = %#v, want %#v", gotQuery, tt.wantQuery)
			}
			if gotAccept != tt.wantAccept {
				t.Fatalf("Accept = %q, want %q", gotAccept, tt.wantAccept)
			}
			if len(gotBody) != 0 {
				t.Fatalf("GET body = %s, want empty", string(gotBody))
			}
		})
	}
}

func TestExportsService_ListRequestConstruction(t *testing.T) {
	var gotMethod string
	var gotPath string
	var gotQuery url.Values

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.EscapedPath()
		gotQuery = r.URL.Query()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"count": 1,
			"next": "https://api.example.test/service/v2/exports/list?page=3&page_size=50",
			"previous": "https://api.example.test/service/v2/exports/list?page=1&page_size=50",
			"results": [{
				"id": "local-1",
				"job_id": "export-1",
				"status": "completed",
				"created_at": "2026-07-06T00:00:00Z",
				"started_at": "2026-07-06T00:00:01Z",
				"completed_at": "2026-07-06T00:00:02Z",
				"expires_at": "2026-07-07T00:00:02Z",
				"progress": {"records_done": 10, "records_total": 10, "bytes_done": 100, "percent": 100, "updated_at": "2026-07-06T00:00:02Z"},
				"result": {"format": "jsonl", "file_name": "export.jsonl", "file_path": "/exports/export.jsonl", "file_size": 100, "records": 10, "ready_at": "2026-07-06T00:00:02Z", "expires_at": "2026-07-07T00:00:02Z"},
				"last_error": "",
				"request": {"type": "docs", "service": "stealer", "format": "jsonl", "limit": 10, "fields": ["email"], "request_count": 1},
				"next_poll_after_ms": 1000,
				"metadata": {"trace_id": "trace-1"}
			}]
		}`))
	}))
	defer server.Close()

	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.Exports.ListExportsV2(2, 50)
	if err != nil {
		t.Fatalf("ListExportsV2() error = %v", err)
	}

	if gotMethod != http.MethodGet {
		t.Fatalf("method = %s, want GET", gotMethod)
	}
	if gotPath != "/service/v2/exports/list" {
		t.Fatalf("path = %s, want /service/v2/exports/list", gotPath)
	}
	wantQuery := url.Values{"page": {"2"}, "page_size": {"50"}}
	if !reflect.DeepEqual(gotQuery, wantQuery) {
		t.Fatalf("query = %#v, want %#v", gotQuery, wantQuery)
	}
	if resp == nil || resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].JobID != "export-1" {
		t.Fatalf("unexpected list response: %#v", resp)
	}
	if resp.Results[0].Result == nil || resp.Results[0].Result.FilePath != "/exports/export.jsonl" || resp.Results[0].Result.ReadyAt == "" {
		t.Fatalf("unexpected export result: %#v", resp.Results[0].Result)
	}
	if resp.Results[0].Request == nil || resp.Results[0].Request.Service != "stealer" || resp.Results[0].Metadata["trace_id"] != "trace-1" {
		t.Fatalf("unexpected export request metadata: %#v", resp.Results[0])
	}
}
