# OathNet Go SDK

Official Go SDK and CLI for the OathNet API.

## Installation

```bash
go get github.com/oathnet/oathnet-go
```

Or clone and build from source:

```bash
git clone https://github.com/oathnet/oathnet-go
cd oathnet-go
go build -o oathnet ./cmd/oathnet
```

## Quick Start

### SDK Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/oathnet/oathnet-go/pkg/oathnet"
)

func main() {
    client, err := oathnet.NewClient("your-api-key")
    if err != nil {
        log.Fatal(err)
    }

    // Search breach database
    result, err := client.Search.Breach("winterfox", nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d results\n", result.Data.ResultsFound)

    // Get IP info
    ipInfo, err := client.OSINT.IPInfo("174.235.65.156")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Country: %s, City: %s\n", ipInfo.Data.Country, ipInfo.Data.City)

    // Discord lookup
    discord, err := client.OSINT.DiscordUserinfo("300760994454437890")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Username: %s\n", discord.Data.Username)
}
```

### CLI Usage

```bash
# Build the CLI
go build -o oathnet ./cmd/oathnet

# Set API key
export OATHNET_API_KEY=your-api-key

# Search breach database
./oathnet search breach -q "winterfox"

# Get IP info
./oathnet osint ip 174.235.65.156

# Discord lookup
./oathnet osint discord user 300760994454437890

# Export to JSON
./oathnet search breach -q "winterfox" -f json
```

## Features

### Search Services
- **Breach Search**: Search leaked credentials across 50B+ records
- **Stealer Search**: Search stealer log databases
- **Search Sessions**: Optimize quota with grouped lookups

### V2 Services
- **V2 Stealer**: Enhanced stealer search with filtering
- **V2 Victims**: Search victim profiles with device info
- **V2 File Search**: Regex search within victim files
- **V2 Exports**: Export results to CSV/JSONL

### OSINT Lookups
- Discord (user info, username history, linked Roblox)
- Steam profiles
- Xbox Live profiles
- Roblox user info
- IP geolocation
- Email existence (Holehe)
- Google accounts (GHunt)
- Subdomain extraction
- Minecraft username history

## SDK Reference

### Client

```go
import "github.com/oathnet/oathnet-go/pkg/oathnet"

// Create client with default options
client, err := oathnet.NewClient("your-api-key")

// Create client with custom options
client, err := oathnet.NewClient("your-api-key",
    oathnet.WithBaseURL("https://oathnet.org/api"),
    oathnet.WithTimeout(60 * time.Second),
)
```

### Services

#### Search

```go
// Initialize search session
session, err := client.Search.InitSession("query")

// Search breach database
result, err := client.Search.Breach("query", &oathnet.BreachSearchOptions{
    Page:   1,
    Limit:  25,
    DBName: "linkedin",  // Optional database filter
})

// Search stealer database
result, err := client.Search.Stealer("query", nil)
```

#### OSINT

```go
// IP lookup
ipInfo, err := client.OSINT.IPInfo("8.8.8.8")

// Steam profile
steam, err := client.OSINT.Steam("steam_id")

// Xbox profile
xbox, err := client.OSINT.Xbox("gamertag")

// Discord user
discord, err := client.OSINT.DiscordUserinfo("discord_id")

// Discord username history
history, err := client.OSINT.DiscordUsernameHistory("discord_id")

// Discord to Roblox
roblox, err := client.OSINT.DiscordToRoblox("discord_id")

// Roblox user (by username)
robloxUser, err := client.OSINT.RobloxUserinfo("", "username")

// Holehe email check
holehe, err := client.OSINT.Holehe("email@example.com")

// Subdomain extraction
subdomains, err := client.OSINT.ExtractSubdomain("example.com")
```

#### V2 Stealer

```go
// Enhanced search with filters
result, err := client.Stealer.Search("query", &oathnet.StealerSearchOptions{
    Domain:   "facebook.com",
    HasLogID: true,
    Wildcard: true,
    PageSize: 25,
    Cursor:   "next_page_cursor",
})

// Extract subdomains from stealer data
subs, err := client.Stealer.Subdomain("example.com", "")
```

#### V2 Victims

```go
// Search victim profiles
result, err := client.Victims.Search("query", &oathnet.VictimsSearchOptions{
    Email:    "user@gmail.com",
    PageSize: 25,
})

// Get file manifest
manifest, err := client.Victims.GetManifest("log_id")

// Download file
content, err := client.Victims.GetFile("log_id", "file_id")

// Download archive
err := client.Victims.DownloadArchive("log_id", "output.zip")
```

#### File Search (Async)

```go
// Create search job
job, err := client.FileSearch.Create("password", &oathnet.FileSearchCreateOptions{
    SearchMode: "regex",  // "literal", "regex", "wildcard"
    MaxMatches: 100,
})

// Wait for results
result, err := client.FileSearch.WaitForCompletion(
    job.Data.JobID,
    2*time.Second,  // poll interval
    60*time.Second, // timeout
)

// Or use convenience method
result, err := client.FileSearch.Search("api_key", &oathnet.FileSearchCreateOptions{
    SearchMode: "literal",
    MaxMatches: 50,
}, 60*time.Second)
```

#### Exports (Async)

```go
// Create export
job, err := client.Exports.Create(&oathnet.ExportCreateOptions{
    ExportType: "docs",  // "docs" or "victims"
    Format:     "csv",   // "csv" or "jsonl"
    Limit:      1000,
    Fields:     []string{"email", "password", "domain"},
    Search:     map[string]interface{}{"query": "example.com"},
})

// Wait for completion
result, err := client.Exports.WaitForCompletion(
    job.Data.JobID,
    time.Second,
    120*time.Second,
)

// Download as bytes
data, err := client.Exports.Download(job.Data.JobID)

// Or download to file
err := client.Exports.DownloadToFile(job.Data.JobID, "export.csv")
```

### Error Handling

```go
import "github.com/oathnet/oathnet-go/pkg/oathnet"

result, err := client.Search.Breach("query", nil)
if err != nil {
    switch e := err.(type) {
    case *oathnet.AuthenticationError:
        log.Printf("Invalid API key: %v", e)
    case *oathnet.ValidationError:
        log.Printf("Invalid parameters: %v", e)
    case *oathnet.NotFoundError:
        log.Printf("Resource not found: %v", e)
    case *oathnet.RateLimitError:
        log.Printf("Rate limited: %v", e)
    case *oathnet.OathNetError:
        log.Printf("API error: %v", e)
    default:
        log.Printf("Unknown error: %v", err)
    }
}
```

## CLI Reference

```bash
# Global options
oathnet --api-key KEY --format json|table|raw COMMAND

# Search commands
oathnet search breach -q "query" [--page N] [--limit N] [--dbnames name]
oathnet search stealer -q "query"
oathnet search init -q "query"

# Stealer V2
oathnet stealer search -q "query" [--domain] [--wildcard] [--has-log-id]
oathnet stealer subdomain -d "domain.com"

# Victims V2
oathnet victims search -q "query" [--email] [--ip]
oathnet victims manifest LOG_ID
oathnet victims file LOG_ID FILE_ID
oathnet victims archive LOG_ID

# File Search
oathnet file-search create -e "expression" [--mode literal|regex|wildcard]
oathnet file-search status JOB_ID
oathnet file-search search -e "expression" [--timeout 60s]

# Exports
oathnet export create --type docs --format csv [--limit 1000]
oathnet export status JOB_ID
oathnet export download JOB_ID -o output.csv

# OSINT
oathnet osint ip IP_ADDRESS
oathnet osint steam STEAM_ID
oathnet osint xbox GAMERTAG
oathnet osint discord user DISCORD_ID
oathnet osint discord history DISCORD_ID
oathnet osint discord roblox DISCORD_ID
oathnet osint roblox [--user-id ID | --username NAME]
oathnet osint holehe EMAIL
oathnet osint ghunt EMAIL
oathnet osint subdomain DOMAIN
oathnet osint minecraft USERNAME

# Utility
oathnet util dbnames -q "linked"
```

## Configuration

API key can be set via:

1. CLI flag: `--api-key KEY`
2. Environment variable: `OATHNET_API_KEY`

## Examples

The `examples/` directory contains comprehensive examples:

| Example | Description |
|---------|-------------|
| `basic_usage/` | Client initialization and simple search |
| `breach_search/` | Breach search with filters and pagination |
| `stealer_search/` | V2 stealer search with log ID access |
| `victims_search/` | Victim profiles, manifests, and file access |
| `file_search/` | Async file search within victim logs |
| `osint_lookups/` | All OSINT methods demonstrated |
| `exports/` | Async export to CSV/JSONL |
| `error_handling/` | Error type handling and retry logic |
| `pagination/` | Cursor-based pagination patterns |

Run an example:

```bash
export OATHNET_API_KEY="your-api-key"
cd examples/basic_usage
go run main.go
```

## Development

```bash
# Run all tests (requires API key)
export OATHNET_API_KEY="your-api-key"
go test ./pkg/oathnet/... -v

# Run specific test file
go test ./pkg/oathnet/... -v -run TestSearchService

# Run tests with coverage
go test ./pkg/oathnet/... -v -cover -coverprofile=coverage.out
go tool cover -html=coverage.out

# Build CLI
go build -o oathnet ./cmd/oathnet
```

### Test Structure

```
pkg/oathnet/
  testing_helper_test.go  # Test utilities and constants
  client_test.go          # Client initialization tests
  errors_test.go          # Error type tests
  search_test.go          # Search service tests
  osint_test.go           # OSINT service tests
  stealer_test.go         # V2 stealer tests
  victims_test.go         # V2 victims tests
  filesearch_test.go      # File search tests
  exports_test.go         # Export tests
  utility_test.go         # Utility tests
```

## License

MIT License - See LICENSE file for details.

## Support

- Documentation: https://docs.oathnet.org
- Discord: https://discord.gg/DCjnk9TAMK
- Email: info@oathnet.org
