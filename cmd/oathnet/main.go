package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/oathnet/oathnet-go/pkg/oathnet"
	"github.com/spf13/cobra"
)

var (
	apiKey       string
	outputFormat string
	client       *oathnet.Client
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "oathnet",
	Short: "OathNet CLI - Search breach databases, stealer logs, and OSINT lookups",
	Long:  `OathNet CLI - Command-line interface for the OathNet API.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if apiKey == "" {
			apiKey = os.Getenv("OATHNET_API_KEY")
		}
		if apiKey == "" && cmd.Name() != "help" && cmd.Name() != "version" {
			return fmt.Errorf("API key is required. Use --api-key or set OATHNET_API_KEY")
		}
		if apiKey != "" {
			var err error
			client, err = oathnet.NewClient(apiKey)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&apiKey, "api-key", "k", "", "OathNet API key")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", "table", "Output format (table|json)")

	// Add command groups
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(stealerCmd)
	rootCmd.AddCommand(victimsCmd)
	rootCmd.AddCommand(fileSearchCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(osintCmd)
	rootCmd.AddCommand(bulkCmd)
	rootCmd.AddCommand(utilCmd)
}

func output(data interface{}, tableFormatter func()) {
	if outputFormat == "json" {
		jsonData, _ := json.MarshalIndent(data, "", "  ")
		fmt.Println(string(jsonData))
	} else if tableFormatter != nil {
		tableFormatter()
	}
}

// formatValue formats any value for display
func formatValue(value interface{}) string {
	if value == nil {
		return "N/A"
	}
	switch v := value.(type) {
	case string:
		if v == "" {
			return "N/A"
		}
		return v
	case []string:
		if len(v) == 0 {
			return "N/A"
		}
		if len(v) <= 5 {
			return strings.Join(v, ", ")
		}
		return strings.Join(v[:5], ", ") + fmt.Sprintf(" (+%d more)", len(v)-5)
	case []interface{}:
		if len(v) == 0 {
			return "N/A"
		}
		strs := make([]string, 0, len(v))
		for _, item := range v {
			strs = append(strs, fmt.Sprintf("%v", item))
		}
		if len(strs) <= 5 {
			return strings.Join(strs, ", ")
		}
		return strings.Join(strs[:5], ", ") + fmt.Sprintf(" (+%d more)", len(strs)-5)
	case map[string]interface{}:
		jsonBytes, _ := json.Marshal(v)
		return string(jsonBytes)
	case float64:
		if v == float64(int(v)) {
			return fmt.Sprintf("%d", int(v))
		}
		return fmt.Sprintf("%.2f", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// structToMap converts a struct to map[string]interface{} via JSON
func structToMap(v interface{}) map[string]interface{} {
	jsonBytes, _ := json.Marshal(v)
	var m map[string]interface{}
	json.Unmarshal(jsonBytes, &m)
	return m
}

// printDynamicFields prints all fields from a map dynamically
func printDynamicFields(data map[string]interface{}, priorityFields []string) {
	shown := make(map[string]bool)

	// Show priority fields first
	for _, field := range priorityFields {
		if val, ok := data[field]; ok && !isEmpty(val) {
			if field == "log_id" || field == "id" {
				fmt.Printf("  \033[1;33m%s:\033[0m \033[32m%s\033[0m\n", field, formatValue(val))
			} else {
				fmt.Printf("  \033[1m%s:\033[0m %s\n", field, formatValue(val))
			}
			shown[field] = true
		}
	}

	// Collect and sort remaining fields
	var remaining []string
	for k := range data {
		if !shown[k] && !strings.HasPrefix(k, "_") && !isEmpty(data[k]) {
			remaining = append(remaining, k)
		}
	}
	sort.Strings(remaining)

	// Show remaining fields
	for _, k := range remaining {
		fmt.Printf("  %s: %s\n", k, formatValue(data[k]))
	}
}

// isEmpty checks if a value should be considered empty
func isEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String:
		return rv.String() == ""
	case reflect.Slice, reflect.Array:
		return rv.Len() == 0
	case reflect.Map:
		return rv.Len() == 0
	}
	return false
}

// ============================================
// SEARCH COMMANDS
// ============================================

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search breach and stealer databases",
}

var searchBreachCmd = &cobra.Command{
	Use:   "breach",
	Short: "Search breach database",
	RunE: func(cmd *cobra.Command, args []string) error {
		query, _ := cmd.Flags().GetString("query")
		cursor, _ := cmd.Flags().GetString("cursor")
		dbnames, _ := cmd.Flags().GetString("dbnames")
		outputFile, _ := cmd.Flags().GetString("output")

		result, err := client.Search.Breach(query, &oathnet.SearchOptions{
			Cursor:  cursor,
			DBNames: dbnames,
		})
		if err != nil {
			return err
		}

		if outputFile != "" {
			jsonData, _ := json.MarshalIndent(result, "", "  ")
			os.WriteFile(outputFile, jsonData, 0644)
			fmt.Printf("Saved to %s\n", outputFile)
		}

		output(result, func() {
			if result.Data != nil {
				fmt.Printf("\nFound %d results (showing %d)\n\n", result.Data.ResultsFound, result.Data.ResultsShown)

				// Priority fields for breach results
				priorityFields := []string{"id", "dbname", "email", "username", "password", "password_hash",
					"phone_number", "ip", "domain", "country", "city", "full_name",
					"first_name", "last_name", "date"}

				for i, r := range result.Data.Results {
					fmt.Printf("\033[36m━━━ Result %d ━━━\033[0m\n", i+1)
					data := structToMap(r)
					printDynamicFields(data, priorityFields)
					fmt.Println()
				}

				if result.Data.Cursor != "" {
					fmt.Printf("\033[33mNext cursor:\033[0m %s\n", result.Data.Cursor)
					fmt.Println("Use --cursor to fetch next page")
				}
			} else {
				fmt.Println("No results found")
			}
		})
		return nil
	},
}

var searchStealerCmd = &cobra.Command{
	Use:   "stealer",
	Short: "Search stealer database (legacy)",
	RunE: func(cmd *cobra.Command, args []string) error {
		query, _ := cmd.Flags().GetString("query")
		cursor, _ := cmd.Flags().GetString("cursor")

		result, err := client.Search.Stealer(query, &oathnet.SearchOptions{Cursor: cursor})
		if err != nil {
			return err
		}

		output(result, func() {
			if result.Data != nil {
				fmt.Printf("\nFound %d results (showing %d)\n\n", result.Data.ResultsFound, result.Data.ResultsShown)
				for i, r := range result.Data.Results {
					fmt.Printf("Result %d:\n", i+1)
					fmt.Printf("  LOG: %s\n", r.LOG)
					if len(r.Domain) > 0 {
						fmt.Printf("  Domain: %s\n", strings.Join(r.Domain, ", "))
					}
					if len(r.Email) > 0 {
						fmt.Printf("  Email: %s\n", strings.Join(r.Email, ", "))
					}
					fmt.Println()
				}
			}
		})
		return nil
	},
}

var searchInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a search session",
	RunE: func(cmd *cobra.Command, args []string) error {
		query, _ := cmd.Flags().GetString("query")

		result, err := client.Search.InitSession(query)
		if err != nil {
			return err
		}

		output(result, func() {
			if result.Data != nil {
				fmt.Println("\nSession initialized")
				fmt.Printf("Session ID: %s\n", result.Data.Session.ID)
				fmt.Printf("Query: %s\n", result.Data.Session.Query)
				fmt.Printf("Search Type: %s\n", result.Data.Session.SearchType)
				fmt.Printf("Expires: %s\n", result.Data.Session.ExpiresAt)
			}
		})
		return nil
	},
}

func init() {
	searchBreachCmd.Flags().StringP("query", "q", "", "Search query (required)")
	searchBreachCmd.Flags().String("cursor", "", "Pagination cursor")
	searchBreachCmd.Flags().String("dbnames", "", "Filter by database names")
	searchBreachCmd.Flags().StringP("output", "o", "", "Save results to file")
	searchBreachCmd.MarkFlagRequired("query")

	searchStealerCmd.Flags().StringP("query", "q", "", "Search query (required)")
	searchStealerCmd.Flags().String("cursor", "", "Pagination cursor")
	searchStealerCmd.MarkFlagRequired("query")

	searchInitCmd.Flags().StringP("query", "q", "", "Search query (required)")
	searchInitCmd.MarkFlagRequired("query")

	searchCmd.AddCommand(searchBreachCmd)
	searchCmd.AddCommand(searchStealerCmd)
	searchCmd.AddCommand(searchInitCmd)
}

// ============================================
// OSINT COMMANDS
// ============================================

var osintCmd = &cobra.Command{
	Use:   "osint",
	Short: "OSINT lookups",
}

var osintIPCmd = &cobra.Command{
	Use:   "ip [ip_address]",
	Short: "Get IP address information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := client.OSINT.IPInfo(args[0])
		if err != nil {
			return err
		}

		output(result, func() {
			if result.Data != nil {
				fmt.Println("\nIP Information")
				fmt.Printf("IP: %s\n", result.Data.Query)
				fmt.Printf("Country: %s (%s)\n", result.Data.Country, result.Data.CountryCode)
				fmt.Printf("Region: %s\n", result.Data.RegionName)
				fmt.Printf("City: %s\n", result.Data.City)
				fmt.Printf("ISP: %s\n", result.Data.ISP)
				fmt.Printf("Org: %s\n", result.Data.Org)
				fmt.Printf("Mobile: %v\n", result.Data.Mobile)
				fmt.Printf("Proxy: %v\n", result.Data.Proxy)
			}
		})
		return nil
	},
}

var osintSteamCmd = &cobra.Command{
	Use:   "steam [steam_id]",
	Short: "Get Steam profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := client.OSINT.Steam(args[0])
		if err != nil {
			return err
		}

		output(result, func() {
			if result.Data != nil {
				fmt.Println("\nSteam Profile")
				fmt.Printf("Username: %s\n", result.Data.Username)
				fmt.Printf("Steam ID: %s\n", result.Data.SteamID)
				fmt.Printf("Profile URL: %s\n", result.Data.ProfileURL)
				if result.Data.RealName != "" {
					fmt.Printf("Real Name: %s\n", result.Data.RealName)
				}
			}
		})
		return nil
	},
}

var osintXboxCmd = &cobra.Command{
	Use:   "xbox [gamertag]",
	Short: "Get Xbox Live profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := client.OSINT.Xbox(args[0])
		if err != nil {
			return err
		}

		output(result, func() {
			if result.Data != nil {
				fmt.Println("\nXbox Profile")
				fmt.Printf("Gamertag: %s\n", result.Data.Username)
				fmt.Printf("Gamerscore: %d\n", result.Data.Gamerscore)
				fmt.Printf("Account Tier: %s\n", result.Data.AccountTier)
				if result.Data.Bio != "" {
					fmt.Printf("Bio: %s\n", result.Data.Bio)
				}
			}
		})
		return nil
	},
}

var osintDiscordCmd = &cobra.Command{
	Use:   "discord",
	Short: "Discord lookups",
}

var osintDiscordUserCmd = &cobra.Command{
	Use:   "user [discord_id]",
	Short: "Get Discord user info",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := client.OSINT.DiscordUserinfo(args[0])
		if err != nil {
			return err
		}

		output(result, func() {
			if result.Data != nil {
				fmt.Println("\nDiscord User")
				fmt.Printf("ID: %s\n", result.Data.ID)
				fmt.Printf("Username: %s\n", result.Data.Username)
				if result.Data.GlobalName != "" {
					fmt.Printf("Display Name: %s\n", result.Data.GlobalName)
				}
				if result.Data.Bio != "" {
					fmt.Printf("Bio: %s\n", result.Data.Bio)
				}
				fmt.Printf("Created: %s\n", result.Data.CreatedAt)
			}
		})
		return nil
	},
}

var osintDiscordHistoryCmd = &cobra.Command{
	Use:   "history [discord_id]",
	Short: "Get Discord username history",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := client.OSINT.DiscordUsernameHistory(args[0])
		if err != nil {
			return err
		}

		output(result, func() {
			if result.Data != nil {
				fmt.Println("\nDiscord Username History")
				fmt.Printf("User ID: %s\n", result.Data.UserID)
				if len(result.Data.History) > 0 {
					for _, h := range result.Data.History {
						fmt.Printf("  %s (%s)\n", h.Username, h.ChangedAt)
					}
				} else {
					fmt.Println("No history found")
				}
			}
		})
		return nil
	},
}

var osintDiscordRobloxCmd = &cobra.Command{
	Use:   "roblox [discord_id]",
	Short: "Get Roblox account linked to Discord",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := client.OSINT.DiscordToRoblox(args[0])
		if err != nil {
			return err
		}

		output(result, func() {
			if result.Data != nil {
				fmt.Println("\nDiscord to Roblox")
				fmt.Printf("Discord ID: %s\n", result.Data.DiscordID)
				fmt.Printf("Roblox ID: %s\n", result.Data.RobloxID)
				if result.Data.RobloxUsername != "" {
					fmt.Printf("Roblox Username: %s\n", result.Data.RobloxUsername)
				}
			}
		})
		return nil
	},
}

var osintHoleheCmd = &cobra.Command{
	Use:   "holehe [email]",
	Short: "Check email account existence",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := client.OSINT.Holehe(args[0])
		if err != nil {
			return err
		}

		output(result, func() {
			if result.Data != nil {
				fmt.Println("\nHolehe Results")
				fmt.Printf("Email: %s\n", result.Data.Email)
				found := 0
				for _, d := range result.Data.Domains {
					if d.Exists {
						found++
					}
				}
				fmt.Printf("Found on %d services:\n", found)
				for _, d := range result.Data.Domains {
					if d.Exists {
						fmt.Printf("  + %s\n", d.Domain)
					}
				}
			}
		})
		return nil
	},
}

var osintSubdomainCmd = &cobra.Command{
	Use:   "subdomain [domain]",
	Short: "Extract subdomains for a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		alive, _ := cmd.Flags().GetBool("alive")
		var alivePtr *bool
		if cmd.Flags().Changed("alive") {
			alivePtr = &alive
		}

		result, err := client.OSINT.ExtractSubdomain(args[0], alivePtr)
		if err != nil {
			return err
		}

		output(result, func() {
			if result.Data != nil {
				fmt.Printf("\nSubdomains for %s\n", result.Data.Domain)
				fmt.Printf("Found: %d\n\n", result.Data.Count)
				for _, s := range result.Data.Subdomains {
					fmt.Printf("  %s\n", s)
				}
			}
		})
		return nil
	},
}

func init() {
	osintSubdomainCmd.Flags().Bool("alive", false, "Only return alive subdomains")

	osintDiscordCmd.AddCommand(osintDiscordUserCmd)
	osintDiscordCmd.AddCommand(osintDiscordHistoryCmd)
	osintDiscordCmd.AddCommand(osintDiscordRobloxCmd)

	osintCmd.AddCommand(osintIPCmd)
	osintCmd.AddCommand(osintSteamCmd)
	osintCmd.AddCommand(osintXboxCmd)
	osintCmd.AddCommand(osintDiscordCmd)
	osintCmd.AddCommand(osintHoleheCmd)
	osintCmd.AddCommand(osintSubdomainCmd)
}

// ============================================
// STEALER V2 COMMANDS
// ============================================

var stealerCmd = &cobra.Command{
	Use:   "stealer",
	Short: "V2 Stealer search commands",
}

var stealerSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search V2 stealer database",
	RunE: func(cmd *cobra.Command, args []string) error {
		query, _ := cmd.Flags().GetString("query")
		domain, _ := cmd.Flags().GetString("domain")
		wildcard, _ := cmd.Flags().GetBool("wildcard")
		hasLogID, _ := cmd.Flags().GetBool("has-log-id")
		pageSize, _ := cmd.Flags().GetInt("page-size")
		fileSearch, _ := cmd.Flags().GetString("file-search")
		fileSearchMode, _ := cmd.Flags().GetString("file-search-mode")
		outputFile, _ := cmd.Flags().GetString("output")

		result, err := client.Stealer.Search(query, &oathnet.StealerSearchOptions{
			Domain:   domain,
			Wildcard: wildcard,
			HasLogID: hasLogID,
			PageSize: pageSize,
		})
		if err != nil {
			return err
		}

		if outputFile != "" {
			jsonData, _ := json.MarshalIndent(result, "", "  ")
			os.WriteFile(outputFile, jsonData, 0644)
			fmt.Printf("Saved to %s\n", outputFile)
		}

		output(result, func() {
			if result.Data != nil && len(result.Data.Items) > 0 {
				total := len(result.Data.Items)
				if result.Data.Meta != nil {
					total = result.Data.Meta.Total
				}
				fmt.Printf("\nFound %d results (showing %d)\n\n", total, len(result.Data.Items))

				// Priority fields for stealer results
				priorityFields := []string{"id", "log_id", "url", "username", "password", "email", "domain",
					"subdomain", "path", "log", "pwned_at", "indexed_at"}

				// Collect log_ids
				var logIDs []string

				for i, r := range result.Data.Items {
					fmt.Printf("\033[36m━━━ Result %d ━━━\033[0m\n", i+1)
					data := structToMap(r)
					printDynamicFields(data, priorityFields)
					fmt.Println()

					// Collect log_id
					if r.LogID != "" {
						logIDs = append(logIDs, r.LogID)
					}
				}

				// Show all collected log_ids
				if len(logIDs) > 0 {
					fmt.Println("\033[1;33m\n═══ LOG IDs (use with file-search/victims) ═══\033[0m")
					for _, lid := range logIDs {
						fmt.Printf("  \033[32m%s\033[0m\n", lid)
					}
					fmt.Println()
				}

				if result.Data.NextCursor != "" {
					fmt.Printf("\033[33mNext cursor:\033[0m %s\n", result.Data.NextCursor)
				}

				// Auto file-search if requested
				if fileSearch != "" && len(logIDs) > 0 {
					fmt.Printf("\033[1;35m\n═══ Auto File Search: '%s' ═══\033[0m\n", fileSearch)
					searchLogIDs := logIDs
					if len(searchLogIDs) > 10 {
						searchLogIDs = searchLogIDs[:10]
					}
					fsResult, err := client.FileSearch.Search(fileSearch, &oathnet.FileSearchCreateOptions{
						SearchMode: fileSearchMode,
						LogIDs:     searchLogIDs,
						MaxMatches: 50,
					}, 60*time.Second)
					if err != nil {
						fmt.Printf("\033[31mFile search error: %s\033[0m\n", err)
					} else if fsResult.Data != nil && len(fsResult.Data.Matches) > 0 {
						fmt.Printf("\033[32mFound %d matches!\033[0m\n\n", len(fsResult.Data.Matches))
						for i, m := range fsResult.Data.Matches {
							if i >= 20 {
								break
							}
							fmt.Printf("  \033[36m%s\033[0m (log: %s)\n", m.FileName, m.LogID)
							if m.MatchText != "" {
								text := m.MatchText
								if len(text) > 100 {
									text = text[:100] + "..."
								}
								fmt.Printf("    → %s\n", text)
							}
						}
					} else {
						fmt.Println("\033[33mNo file matches found\033[0m")
					}
				}
			} else {
				fmt.Println("No results found")
			}
		})
		return nil
	},
}

var stealerSubdomainCmd = &cobra.Command{
	Use:   "subdomain",
	Short: "Extract subdomains from stealer data",
	RunE: func(cmd *cobra.Command, args []string) error {
		domain, _ := cmd.Flags().GetString("domain")

		result, err := client.Stealer.Subdomain(domain, "")
		if err != nil {
			return err
		}

		output(result, func() {
			if result.Data != nil {
				fmt.Printf("\nFound %d subdomains for %s\n\n", result.Data.Count, result.Data.Domain)
				for _, s := range result.Data.Subdomains {
					fmt.Printf("  %s\n", s)
				}
			}
		})
		return nil
	},
}

func init() {
	stealerSearchCmd.Flags().StringP("query", "q", "", "Search query")
	stealerSearchCmd.Flags().String("domain", "", "Filter by domain")
	stealerSearchCmd.Flags().Bool("wildcard", false, "Enable wildcard search")
	stealerSearchCmd.Flags().Bool("has-log-id", false, "Only results with log ID")
	stealerSearchCmd.Flags().Int("page-size", 25, "Results per page")
	stealerSearchCmd.Flags().String("file-search", "", "Auto file-search pattern in results")
	stealerSearchCmd.Flags().String("file-search-mode", "literal", "Search mode (literal|regex|wildcard)")
	stealerSearchCmd.Flags().StringP("output", "o", "", "Save results to file")

	stealerSubdomainCmd.Flags().StringP("domain", "d", "", "Domain to search (required)")
	stealerSubdomainCmd.MarkFlagRequired("domain")

	stealerCmd.AddCommand(stealerSearchCmd)
	stealerCmd.AddCommand(stealerSubdomainCmd)
}

// ============================================
// VICTIMS V2 COMMANDS
// ============================================

var victimsCmd = &cobra.Command{
	Use:   "victims",
	Short: "V2 Victims commands",
}

var victimsSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search victim profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		query, _ := cmd.Flags().GetString("query")
		email, _ := cmd.Flags().GetString("email")
		ip, _ := cmd.Flags().GetString("ip")
		discordID, _ := cmd.Flags().GetString("discord-id")
		wildcard, _ := cmd.Flags().GetBool("wildcard")
		pageSize, _ := cmd.Flags().GetInt("page-size")
		fileSearch, _ := cmd.Flags().GetString("file-search")
		fileSearchMode, _ := cmd.Flags().GetString("file-search-mode")
		outputFile, _ := cmd.Flags().GetString("output")

		result, err := client.Victims.Search(query, &oathnet.VictimsSearchOptions{
			Email:     email,
			IP:        ip,
			DiscordID: discordID,
			Wildcard:  wildcard,
			PageSize:  pageSize,
		})
		if err != nil {
			return err
		}

		if outputFile != "" {
			jsonData, _ := json.MarshalIndent(result, "", "  ")
			os.WriteFile(outputFile, jsonData, 0644)
			fmt.Printf("Saved to %s\n", outputFile)
		}

		output(result, func() {
			if result.Data != nil && len(result.Data.Items) > 0 {
				total := len(result.Data.Items)
				if result.Data.Meta != nil {
					total = result.Data.Meta.Total
				}
				fmt.Printf("\nFound %d victims (showing %d)\n\n", total, len(result.Data.Items))

				// Priority fields for victim results
				priorityFields := []string{"log_id", "device_users", "device_emails", "device_ips",
					"discord_ids", "hwids", "total_docs", "pwned_at", "indexed_at"}

				// Collect log_ids
				var logIDs []string

				for i, v := range result.Data.Items {
					fmt.Printf("\033[36m━━━ Victim %d ━━━\033[0m\n", i+1)
					data := structToMap(v)
					printDynamicFields(data, priorityFields)
					fmt.Println()

					// Collect log_id
					if v.LogID != "" {
						logIDs = append(logIDs, v.LogID)
					}
				}

				// Show all collected log_ids
				if len(logIDs) > 0 {
					fmt.Println("\033[1;33m\n═══ LOG IDs (use with file-search/manifest/archive) ═══\033[0m")
					for _, lid := range logIDs {
						fmt.Printf("  \033[32m%s\033[0m\n", lid)
					}
					fmt.Println()
					fmt.Println("Usage: oathnet victims manifest <log_id>")
					fmt.Println("       oathnet file-search create -e \"password\" --log-id <log_id>")
				}

				if result.Data.NextCursor != "" {
					fmt.Printf("\n\033[33mNext cursor:\033[0m %s\n", result.Data.NextCursor)
				}

				// Auto file-search if requested
				if fileSearch != "" && len(logIDs) > 0 {
					fmt.Printf("\033[1;35m\n═══ Auto File Search: '%s' ═══\033[0m\n", fileSearch)
					searchLogIDs := logIDs
					if len(searchLogIDs) > 10 {
						searchLogIDs = searchLogIDs[:10]
					}
					fsResult, err := client.FileSearch.Search(fileSearch, &oathnet.FileSearchCreateOptions{
						SearchMode: fileSearchMode,
						LogIDs:     searchLogIDs,
						MaxMatches: 50,
					}, 60*time.Second)
					if err != nil {
						fmt.Printf("\033[31mFile search error: %s\033[0m\n", err)
					} else if fsResult.Data != nil && len(fsResult.Data.Matches) > 0 {
						fmt.Printf("\033[32mFound %d matches!\033[0m\n\n", len(fsResult.Data.Matches))
						for i, m := range fsResult.Data.Matches {
							if i >= 20 {
								break
							}
							fmt.Printf("  \033[36m%s\033[0m (log: %s)\n", m.FileName, m.LogID)
							if m.MatchText != "" {
								text := m.MatchText
								if len(text) > 100 {
									text = text[:100] + "..."
								}
								fmt.Printf("    → %s\n", text)
							}
						}
					} else {
						fmt.Println("\033[33mNo file matches found\033[0m")
					}
				}
			} else {
				fmt.Println("No victims found")
			}
		})
		return nil
	},
}

var victimsManifestCmd = &cobra.Command{
	Use:   "manifest [log_id]",
	Short: "Get victim file manifest (file tree)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := client.Victims.GetManifest(args[0])
		if err != nil {
			return err
		}

		output(result, func() {
			if result != nil {
				fmt.Printf("\nManifest for %s\n", result.LogID)
				if result.LogName != "" {
					fmt.Printf("Log Name: %s\n", result.LogName)
				}
				if result.VictimTree != nil {
					fmt.Printf("\nFile Tree:\n")
					printVictimTree(result.VictimTree, 0)
				}
			}
		})
		return nil
	},
}

func printVictimTree(node *oathnet.VictimManifestNode, indent int) {
	if node == nil {
		return
	}
	prefix := ""
	for i := 0; i < indent; i++ {
		prefix += "  "
	}
	if node.Type == "directory" {
		fmt.Printf("%s[DIR] %s\n", prefix, node.Name)
		for i := range node.Children {
			printVictimTree(&node.Children[i], indent+1)
		}
	} else {
		fmt.Printf("%s%s (%d bytes)\n", prefix, node.Name, node.SizeBytes)
	}
}

var victimsArchiveCmd = &cobra.Command{
	Use:   "archive [log_id]",
	Short: "Download victim archive",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")
		fmt.Printf("Downloading archive for %s...\n", args[0])
		err := client.Victims.DownloadArchive(args[0], output)
		if err != nil {
			return err
		}
		fmt.Printf("Downloaded to %s\n", output)
		return nil
	},
}

func init() {
	victimsSearchCmd.Flags().StringP("query", "q", "", "Search query")
	victimsSearchCmd.Flags().String("email", "", "Filter by email")
	victimsSearchCmd.Flags().String("ip", "", "Filter by IP")
	victimsSearchCmd.Flags().String("discord-id", "", "Filter by Discord ID")
	victimsSearchCmd.Flags().Bool("wildcard", false, "Enable wildcard matching")
	victimsSearchCmd.Flags().Int("page-size", 25, "Results per page")
	victimsSearchCmd.Flags().String("file-search", "", "Auto file-search pattern in results")
	victimsSearchCmd.Flags().String("file-search-mode", "literal", "Search mode (literal|regex|wildcard)")
	victimsSearchCmd.Flags().StringP("output", "o", "", "Save results to file")

	victimsArchiveCmd.Flags().StringP("output", "o", "", "Output file path (required)")
	victimsArchiveCmd.MarkFlagRequired("output")

	victimsCmd.AddCommand(victimsSearchCmd)
	victimsCmd.AddCommand(victimsManifestCmd)
	victimsCmd.AddCommand(victimsArchiveCmd)
}

// ============================================
// FILE-SEARCH V2 COMMANDS
// ============================================

var fileSearchCmd = &cobra.Command{
	Use:   "file-search",
	Short: "V2 File search commands",
}

var fileSearchCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a file search job",
	RunE: func(cmd *cobra.Command, args []string) error {
		expression, _ := cmd.Flags().GetString("expression")
		mode, _ := cmd.Flags().GetString("mode")
		maxMatches, _ := cmd.Flags().GetInt("max-matches")

		result, err := client.FileSearch.Create(expression, &oathnet.FileSearchCreateOptions{
			SearchMode: mode,
			MaxMatches: maxMatches,
		})
		if err != nil {
			return err
		}

		output(result, func() {
			if result.Data != nil {
				fmt.Println("\nFile search job created")
				fmt.Printf("Job ID: %s\n", result.Data.JobID)
				fmt.Printf("Status: %s\n", result.Data.Status)
				fmt.Printf("\nUse: oathnet file-search status %s\n", result.Data.JobID)
			}
		})
		return nil
	},
}

var fileSearchStatusCmd = &cobra.Command{
	Use:   "status [job_id]",
	Short: "Get file search job status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := client.FileSearch.GetStatus(args[0])
		if err != nil {
			return err
		}

		output(result, func() {
			if result.Data != nil {
				fmt.Println("\nFile Search Job Status")
				fmt.Printf("Job ID: %s\n", result.Data.JobID)
				fmt.Printf("Status: %s\n", result.Data.Status)
				if result.Data.Summary != nil {
					fmt.Printf("\nSummary:\n")
					fmt.Printf("  Files scanned: %d/%d\n", result.Data.Summary.FilesScanned, result.Data.Summary.FilesTotal)
					fmt.Printf("  Files matched: %d\n", result.Data.Summary.FilesMatched)
					fmt.Printf("  Total matches: %d\n", result.Data.Summary.Matches)
				}
				if len(result.Data.Matches) > 0 {
					fmt.Printf("\nMatches (%d):\n", len(result.Data.Matches))
					for i, m := range result.Data.Matches[:min(10, len(result.Data.Matches))] {
						fmt.Printf("  %d. %s (%s)\n", i+1, m.FileName, m.LogID)
					}
				}
			}
		})
		return nil
	},
}

func init() {
	fileSearchCreateCmd.Flags().StringP("expression", "e", "", "Search expression (required)")
	fileSearchCreateCmd.Flags().String("mode", "literal", "Search mode (literal|regex|wildcard)")
	fileSearchCreateCmd.Flags().Int("max-matches", 100, "Maximum matches")
	fileSearchCreateCmd.MarkFlagRequired("expression")

	fileSearchCmd.AddCommand(fileSearchCreateCmd)
	fileSearchCmd.AddCommand(fileSearchStatusCmd)
}

// ============================================
// EXPORT V2 COMMANDS
// ============================================

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "V2 Export commands",
}

var exportCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an export job",
	RunE: func(cmd *cobra.Command, args []string) error {
		exportType, _ := cmd.Flags().GetString("type")
		format, _ := cmd.Flags().GetString("format")
		limit, _ := cmd.Flags().GetInt("limit")
		query, _ := cmd.Flags().GetString("query")

		opts := &oathnet.ExportCreateOptions{
			Format: format,
			Limit:  limit,
		}
		if query != "" {
			opts.Search = map[string]string{"query": query}
		}

		result, err := client.Exports.Create(exportType, opts)
		if err != nil {
			return err
		}

		output(result, func() {
			if result.Data != nil {
				fmt.Println("\nExport job created")
				fmt.Printf("Job ID: %s\n", result.Data.JobID)
				fmt.Printf("Status: %s\n", result.Data.Status)
				fmt.Printf("\nUse: oathnet export status %s\n", result.Data.JobID)
			}
		})
		return nil
	},
}

var exportStatusCmd = &cobra.Command{
	Use:   "status [job_id]",
	Short: "Get export job status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := client.Exports.GetStatus(args[0])
		if err != nil {
			return err
		}

		output(result, func() {
			if result.Data != nil {
				fmt.Println("\nExport Job Status")
				fmt.Printf("Job ID: %s\n", result.Data.JobID)
				fmt.Printf("Status: %s\n", result.Data.Status)
				if result.Data.Progress != nil {
					fmt.Printf("\nProgress:\n")
					fmt.Printf("  Progress: %.1f%%\n", result.Data.Progress.Percent)
					fmt.Printf("  Records: %d/%d\n", result.Data.Progress.RecordsDone, result.Data.Progress.RecordsTotal)
				}
				if result.Data.Result != nil && result.Data.Status == "completed" {
					fmt.Printf("\nResult:\n")
					fmt.Printf("  File: %s\n", result.Data.Result.FileName)
					fmt.Printf("  Size: %d bytes\n", result.Data.Result.FileSize)
					fmt.Printf("  Records: %d\n", result.Data.Result.Records)
				}
			}
		})
		return nil
	},
}

var exportDownloadCmd = &cobra.Command{
	Use:   "download [job_id]",
	Short: "Download completed export",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")
		fmt.Printf("Downloading export %s...\n", args[0])
		err := client.Exports.Download(args[0], output)
		if err != nil {
			return err
		}
		fmt.Printf("Downloaded to %s\n", output)
		return nil
	},
}

func init() {
	exportCreateCmd.Flags().String("type", "", "Export type (docs|victims) (required)")
	exportCreateCmd.Flags().String("format", "jsonl", "Output format (jsonl|csv)")
	exportCreateCmd.Flags().Int("limit", 0, "Maximum records")
	exportCreateCmd.Flags().StringP("query", "q", "", "Search query")
	exportCreateCmd.MarkFlagRequired("type")

	exportDownloadCmd.Flags().StringP("output", "o", "", "Output file path (required)")
	exportDownloadCmd.MarkFlagRequired("output")

	exportCmd.AddCommand(exportCreateCmd)
	exportCmd.AddCommand(exportStatusCmd)
	exportCmd.AddCommand(exportDownloadCmd)
}

// ============================================
// BULK COMMANDS
// ============================================

var bulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Bulk search commands",
}

var bulkCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a bulk search job",
	RunE: func(cmd *cobra.Command, args []string) error {
		terms, _ := cmd.Flags().GetStringSlice("term")
		service, _ := cmd.Flags().GetString("service")
		format, _ := cmd.Flags().GetString("format")

		result, err := client.Bulk.Create(terms, service, &oathnet.BulkCreateOptions{
			Format: format,
		})
		if err != nil {
			return err
		}

		output(result, func() {
			if result.Data != nil {
				fmt.Println("\nBulk job created")
				fmt.Printf("Job ID: %s\n", result.Data.JobID)
				fmt.Printf("Status: %s\n", result.Data.Status)
			}
		})
		return nil
	},
}

var bulkStatusCmd = &cobra.Command{
	Use:   "status [job_id]",
	Short: "Get bulk job status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := client.Bulk.GetStatus(args[0])
		if err != nil {
			return err
		}

		output(result, func() {
			if result.Data != nil {
				fmt.Println("\nBulk Job Status")
				fmt.Printf("Job ID: %s\n", result.Data.JobID)
				fmt.Printf("Status: %s\n", result.Data.Status)
				fmt.Printf("Terms: %d\n", result.Data.TermsCount)
				fmt.Printf("Results: %d\n", result.Data.ResultsCount)
			}
		})
		return nil
	},
}

var bulkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List bulk search jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := client.Bulk.List(1, 20)
		if err != nil {
			return err
		}

		output(result, func() {
			fmt.Println("\nBulk Jobs")
			if len(result.Results) > 0 {
				for _, j := range result.Results {
					fmt.Printf("  %s: %s (%d terms)\n", j.JobID, j.Status, j.TermsCount)
				}
			} else {
				fmt.Println("No jobs found")
			}
		})
		return nil
	},
}

func init() {
	bulkCreateCmd.Flags().StringSliceP("term", "t", nil, "Search terms (required)")
	bulkCreateCmd.Flags().String("service", "", "Service (breach|stealer) (required)")
	bulkCreateCmd.Flags().String("format", "jsonl", "Output format")
	bulkCreateCmd.MarkFlagRequired("term")
	bulkCreateCmd.MarkFlagRequired("service")

	bulkCmd.AddCommand(bulkCreateCmd)
	bulkCmd.AddCommand(bulkStatusCmd)
	bulkCmd.AddCommand(bulkListCmd)
}

// ============================================
// UTILITY COMMANDS
// ============================================

var utilCmd = &cobra.Command{
	Use:   "util",
	Short: "Utility commands",
}

var utilDBNamesCmd = &cobra.Command{
	Use:   "dbnames",
	Short: "Autocomplete database names",
	RunE: func(cmd *cobra.Command, args []string) error {
		query, _ := cmd.Flags().GetString("query")

		result, err := client.Utility.DBNameAutocomplete(query)
		if err != nil {
			return err
		}

		output(result, func() {
			fmt.Println("\nDatabase Names")
			for _, name := range result {
				fmt.Printf("  %s\n", name)
			}
		})
		return nil
	},
}

func init() {
	utilDBNamesCmd.Flags().StringP("query", "q", "", "Search query (required)")
	utilDBNamesCmd.MarkFlagRequired("query")

	utilCmd.AddCommand(utilDBNamesCmd)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
