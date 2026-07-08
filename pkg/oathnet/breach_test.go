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

func TestBreachV2Service_SearchRequestConstruction(t *testing.T) {
	filter := StructuredFilter{
		"and": []interface{}{
			map[string]interface{}{
				"field":    "country",
				"operator": "eq",
				"value":    "US",
			},
		},
	}

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
			"success": true,
			"message": "ok",
			"data": {
				"items": [{
					"id": "breach-row-1",
					"email": "alice@example.com",
					"email_domain": "example.com",
					"username": "alice",
					"password": "secret",
					"password_hash": "hash",
					"salt": "salt",
					"full_name": "Alice Example",
					"first_name": "Alice",
					"last_name": "Example",
					"middle_name": "Q",
					"display_name": "A. Example",
					"phone_number": "+15551234567",
					"phone_national": "5551234567",
					"address_street": "1 Test St",
					"city": "Austin",
					"state": "TX",
					"postal_code": "78701",
					"country": "US",
					"date_birth": "1990-01-01",
					"age": 36,
					"created_at": "2026-01-01T00:00:00Z",
					"last_login": "2026-01-02T00:00:00Z",
					"indexed_at": "2026-01-03T00:00:00Z",
					"ip": "192.0.2.10",
					"discordid": "1234567890",
					"instagram": "alicegram",
					"linkedin": "alice-link",
					"iban": "DE89370400440532013000",
					"ssn": "123-45-6789",
					"dbname": "linkedin.com",
					"gender": "f",
					"language": "en",
					"bio": "hello",
					"location": "Austin, TX",
					"extra": {"source": "fixture"}
				}],
				"meta": {
					"total": 42,
					"count": 1,
					"took_ms": 12,
					"has_more": true,
					"total_pages": 2,
					"max_score": 9.5,
					"filter_id": "0123456789abcdef01234567"
				},
				"dbname_info": {
					"linkedin.com": {
						"Title": "LinkedIn",
						"Domain": "linkedin.com",
						"BreachDate": "2021-06-01",
						"PwnCount": 700000000,
						"Description": "Public test metadata"
					}
				},
				"next_cursor": "next cursor",
				"_meta": {"trace_id": "trace-1"}
			}
		}`))
	}))
	defer server.Close()

	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	result, err := client.Breach.Search("alice+tag@example.com", &BreachV2SearchOptions{
		Cursor:         "cursor/with space",
		PageSize:       25,
		Sort:           "-indexed_at",
		From:           "2026-01-01T00:00:00Z",
		To:             "2026-02-01T00:00:00Z",
		DateField:      "indexed_at",
		Wildcard:       true,
		Logic:          "and",
		Filter:         filter,
		FilterID:       "0123456789abcdef01234567",
		Email:          []string{"alice+tag@example.com", "bob@example.com"},
		EmailDomains:   []string{"example.com"},
		Domains:        []string{"app.example.com"},
		Usernames:      []string{"alice"},
		Passwords:      []string{"secret"},
		PasswordHashes: []string{"hash"},
		IPs:            []string{"192.0.2.10"},
		Phones:         []string{"+1 555 123 4567"},
		FirstNames:     []string{"Alice"},
		LastNames:      []string{"Example"},
		FullNames:      []string{"Alice Example"},
		Cities:         []string{"Austin"},
		Countries:      []string{"US"},
		States:         []string{"TX"},
		PostalCodes:    []string{"78701"},
		DBNames:        []string{"linkedin.com", "twitter.com"},
		DiscordIDs:     []string{"1234567890"},
		IBANs:          []string{"DE89370400440532013000"},
		SSNs:           []string{"123-45-6789"},
		DateBirthFrom:  "1980-01-01",
		DateBirthTo:    "2000-12-31",
		Names:          []string{"Alice E"},
		Genders:        []string{"f"},
		Addresses:      []string{"1 Test St"},
		Discord:        []string{"alice#0001"},
		Social:         []string{"alice-social"},
		Financial:      []string{"alice-financial"},
		Gaming:         []string{"alice-game"},
		Fields:         []string{"email", "dbname"},
		SearchID:       "search 123",
		ExtraQuery: map[string][]string{
			"steam[]": {"steam/with space"},
		},
	})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Fatalf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/service/v2/breach/search" {
		t.Fatalf("path = %s, want /service/v2/breach/search", gotPath)
	}
	if gotAPIKey != "test-key" {
		t.Fatalf("x-api-key = %q, want test-key", gotAPIKey)
	}
	var gotBodyJSON map[string]interface{}
	if err := json.Unmarshal(gotBody, &gotBodyJSON); err != nil {
		t.Fatalf("Unmarshal request body error = %v; body=%s", err, string(gotBody))
	}
	wantBody := map[string]interface{}{
		"q": "alice+tag@example.com",
		"filter": map[string]interface{}{
			"and": []interface{}{
				map[string]interface{}{
					"field":    "country",
					"operator": "eq",
					"value":    "US",
				},
			},
		},
		"filter_id": "0123456789abcdef01234567",
	}
	if !reflect.DeepEqual(gotBodyJSON, wantBody) {
		t.Fatalf("body = %#v, want %#v", gotBodyJSON, wantBody)
	}

	wantQuery := url.Values{
		"cursor":          {"cursor/with space"},
		"page_size":       {"25"},
		"sort":            {"-indexed_at"},
		"from":            {"2026-01-01T00:00:00Z"},
		"to":              {"2026-02-01T00:00:00Z"},
		"date_field":      {"indexed_at"},
		"wildcard":        {"true"},
		"logic":           {"and"},
		"email[]":         {"alice+tag@example.com", "bob@example.com"},
		"email_domain[]":  {"example.com"},
		"domain[]":        {"app.example.com"},
		"username[]":      {"alice"},
		"password[]":      {"secret"},
		"password_hash[]": {"hash"},
		"ip[]":            {"192.0.2.10"},
		"phone[]":         {"+1 555 123 4567"},
		"first_name[]":    {"Alice"},
		"last_name[]":     {"Example"},
		"full_name[]":     {"Alice Example"},
		"city[]":          {"Austin"},
		"country[]":       {"US"},
		"state[]":         {"TX"},
		"postal_code[]":   {"78701"},
		"dbname[]":        {"linkedin.com", "twitter.com"},
		"discord_id[]":    {"1234567890"},
		"iban[]":          {"DE89370400440532013000"},
		"ssn[]":           {"123-45-6789"},
		"date_birth_from": {"1980-01-01"},
		"date_birth_to":   {"2000-12-31"},
		"name[]":          {"Alice E"},
		"gender[]":        {"f"},
		"address[]":       {"1 Test St"},
		"discord[]":       {"alice#0001"},
		"social[]":        {"alice-social"},
		"financial[]":     {"alice-financial"},
		"gaming[]":        {"alice-game"},
		"fields[]":        {"email", "dbname"},
		"search_id":       {"search 123"},
		"steam[]":         {"steam/with space"},
	}
	if !reflect.DeepEqual(gotQuery, wantQuery) {
		t.Fatalf("query = %#v, want %#v", gotQuery, wantQuery)
	}
	for _, encoded := range []string{
		"cursor=cursor%2Fwith+space",
		"email%5B%5D=alice%2Btag%40example.com",
		"fields%5B%5D=email",
		"steam%5B%5D=steam%2Fwith+space",
	} {
		if !strings.Contains(gotRawQuery, encoded) {
			t.Fatalf("raw query %q does not contain %q", gotRawQuery, encoded)
		}
	}

	if result == nil || !result.Success || result.Data == nil {
		t.Fatalf("unexpected response: %#v", result)
	}
	if result.Data.Meta == nil || result.Data.Meta.FilterID != "0123456789abcdef01234567" || !result.Data.Meta.HasMore {
		t.Fatalf("unexpected metadata: %#v", result.Data.Meta)
	}
	if result.Data.NextCursor != "next cursor" {
		t.Fatalf("next cursor = %q, want next cursor", result.Data.NextCursor)
	}
	if len(result.Data.Items) != 1 || result.Data.Items[0].Email != "alice@example.com" || result.Data.Items[0].DBName != "linkedin.com" {
		t.Fatalf("unexpected items: %#v", result.Data.Items)
	}
	if result.Data.DBNameInfo["linkedin.com"].PwnCount != 700000000 {
		t.Fatalf("unexpected dbname_info: %#v", result.Data.DBNameInfo)
	}
	if result.Data.APIMeta["trace_id"] != "trace-1" {
		t.Fatalf("unexpected _meta: %#v", result.Data.APIMeta)
	}
}

func TestBreachV2Service_SearchPostRequestConstruction(t *testing.T) {
	var gotMethod string
	var gotPath string
	var gotRawQuery string
	var gotQuery url.Values
	var gotBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.EscapedPath()
		gotRawQuery = r.URL.RawQuery
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
		_, _ = w.Write([]byte(`{"success":true,"data":{"items":[],"meta":{"total":0,"count":0,"took_ms":3,"filter_id":"abcdefabcdefabcdefabcdef"},"next_cursor":""}}`))
	}))
	defer server.Close()

	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	result, err := client.Breach.SearchPost(V2SearchPostBody{
		"filter_id": "abcdefabcdefabcdefabcdef",
		"filter": map[string]interface{}{
			"field":    "email_domain",
			"operator": "eq",
			"value":    "example.com",
		},
	}, &BreachV2SearchOptions{
		Cursor:    "post cursor",
		PageSize:  50,
		Sort:      "-pwned_at",
		From:      "2026-03-01T00:00:00Z",
		To:        "2026-04-01T00:00:00Z",
		DateField: "pwned_at",
		Fields:    []string{"email", "email_domain"},
		SearchID:  "search/post",
		Filter: StructuredFilter{
			"field":    "country",
			"operator": "eq",
			"value":    "CA",
		},
		FilterID: "should-not-be-query",
	})
	if err != nil {
		t.Fatalf("SearchPost() error = %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Fatalf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/service/v2/breach/search" {
		t.Fatalf("path = %s, want /service/v2/breach/search", gotPath)
	}
	wantQuery := url.Values{
		"cursor":     {"post cursor"},
		"page_size":  {"50"},
		"sort":       {"-pwned_at"},
		"from":       {"2026-03-01T00:00:00Z"},
		"to":         {"2026-04-01T00:00:00Z"},
		"date_field": {"pwned_at"},
		"fields[]":   {"email", "email_domain"},
		"search_id":  {"search/post"},
	}
	if !reflect.DeepEqual(gotQuery, wantQuery) {
		t.Fatalf("query = %#v, want %#v", gotQuery, wantQuery)
	}
	if gotQuery.Get("filter") != "" || gotQuery.Get("filter_id") != "" {
		t.Fatalf("POST query unexpectedly included structured filter params: %q", gotRawQuery)
	}
	if !strings.Contains(gotRawQuery, "search_id=search%2Fpost") || !strings.Contains(gotRawQuery, "fields%5B%5D=email_domain") {
		t.Fatalf("raw query %q did not preserve escaping for search_id/fields", gotRawQuery)
	}

	wantBody := map[string]interface{}{
		"filter_id": "abcdefabcdefabcdefabcdef",
		"filter": map[string]interface{}{
			"field":    "email_domain",
			"operator": "eq",
			"value":    "example.com",
		},
	}
	if !reflect.DeepEqual(gotBody, wantBody) {
		t.Fatalf("body = %#v, want %#v", gotBody, wantBody)
	}
	if result == nil || !result.Success || result.Data == nil || result.Data.Meta.FilterID != "abcdefabcdefabcdefabcdef" {
		t.Fatalf("unexpected response: %#v", result)
	}
}

func TestBreachV2Service_AutocompleteRequestConstruction(t *testing.T) {
	tests := []struct {
		name      string
		call      func(*Client) error
		wantPath  string
		wantQuery url.Values
		response  string
	}{
		{
			name: "value autocomplete sends field query limit and include_info",
			call: func(client *Client) error {
				resp, err := client.Breach.Autocomplete(&BreachV2AutocompleteOptions{
					Field:       "email_domain",
					Query:       "g mail",
					Limit:       7,
					IncludeInfo: true,
				})
				if err != nil {
					return err
				}
				if resp == nil || resp.TookMs != 2 || len(resp.Items) != 1 || resp.Items[0].Value != "gmail.com" {
					t.Fatalf("unexpected autocomplete response: %#v", resp)
				}
				return nil
			},
			wantPath: "/service/v2/breach/autocomplete",
			wantQuery: url.Values{
				"field":        {"email_domain"},
				"q":            {"g mail"},
				"limit":        {"7"},
				"include_info": {"true"},
			},
			response: `{"items":[{"field":"email_domain","value":"gmail.com","count":1200}],"took_ms":2}`,
		},
		{
			name: "DB-name autocomplete sends q and limit",
			call: func(client *Client) error {
				resp, err := client.Breach.AutocompleteDBNames(&BreachV2AutocompleteDBNamesOptions{
					Query: "linked in",
					Limit: 3,
				})
				if err != nil {
					return err
				}
				if resp == nil || resp.TookMs != 4 || len(resp.Items) != 1 || resp.Items[0].Info == nil || resp.Items[0].Info.Title != "LinkedIn" {
					t.Fatalf("unexpected DB-name autocomplete response: %#v", resp)
				}
				return nil
			},
			wantPath: "/service/v2/breach/autocomplete/dbnames",
			wantQuery: url.Values{
				"q":     {"linked in"},
				"limit": {"3"},
			},
			response: `{"items":[{"name":"linkedin.com","count":700000000,"fields":["email","password"],"info":{"Title":"LinkedIn","Domain":"linkedin.com","BreachDate":"2021-06-01","PwnCount":700000000,"Description":"Public test metadata"}}],"took_ms":4}`,
		},
		{
			name: "field coverage autocomplete sends field and limit",
			call: func(client *Client) error {
				resp, err := client.Breach.AutocompleteFields(&BreachV2AutocompleteFieldsOptions{
					Field: "discord/id",
					Limit: 5,
				})
				if err != nil {
					return err
				}
				if resp == nil || resp.Field != "discord/id" || resp.Total != 1 || len(resp.Items) != 1 || resp.Items[0].DBName != "gaming.example" {
					t.Fatalf("unexpected field autocomplete response: %#v", resp)
				}
				return nil
			},
			wantPath: "/service/v2/breach/autocomplete/fields",
			wantQuery: url.Values{
				"field": {"discord/id"},
				"limit": {"5"},
			},
			response: `{"field":"discord/id","items":[{"dbname":"gaming.example","count":99}],"total":1,"took_ms":5}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotMethod string
			var gotPath string
			var gotRawQuery string
			var gotQuery url.Values

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotMethod = r.Method
				gotPath = r.URL.EscapedPath()
				gotRawQuery = r.URL.RawQuery
				gotQuery = r.URL.Query()

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client, err := NewClient("test-key", WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			if err := tt.call(client); err != nil {
				t.Fatalf("autocomplete call error = %v", err)
			}
			if gotMethod != http.MethodGet {
				t.Fatalf("method = %s, want GET", gotMethod)
			}
			if gotPath != tt.wantPath {
				t.Fatalf("path = %s, want %s", gotPath, tt.wantPath)
			}
			if !reflect.DeepEqual(gotQuery, tt.wantQuery) {
				t.Fatalf("query = %#v, want %#v", gotQuery, tt.wantQuery)
			}
			if strings.Contains(gotRawQuery, " ") || strings.Contains(gotRawQuery, "/") {
				t.Fatalf("raw query contains unescaped space or slash: %q", gotRawQuery)
			}
		})
	}
}
