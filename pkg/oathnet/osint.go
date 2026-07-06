package oathnet

import "net/url"

// OSINTService handles OSINT lookups.
type OSINTService struct {
	client *Client
}

// OSINTOptions contains optional query params shared by OSINT lookups.
type OSINTOptions struct {
	SearchID string
}

func applyOSINTOptions(params url.Values, opts []OSINTOptions) {
	for _, opt := range opts {
		if opt.SearchID != "" {
			params.Set("search_id", opt.SearchID)
		}
	}
}

// IPInfo gets IP address information.
func (s *OSINTService) IPInfo(ip string, opts ...OSINTOptions) (*IPInfoResponse, error) {
	params := url.Values{}
	params.Set("ip", ip)
	applyOSINTOptions(params, opts)

	var resp IPInfoResponse
	err := s.client.get("/service/ip-info", params, &resp)
	return &resp, err
}

// Steam gets Steam profile information.
func (s *OSINTService) Steam(steamID string, opts ...OSINTOptions) (*SteamProfileResponse, error) {
	params := url.Values{}
	params.Set("steam_id", steamID)
	applyOSINTOptions(params, opts)

	var resp SteamProfileResponse
	err := s.client.get("/service/steam", params, &resp)
	return &resp, err
}

// Xbox gets Xbox Live profile information.
func (s *OSINTService) Xbox(xblID string, opts ...OSINTOptions) (*XboxProfileResponse, error) {
	params := url.Values{}
	params.Set("xbl_id", xblID)
	applyOSINTOptions(params, opts)

	var resp XboxProfileResponse
	err := s.client.get("/service/xbox", params, &resp)
	return &resp, err
}

// DiscordUserinfo gets Discord user information.
func (s *OSINTService) DiscordUserinfo(discordID string, opts ...OSINTOptions) (*DiscordUserResponse, error) {
	params := url.Values{}
	params.Set("discord_id", discordID)
	applyOSINTOptions(params, opts)

	var resp DiscordUserResponse
	err := s.client.get("/service/discord-userinfo", params, &resp)
	return &resp, err
}

// DiscordUsernameHistory gets Discord username history.
func (s *OSINTService) DiscordUsernameHistory(discordID string, opts ...OSINTOptions) (*DiscordUsernameHistoryResponse, error) {
	params := url.Values{}
	params.Set("discord_id", discordID)
	applyOSINTOptions(params, opts)

	var resp DiscordUsernameHistoryResponse
	err := s.client.get("/service/discord-username-history", params, &resp)
	return &resp, err
}

// DiscordToRoblox gets Roblox account linked to Discord.
func (s *OSINTService) DiscordToRoblox(discordID string, opts ...OSINTOptions) (*DiscordToRobloxResponse, error) {
	params := url.Values{}
	params.Set("discord_id", discordID)
	applyOSINTOptions(params, opts)

	var resp DiscordToRobloxResponse
	err := s.client.get("/service/discord-to-roblox", params, &resp)
	return &resp, err
}

// RobloxUserinfoOptions contains options for Roblox user lookup.
type RobloxUserinfoOptions struct {
	UserID   string
	Username string
	SearchID string
}

// RobloxUserinfo gets Roblox user information.
func (s *OSINTService) RobloxUserinfo(opts RobloxUserinfoOptions) (*RobloxUserResponse, error) {
	params := url.Values{}
	if opts.UserID != "" {
		params.Set("user_id", opts.UserID)
	}
	if opts.Username != "" {
		params.Set("username", opts.Username)
	}
	if opts.SearchID != "" {
		params.Set("search_id", opts.SearchID)
	}

	var resp RobloxUserResponse
	err := s.client.get("/service/roblox-userinfo", params, &resp)
	return &resp, err
}

// Holehe checks email account existence across services.
func (s *OSINTService) Holehe(email string, opts ...OSINTOptions) (*HoleheResponse, error) {
	params := url.Values{}
	params.Set("email", email)
	applyOSINTOptions(params, opts)

	var resp HoleheResponse
	err := s.client.get("/service/holehe", params, &resp)
	return &resp, err
}

// GHunt gets Google account information.
func (s *OSINTService) GHunt(email string, opts ...OSINTOptions) (*GHuntResponse, error) {
	params := url.Values{}
	params.Set("email", email)
	applyOSINTOptions(params, opts)

	var resp GHuntResponse
	err := s.client.get("/service/ghunt", params, &resp)
	return &resp, err
}

// ExtractSubdomain extracts subdomains for a domain.
func (s *OSINTService) ExtractSubdomain(domain string, isAlive *bool, opts ...OSINTOptions) (*ExtractSubdomainResponse, error) {
	params := url.Values{}
	params.Set("domain", domain)
	if isAlive != nil {
		if *isAlive {
			params.Set("is_alive", "true")
		} else {
			params.Set("is_alive", "false")
		}
	}
	applyOSINTOptions(params, opts)

	var resp ExtractSubdomainResponse
	err := s.client.get("/service/extract-subdomain", params, &resp)
	return &resp, err
}

// MinecraftHistory gets Minecraft username history.
func (s *OSINTService) MinecraftHistory(username string, opts ...OSINTOptions) (*MinecraftHistoryResponse, error) {
	params := url.Values{}
	params.Set("username", username)
	applyOSINTOptions(params, opts)

	var resp MinecraftHistoryResponse
	err := s.client.get("/service/mc-history", params, &resp)
	return &resp, err
}
