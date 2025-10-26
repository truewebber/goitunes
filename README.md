# GoiTunes

Go library for Apple App Store API.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Installation

```bash
go get github.com/truewebber/goitunes
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/truewebber/goitunes/pkg/goitunes"
)

func main() {
    // Create client for US App Store
    client, err := goitunes.New("us")
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Get top 10 free apps
    charts, err := client.Charts().GetTop200(
        ctx,
        goitunes.GenreAll,
        goitunes.ChartTypeTopFree,
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Top 10 apps in %s:\n", goitunes.GenreAll.Name())
    for i, item := range charts[:10] {
        fmt.Printf("%d. %s\n", i+1, item.App.Name)
    }
}
```

## API Overview

### Charts Service

Get top applications from App Store charts.

```go
// Get top 200 applications
charts, err := client.Charts().GetTop200(
    ctx,
    genre,        // Genre type (see Genre IDs section below)
    chartType,    // TopFree, TopPaid, TopGrossing
    options...    // Optional: WithRange(), WithKidPrefix()
)

// Get top 1500 applications with pagination
charts, err := client.Charts().GetTop1500(
    ctx,
    genre,        // Genre type
    chartType,
    page,         // 0-based page number
    pageSize,     // Results per page
)
```

**Chart Types:**
- `ChartTypeTopFree` - Top free applications
- `ChartTypeTopPaid` - Top paid applications
- `ChartTypeTopGrossing` - Top grossing applications

**Options:**
- `WithRange(from, limit)` - Get specific range of positions
- `WithKidPrefix(prefix)` - Filter by age band

### Genre IDs

The library provides a strongly-typed `Genre` enum for all App Store genres with 72 categories:

**Main Categories:**
```go
goitunes.GenreAll                // All categories
goitunes.GenreGames              // Games
goitunes.GenreBusiness           // Business
goitunes.GenreEducation          // Education
goitunes.GenreEntertainment      // Entertainment
goitunes.GenreFinance            // Finance
goitunes.GenreHealthFitness      // Health & Fitness
goitunes.GenreLifestyle          // Lifestyle
goitunes.GenreSocialNetworking   // Social Networking
goitunes.GenrePhotoVideo         // Photo & Video
// ... and 15 more main categories (25 total)
```

**Game Sub-genres:**
```go
goitunes.GenreGamesAction        // Action games
goitunes.GenreGamesAdventure     // Adventure games
goitunes.GenreGamesRacing        // Racing games
goitunes.GenreGamesStrategy      // Strategy games
// ... and 14 more game sub-genres (18 total)
```

**Magazines Sub-genres:**
```go
goitunes.GenreMagazinesNewsPolitics      // News & Politics
goitunes.GenreMagazinesFashionStyle      // Fashion & Style
goitunes.GenreMagazinesArtsPhotography   // Arts & Photography
// ... and 22 more magazine sub-genres (25 total)
```

**Kids Sub-genres:**
```go
goitunes.GenreKidsLess5          // Kids 5 & Under
goitunes.GenreKids6To8           // Kids 6–8
goitunes.GenreKids9To11          // Kids 9–11
```

**Genre Methods:**
```go
genre := goitunes.GenreGames
genre.String()   // Returns: "6014"
genre.Name()     // Returns: "Games"
genre.IsValid()  // Returns: true
```

**Example:**
```go
// Get top free games
charts, _ := client.Charts().GetTop200(ctx, goitunes.GenreGames, goitunes.ChartTypeTopFree)

// Get top action games
charts, _ := client.Charts().GetTop200(ctx, goitunes.GenreGamesAction, goitunes.ChartTypeTopFree)

// Custom genre with validation
customGenre := goitunes.Genre("7001")
if customGenre.IsValid() {
    fmt.Println(customGenre.Name()) // Output: "Action"
}
```

See [`pkg/goitunes/constants.go`](pkg/goitunes/constants.go) for the complete list of all 73 genre constants.

### Application Service

Get detailed information about applications.

```go
// Get by Adam ID (Apple's internal ID)
apps, err := client.Applications().GetByAdamID(ctx, "564177498", "284882215")

// Get by Bundle ID
apps, err := client.Applications().GetByBundleID(ctx, "com.facebook.Facebook")

// Get rating information
rating, err := client.Applications().GetRating(ctx, adamID)

// Get overall rating from public API
rating, err := client.Applications().GetOverallRating(ctx, adamID)
```

**Application Response includes:**
- Adam ID, Bundle ID, Name
- Version information
- Price and currency
- Rating and review count
- Release date
- Genre information
- Device families (iPhone, iPad)
- File size
- Screenshots and icon URLs
- Description

### Authentication Service

Authenticate with Apple ID to access purchase functionality.

```go
// Create client with credentials
client, err := goitunes.New("us",
    goitunes.WithCredentials(appleID, "", ""),
    goitunes.WithDevice(guid, machineName, userAgent),
)

// Login with password
authResp, err := client.Auth().Login(ctx, password)
// Returns: PasswordToken, DSID, Authenticated status

// Check authentication status
isAuth := client.IsAuthenticated()
```

### Purchase Service

Purchase and download applications (requires kbsync certificate).

```go
// Setup client with kbsync
client, err := goitunes.New("us",
    goitunes.WithCredentials(appleID, "", ""),
    goitunes.WithKbsync(kbsyncCertificate),
    goitunes.WithDevice(guid, machineName, userAgent),
)

// Login first
client.Auth().Login(ctx, password)

// Purchase application
downloadInfo, err := client.Purchase().Buy(ctx, adamID, versionID)
```

**Download Info includes:**
- Download URL
- Download key and headers
- SINF (DRM information)
- iTunes metadata plist
- Bundle ID and version

## Supported Regions

29 App Store regions are supported:

| Region | Code | Region | Code | Region | Code |
|--------|------|--------|------|--------|------|
| United States | `us` | Russia | `ru` | Japan | `jp` |
| United Kingdom | `gb` | Germany | `de` | China | `cn` |
| Canada | `ca` | France | `fr` | Hong Kong | `hk` |
| Australia | `au` | Italy | `it` | Taiwan | `tw` |
| Brazil | `br` | Spain | `es` | Korea | `kr` |
| Mexico | `mx` | Netherlands | `nl` | Singapore | `sg` |
| India | `in` | Turkey | `tr` | Malaysia | `my` |
| Argentina | `ar` | Portugal | `pt` | Indonesia | `id` |
| UAE | `ae` | South Africa | `za` | Thailand | `th` |
| New Zealand | `nz` | | | Vietnam | `vn` |

```go
// Get list of all supported regions
regions := client.SupportedRegions()

// Create clients for different regions
usClient, _ := goitunes.New("us")
jpClient, _ := goitunes.New("jp")
```

## Configuration Options

### Custom HTTP Client

Provide your own HTTP client for custom behavior (retry logic, logging, etc.).

```go
type MyHTTPClient struct {
    // Your implementation
}

func (c *MyHTTPClient) Do(req *http.Request) (*http.Response, error) {
    // Custom logic
    return http.DefaultClient.Do(req)
}

client, err := goitunes.New("us",
    goitunes.WithHTTPClient(&MyHTTPClient{}),
)
```

### Device Configuration

Set device information for authentication and purchases.

```go
client, err := goitunes.New("us",
    goitunes.WithDevice(
        guid,        // Device UUID
        machineName, // Machine name
        userAgent,   // User agent string
    ),
)
```

**Predefined User Agents:**
- `goitunes.UserAgentWindows`
- `goitunes.UserAgentTop200`
- `goitunes.UserAgentTop1500`
- `goitunes.UserAgentDownload`

## Examples

Complete working examples are available in the `examples/` directory:

- **[get_top_charts](examples/get_top_charts/)** - Retrieve top application charts
- **[get_app_info](examples/get_app_info/)** - Get detailed application information
- **[authentication](examples/authentication/)** - Login with Apple ID
- **[purchase_app](examples/purchase_app/)** - Purchase and download applications
- **[custom_http_client](examples/custom_http_client/)** - Custom HTTP client with retry logic
- **[multi_region](examples/multi_region/)** - Work with multiple regions

Each example includes a README with detailed explanation.

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./internal/...

# Run specific tests
go test ./internal/domain/entity/...
```

## Error Handling

The library uses standard Go error handling with context:

```go
charts, err := client.Charts().GetTop200(ctx, genreID, chartType)
if err != nil {
    if errors.Is(err, goitunes.ErrUnsupportedRegion) {
        // Handle unsupported region
    }
    // Handle other errors
}
```

**Defined Errors:**
- `ErrUnsupportedRegion` - Region not supported
- `ErrNotAuthenticated` - Authentication required
- `ErrInvalidCredentials` - Invalid credentials
- `ErrApplicationNotFound` - Application not found
- `ErrPurchaseFailed` - Purchase operation failed
- `ErrInvalidRequest` - Invalid request parameters

## Limitations & Notes

### Kbsync Certificate
The kbsync certificate is required for purchases but is difficult to obtain. It's a certificate that authorizes purchase operations (STDQ) or re-downloads (STDRDL).

### Rate Limiting
Apple may rate limit API requests. Consider implementing appropriate delays between requests.

### Authentication Tokens
Password tokens and DSID are temporary and will expire. Re-authentication may be required.

### Re-downloads
Re-downloading already purchased apps may require a different kbsync certificate (STDRDL vs STDQ).

## Architecture

The library is built with Clean Architecture principles:

- **Domain Layer** - Core entities and business logic
- **Application Layer** - Use cases and DTOs
- **Infrastructure Layer** - HTTP clients and API implementations
- **Public API** - Clean, service-based interface

This architecture ensures:
- Easy testing and mocking
- Clear separation of concerns
- Simple to extend and maintain
- Pluggable dependencies

## Contributing

Contributions are welcome! Please:
1. Follow Go best practices
2. Include tests for new features
3. Update documentation
4. Provide examples for new functionality

## License

MIT License

Copyright (c) 2025 truewebber

See [LICENSE](LICENSE) file for full details.

## Disclaimer

This library is for educational purposes. Using it to purchase applications may violate Apple's Terms of Service. Use at your own risk.
