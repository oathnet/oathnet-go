package oathnet

import "net/url"

// OSINTService handles OSINT lookups.
type OSINTService struct {
	client *Client
}

// IPInfo gets IP address information.
func (s *OSINTService) IPInfo(ip string) (*IPInfoResponse, error) {
	params := url.Values{}
	params.Set("ip", ip)

	var resp IPInfoResponse
	err := s.client.get("/service/ip-info", params, &resp)
	return &resp, err
}

// Steam gets Steam profile information.
func (s *OSINTService) Steam(steamID string) (*SteamProfileResponse, error) {
	params := url.Values{}
	params.Set("steam_id", steamID)

	var resp SteamProfileResponse
	err := s.client.get("/service/steam", params, &resp)
	return &resp, err
}

// Xbox gets Xbox Live profile information.
func (s *OSINTService) Xbox(xblID string) (*XboxProfileResponse, error) {
	params := url.Values{}
	params.Set("xbl_id", xblID)

	var resp XboxProfileResponse
	err := s.client.get("/service/xbox", params, &resp)
	return &resp, err
}

// DiscordUserinfo gets Discord user information.
func (s *OSINTService) DiscordUserinfo(discordID string) (*DiscordUserResponse, error) {
	params := url.Values{}
	params.Set("discord_id", discordID)

	var resp DiscordUserResponse
	err := s.client.get("/service/discord-userinfo", params, &resp)
	return &resp, err
}

// DiscordUsernameHistory gets Discord username history.
func (s *OSINTService) DiscordUsernameHistory(discordID string) (*DiscordUsernameHistoryResponse, error) {
	params := url.Values{}
	params.Set("discord_id", discordID)

	var resp DiscordUsernameHistoryResponse
	err := s.client.get("/service/discord-username-history", params, &resp)
	return &resp, err
}

// DiscordToRoblox gets Roblox account linked to Discord.
func (s *OSINTService) DiscordToRoblox(discordID string) (*DiscordToRobloxResponse, error) {
	params := url.Values{}
	params.Set("discord_id", discordID)

	var resp DiscordToRobloxResponse
	err := s.client.get("/service/discord-to-roblox", params, &resp)
	return &resp, err
}

// RobloxUserinfoOptions contains options for Roblox user lookup.
type RobloxUserinfoOptions struct {
	UserID   string
	Username string
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

	var resp RobloxUserResponse
	err := s.client.get("/service/roblox-userinfo", params, &resp)
	return &resp, err
}

// Holehe checks email account existence across services.
func (s *OSINTService) Holehe(email string) (*HoleheResponse, error) {
	params := url.Values{}
	params.Set("email", email)

	var resp HoleheResponse
	err := s.client.get("/service/holehe", params, &resp)
	return &resp, err
}

// GHunt gets Google account information.
func (s *OSINTService) GHunt(email string) (*GHuntResponse, error) {
	params := url.Values{}
	params.Set("email", email)

	var resp GHuntResponse
	err := s.client.get("/service/ghunt", params, &resp)
	return &resp, err
}

// ExtractSubdomain extracts subdomains for a domain.
func (s *OSINTService) ExtractSubdomain(domain string, isAlive *bool) (*ExtractSubdomainResponse, error) {
	params := url.Values{}
	params.Set("domain", domain)
	if isAlive != nil {
		if *isAlive {
			params.Set("is_alive", "true")
		} else {
			params.Set("is_alive", "false")
		}
	}

	var resp ExtractSubdomainResponse
	err := s.client.get("/service/extract-subdomain", params, &resp)
	return &resp, err
}

// MinecraftHistory gets Minecraft username history.
func (s *OSINTService) MinecraftHistory(username string) (*MinecraftHistoryResponse, error) {
	params := url.Values{}
	params.Set("username", username)

	var resp MinecraftHistoryResponse
	err := s.client.get("/service/minecraft-history", params, &resp)
	return &resp, err
}
