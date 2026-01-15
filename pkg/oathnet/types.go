package oathnet

// ResponseMeta contains metadata about the API response.
type ResponseMeta struct {
	User        *UserMeta        `json:"user,omitempty"`
	Lookups     *LookupsMeta     `json:"lookups,omitempty"`
	Service     *ServiceMeta     `json:"service,omitempty"`
	Performance *PerformanceMeta `json:"performance,omitempty"`
}

type UserMeta struct {
	Plan         string `json:"plan,omitempty"`
	PlanType     string `json:"plan_type,omitempty"`
	IsPlanActive bool   `json:"is_plan_active,omitempty"`
}

type LookupsMeta struct {
	UsedToday   int  `json:"used_today,omitempty"`
	LeftToday   int  `json:"left_today,omitempty"`
	DailyLimit  int  `json:"daily_limit,omitempty"`
	IsUnlimited bool `json:"is_unlimited,omitempty"`
}

type ServiceMeta struct {
	Name         string `json:"name,omitempty"`
	ID           string `json:"id,omitempty"`
	Category     string `json:"category,omitempty"`
	IsPremium    bool   `json:"is_premium,omitempty"`
	IsAvailable  bool   `json:"is_available,omitempty"`
	SessionQuota int    `json:"session_quota,omitempty"`
}

type PerformanceMeta struct {
	DurationMs float64 `json:"duration_ms,omitempty"`
	Timestamp  string  `json:"timestamp,omitempty"`
}

// ============================================
// SEARCH TYPES
// ============================================

type SearchSessionResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message,omitempty"`
	Data    *SearchSessionData `json:"data,omitempty"`
}

type SearchSessionData struct {
	Session *SearchSession     `json:"session"`
	User    *SearchSessionUser `json:"user,omitempty"`
}

type SearchSession struct {
	ID         string `json:"id"`
	Query      string `json:"query"`
	SearchType string `json:"search_type"`
	ExpiresAt  string `json:"expires_at"`
}

type SearchSessionUser struct {
	Plan         string        `json:"plan"`
	PlanType     string        `json:"plan_type"`
	IsPlanActive bool          `json:"is_plan_active"`
	DailyLookups *DailyLookups `json:"daily_lookups,omitempty"`
}

type DailyLookups struct {
	Used        int  `json:"used"`
	Remaining   int  `json:"remaining"`
	Limit       int  `json:"limit"`
	IsUnlimited bool `json:"is_unlimited"`
}

type BreachSearchResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message,omitempty"`
	Data    *BreachSearchData `json:"data,omitempty"`
}

type BreachSearchData struct {
	Results      []BreachResult `json:"results"`
	ResultsFound int            `json:"results_found"`
	ResultsShown int            `json:"results_shown"`
	Cursor       string         `json:"cursor,omitempty"`
	Meta         *ResponseMeta  `json:"_meta,omitempty"`
}

type BreachResult struct {
	ID           string      `json:"id,omitempty"`
	DBName       string      `json:"dbname,omitempty"`
	Email        string      `json:"email,omitempty"`
	Username     interface{} `json:"username,omitempty"` // Can be string or []string
	Password     string      `json:"password,omitempty"`
	PasswordHash string      `json:"password_hash,omitempty"`
	FullName     interface{} `json:"full_name,omitempty"`
	FirstName    interface{} `json:"first_name,omitempty"`
	LastName     interface{} `json:"last_name,omitempty"`
	PhoneNumber  string      `json:"phone_number,omitempty"`
	City         interface{} `json:"city,omitempty"`
	Country      interface{} `json:"country,omitempty"`
	IP           string      `json:"ip,omitempty"`
}

type StealerSearchResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message,omitempty"`
	Data    *StealerSearchData `json:"data,omitempty"`
}

type StealerSearchData struct {
	Results      []StealerResult `json:"results"`
	ResultsFound int             `json:"results_found"`
	ResultsShown int             `json:"results_shown"`
	Cursor       string          `json:"cursor,omitempty"`
	Meta         *ResponseMeta   `json:"_meta,omitempty"`
}

type StealerResult struct {
	LOG      string   `json:"LOG,omitempty"`
	Domain   []string `json:"domain,omitempty"`
	Email    []string `json:"email,omitempty"`
	Username []string `json:"username,omitempty"`
	Password []string `json:"password,omitempty"`
	URL      []string `json:"url,omitempty"`
	IP       string   `json:"ip,omitempty"`
	Country  string   `json:"country,omitempty"`
}

// ============================================
// V2 STEALER TYPES
// ============================================

type V2StealerResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message,omitempty"`
	Data    *V2StealerData `json:"data,omitempty"`
}

type V2SearchMeta struct {
	Total      int     `json:"total,omitempty"`
	Count      int     `json:"count,omitempty"`
	TookMs     int     `json:"took_ms,omitempty"`
	HasMore    bool    `json:"has_more,omitempty"`
	TotalPages int     `json:"total_pages,omitempty"`
	MaxScore   float64 `json:"max_score,omitempty"`
}

type V2StealerData struct {
	Items      []V2StealerResult `json:"items"`
	Meta       *V2SearchMeta     `json:"meta,omitempty"`
	NextCursor string            `json:"next_cursor,omitempty"`
	APIMeta    *ResponseMeta     `json:"_meta,omitempty"`
}

type V2StealerResult struct {
	ID           string   `json:"id,omitempty"`
	LogID        string   `json:"log_id,omitempty"`
	URL          string   `json:"url,omitempty"`
	Domain       []string `json:"domain,omitempty"`
	Subdomain    []string `json:"subdomain,omitempty"`
	EmailDomains []string `json:"email_domains,omitempty"`
	Path         []string `json:"path,omitempty"`
	Username     string   `json:"username,omitempty"`
	Password     string   `json:"password,omitempty"`
	Email        []string `json:"email,omitempty"`
	Log          string   `json:"log,omitempty"`
	PwnedAt      string   `json:"pwned_at,omitempty"`
	IndexedAt    string   `json:"indexed_at,omitempty"`
}

type SubdomainResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message,omitempty"`
	Data    *SubdomainData `json:"data,omitempty"`
}

type SubdomainData struct {
	Subdomains []string      `json:"subdomains"`
	Count      int           `json:"count"`
	Domain     string        `json:"domain"`
	Meta       *ResponseMeta `json:"_meta,omitempty"`
}

// ============================================
// V2 VICTIMS TYPES
// ============================================

type V2VictimsResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message,omitempty"`
	Data    *V2VictimsData `json:"data,omitempty"`
}

type V2VictimsData struct {
	Items      []V2VictimResult `json:"items"`
	Meta       *V2SearchMeta    `json:"meta,omitempty"`
	NextCursor string           `json:"next_cursor,omitempty"`
	APIMeta    *ResponseMeta    `json:"_meta,omitempty"`
}

type V2VictimResult struct {
	LogID        string   `json:"log_id,omitempty"`
	DeviceUsers  []string `json:"device_users,omitempty"`
	HWIDs        []string `json:"hwids,omitempty"`
	DeviceIPs    []string `json:"device_ips,omitempty"`
	DeviceEmails []string `json:"device_emails,omitempty"`
	DiscordIDs   []string `json:"discord_ids,omitempty"`
	TotalDocs    int      `json:"total_docs,omitempty"`
	PwnedAt      string   `json:"pwned_at,omitempty"`
	IndexedAt    string   `json:"indexed_at,omitempty"`
}

// VictimManifestData represents the file tree for a victim log.
// Note: This is returned unwrapped from the API (no success/data wrapper).
type VictimManifestData struct {
	LogID      string              `json:"log_id"`
	LogName    string              `json:"log_name,omitempty"`
	VictimTree *VictimManifestNode `json:"victim_tree"`
	APIMeta    *ResponseMeta       `json:"_meta,omitempty"`
}

type VictimManifestNode struct {
	ID        string               `json:"id"`
	Name      string               `json:"name"`
	Type      string               `json:"type"` // "file" or "directory"
	SizeBytes int64                `json:"size_bytes,omitempty"`
	Children  []VictimManifestNode `json:"children,omitempty"`
}

// ============================================
// V2 FILE SEARCH TYPES
// ============================================

// FileSearchJobResponse wraps a file search job (for wrapped responses).
type FileSearchJobResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message,omitempty"`
	Data    *FileSearchJobData `json:"data,omitempty"`
}

// FileSearchJobData represents file search job status.
// Note: API returns this unwrapped (no success/data wrapper).
type FileSearchJobData struct {
	JobID           string               `json:"job_id,omitempty"`
	Status          string               `json:"status,omitempty"` // queued, running, completed, canceled
	CreatedAt       string               `json:"created_at,omitempty"`
	StartedAt       string               `json:"started_at,omitempty"`
	CompletedAt     string               `json:"completed_at,omitempty"`
	ExpiresAt       string               `json:"expires_at,omitempty"`
	NextPollAfterMs int                  `json:"next_poll_after_ms,omitempty"`
	Progress        *FileSearchProgress  `json:"progress,omitempty"`
	Summary         *FileSearchSummary   `json:"summary,omitempty"`
	Limits          *FileSearchLimits    `json:"limits,omitempty"`
	Matches         []FileSearchMatch    `json:"matches,omitempty"`
	APIMeta         *ResponseMeta        `json:"_meta,omitempty"`
}

type FileSearchProgress struct {
	LogsTotal     int    `json:"logs_total,omitempty"`
	LogsCompleted int    `json:"logs_completed,omitempty"`
	FilesTotal    int    `json:"files_total,omitempty"`
	FilesScanned  int    `json:"files_scanned,omitempty"`
	Percent       int    `json:"percent,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
}

type FileSearchSummary struct {
	FilesTotal     int   `json:"files_total,omitempty"`
	FilesScanned   int   `json:"files_scanned,omitempty"`
	FilesMatched   int   `json:"files_matched,omitempty"`
	Matches        int   `json:"matches,omitempty"`
	DurationMs     int   `json:"duration_ms,omitempty"`
	BytesScanned   int64 `json:"bytes_scanned,omitempty"`
	BudgetExceeded bool  `json:"budget_exceeded,omitempty"`
	Truncated      bool  `json:"truncated,omitempty"`
	Timeouts       int   `json:"timeouts,omitempty"`
}

type FileSearchLimits struct {
	ByteBudgetBytes     int64 `json:"byte_budget_bytes,omitempty"`
	JobTTLSeconds       int   `json:"job_ttl_seconds,omitempty"`
	MaxContextLines     int   `json:"max_context_lines,omitempty"`
	MaxExpressionLength int   `json:"max_expression_length,omitempty"`
	MaxFileSizeBytes    int64 `json:"max_file_size_bytes,omitempty"`
	MaxLogIDs           int   `json:"max_log_ids,omitempty"`
	MaxMatches          int   `json:"max_matches,omitempty"`
}

type FileSearchSnippet struct {
	Line      string   `json:"line,omitempty"`
	Pre       []string `json:"pre,omitempty"`
	Post      []string `json:"post,omitempty"`
	Truncated bool     `json:"truncated,omitempty"`
}

type FileSearchColumnRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type FileSearchMatch struct {
	LogID        string                 `json:"log_id"`
	FileID       string                 `json:"file_id"`
	FileName     string                 `json:"file_name"`
	RelativePath string                 `json:"relative_path"`
	SizeBytes    int64                  `json:"size_bytes,omitempty"`
	LineNumber   int                    `json:"line_number,omitempty"`
	ColumnRange  *FileSearchColumnRange `json:"column_range,omitempty"`
	MatchText    string                 `json:"match_text,omitempty"`
	Snippet      *FileSearchSnippet     `json:"snippet,omitempty"`
}

// ============================================
// V2 EXPORT TYPES
// ============================================

type ExportJobResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message,omitempty"`
	Data    *ExportJobData `json:"data,omitempty"`
}

type ExportJobData struct {
	JobID           string          `json:"job_id,omitempty"`
	Status          string          `json:"status,omitempty"`
	Progress        *ExportProgress `json:"progress,omitempty"`
	Result          *ExportResult   `json:"result,omitempty"`
	CreatedAt       string          `json:"created_at,omitempty"`
	StartedAt       string          `json:"started_at,omitempty"`
	CompletedAt     string          `json:"completed_at,omitempty"`
	ExpiresAt       string          `json:"expires_at,omitempty"`
	NextPollAfterMs int             `json:"next_poll_after_ms,omitempty"`
	Meta            *ResponseMeta   `json:"_meta,omitempty"`
}

type ExportProgress struct {
	RecordsDone  int     `json:"records_done,omitempty"`
	RecordsTotal int     `json:"records_total,omitempty"`
	BytesDone    int64   `json:"bytes_done,omitempty"`
	Percent      float64 `json:"percent,omitempty"`
	UpdatedAt    string  `json:"updated_at,omitempty"`
}

type ExportResult struct {
	FileName    string `json:"file_name,omitempty"`
	FileSize    int64  `json:"file_size,omitempty"`
	Records     int    `json:"records,omitempty"`
	Format      string `json:"format,omitempty"`
	ExpiresAt   string `json:"expires_at,omitempty"`
	DownloadURL string `json:"download_url,omitempty"`
}

// ============================================
// OSINT TYPES
// ============================================

type IPInfoResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    *IPInfoData `json:"data,omitempty"`
}

type IPInfoData struct {
	Status        string        `json:"status,omitempty"`
	Query         string        `json:"query,omitempty"`
	Continent     string        `json:"continent,omitempty"`
	ContinentCode string        `json:"continentCode,omitempty"`
	Country       string        `json:"country,omitempty"`
	CountryCode   string        `json:"countryCode,omitempty"`
	Region        string        `json:"region,omitempty"`
	RegionName    string        `json:"regionName,omitempty"`
	City          string        `json:"city,omitempty"`
	District      string        `json:"district,omitempty"`
	Zip           string        `json:"zip,omitempty"`
	Lat           float64       `json:"lat,omitempty"`
	Lon           float64       `json:"lon,omitempty"`
	Timezone      string        `json:"timezone,omitempty"`
	Offset        int           `json:"offset,omitempty"`
	Currency      string        `json:"currency,omitempty"`
	ISP           string        `json:"isp,omitempty"`
	Org           string        `json:"org,omitempty"`
	ASName        string        `json:"asname,omitempty"`
	Mobile        bool          `json:"mobile,omitempty"`
	Proxy         bool          `json:"proxy,omitempty"`
	Hosting       bool          `json:"hosting,omitempty"`
	Reverse       string        `json:"reverse,omitempty"`
	Meta          *ResponseMeta `json:"_meta,omitempty"`
}

type SteamProfileResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message,omitempty"`
	Data    *SteamProfileData `json:"data,omitempty"`
}

type SteamProfileData struct {
	SteamID                  string        `json:"steam_id,omitempty"`
	Username                 string        `json:"username,omitempty"`
	ProfileURL               string        `json:"profile_url,omitempty"`
	Avatar                   string        `json:"avatar,omitempty"`
	PersonaState             int           `json:"persona_state,omitempty"`
	CommunityVisibilityState int           `json:"community_visibility_state,omitempty"`
	ProfileState             int           `json:"profile_state,omitempty"`
	LastLogoff               int64         `json:"last_logoff,omitempty"`
	TimeCreated              int64         `json:"time_created,omitempty"`
	RealName                 string        `json:"real_name,omitempty"`
	LocCountryCode           string        `json:"loc_country_code,omitempty"`
	LocStateCode             string        `json:"loc_state_code,omitempty"`
	LocCityID                int           `json:"loc_city_id,omitempty"`
	Meta                     *ResponseMeta `json:"_meta,omitempty"`
}

type XboxProfileResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message,omitempty"`
	Data    *XboxProfileData `json:"data,omitempty"`
}

type XboxProfileData struct {
	Username    string        `json:"username,omitempty"`
	XUID        string        `json:"xuid,omitempty"`
	Gamerscore  int           `json:"gamerscore,omitempty"`
	AccountTier string        `json:"account_tier,omitempty"`
	Tenure      string        `json:"tenure,omitempty"`
	Bio         string        `json:"bio,omitempty"`
	Location    string        `json:"location,omitempty"`
	RealName    string        `json:"real_name,omitempty"`
	Avatar      string        `json:"avatar,omitempty"`
	Meta        *ResponseMeta `json:"_meta,omitempty"`
}

type DiscordUserResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message,omitempty"`
	Data    *DiscordUserData `json:"data,omitempty"`
}

type DiscordUserData struct {
	ID            string        `json:"id,omitempty"`
	Username      string        `json:"username,omitempty"`
	GlobalName    string        `json:"global_name,omitempty"`
	Avatar        string        `json:"avatar,omitempty"`
	Discriminator string        `json:"discriminator,omitempty"`
	PublicFlags   int           `json:"public_flags,omitempty"`
	Flags         int           `json:"flags,omitempty"`
	Banner        string        `json:"banner,omitempty"`
	BannerColor   string        `json:"banner_color,omitempty"`
	AccentColor   int           `json:"accent_color,omitempty"`
	Bio           string        `json:"bio,omitempty"`
	CreatedAt     string        `json:"created_at,omitempty"`
	Meta          *ResponseMeta `json:"_meta,omitempty"`
}

type DiscordUsernameHistoryResponse struct {
	Success bool                        `json:"success"`
	Message string                      `json:"message,omitempty"`
	Data    *DiscordUsernameHistoryData `json:"data,omitempty"`
}

type DiscordUsernameHistoryData struct {
	UserID  string                       `json:"user_id"`
	History []DiscordUsernameHistoryEntry `json:"history"`
	Meta    *ResponseMeta                `json:"_meta,omitempty"`
}

type DiscordUsernameHistoryEntry struct {
	Username  string `json:"username"`
	ChangedAt string `json:"changed_at,omitempty"`
}

type DiscordToRobloxResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message,omitempty"`
	Data    *DiscordToRobloxData `json:"data,omitempty"`
}

type DiscordToRobloxData struct {
	DiscordID      string        `json:"discord_id"`
	RobloxID       string        `json:"roblox_id"`
	RobloxUsername string        `json:"roblox_username,omitempty"`
	Verified       bool          `json:"verified,omitempty"`
	Meta           *ResponseMeta `json:"_meta,omitempty"`
}

type RobloxUserResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	Data    *RobloxUserData `json:"data,omitempty"`
}

type RobloxUserData struct {
	UserID           string        `json:"user_id,omitempty"`
	Username         string        `json:"username,omitempty"`
	DisplayName      string        `json:"display_name,omitempty"`
	Description      string        `json:"description,omitempty"`
	Created          string        `json:"created,omitempty"`
	IsBanned         bool          `json:"is_banned,omitempty"`
	HasVerifiedBadge bool          `json:"has_verified_badge,omitempty"`
	Meta             *ResponseMeta `json:"_meta,omitempty"`
}

type HoleheResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    *HoleheData `json:"data,omitempty"`
}

type HoleheData struct {
	Email   string        `json:"email"`
	Domains []string      `json:"domains"`
	Meta    *ResponseMeta `json:"_meta,omitempty"`
}

type GHuntResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message,omitempty"`
	Data    *GHuntData `json:"data,omitempty"`
}

type GHuntData struct {
	Email   string         `json:"email"`
	Found   bool           `json:"found"`
	Profile *GHuntProfile  `json:"profile,omitempty"`
	Meta    *ResponseMeta  `json:"_meta,omitempty"`
}

type GHuntProfile struct {
	Name           string `json:"name,omitempty"`
	ProfilePicture string `json:"profile_picture,omitempty"`
	CoverPicture   string `json:"cover_picture,omitempty"`
	LastEdit       string `json:"last_edit,omitempty"`
	MapsID         string `json:"maps_id,omitempty"`
	CalendarID     string `json:"calendar_id,omitempty"`
}

type ExtractSubdomainResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message,omitempty"`
	Data    *ExtractSubdomainData `json:"data,omitempty"`
}

type ExtractSubdomainData struct {
	Domain     string        `json:"domain"`
	Subdomains []string      `json:"subdomains"`
	Count      int           `json:"count"`
	Meta       *ResponseMeta `json:"_meta,omitempty"`
}

type MinecraftHistoryResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message,omitempty"`
	Data    *MinecraftHistoryData `json:"data,omitempty"`
}

type MinecraftHistoryData struct {
	UUID     string                   `json:"uuid,omitempty"`
	Username string                   `json:"username"`
	History  []MinecraftHistoryEntry  `json:"history"`
	Meta     *ResponseMeta            `json:"_meta,omitempty"`
}

type MinecraftHistoryEntry struct {
	Name      string `json:"name"`
	ChangedAt string `json:"changed_at,omitempty"`
}

// ============================================
// BULK TYPES
// ============================================

type BulkJobResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message,omitempty"`
	Data    *BulkJobData `json:"data,omitempty"`
}

type BulkJobData struct {
	JobID           string        `json:"job_id,omitempty"`
	Status          string        `json:"status,omitempty"`
	TermsCount      int           `json:"terms_count,omitempty"`
	ResultsCount    int           `json:"results_count,omitempty"`
	Format          string        `json:"format,omitempty"`
	Service         string        `json:"service,omitempty"`
	CreatedAt       string        `json:"created_at,omitempty"`
	CompletedAt     string        `json:"completed_at,omitempty"`
	ExpiresAt       string        `json:"expires_at,omitempty"`
	NextPollAfterMs int           `json:"next_poll_after_ms,omitempty"`
	Meta            *ResponseMeta `json:"_meta,omitempty"`
}

type BulkJobListResponse struct {
	Results  []BulkJobData `json:"results"`
	Total    int           `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}
