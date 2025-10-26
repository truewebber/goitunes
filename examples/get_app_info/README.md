# Get Application Info Example

This example demonstrates how to retrieve detailed information about applications from the App Store.

## Features

- Get application information by Adam ID
- Get application information by Bundle ID
- Retrieve rating information
- Get overall rating from public API

## Usage

```bash
cd examples/get_app_info
go run main.go
```

## Code Explanation

### Get by Adam ID

```go
apps, err := client.Applications().GetByAdamID(ctx, "564177498")
```

Retrieve application information using Apple's internal Adam ID. You can pass multiple IDs.

### Get by Bundle ID

```go
apps, err := client.Applications().GetByBundleID(
    ctx,
    "com.facebook.Facebook",
    "com.instagram.android",
)
```

Retrieve application information using the bundle identifier. Supports multiple bundle IDs.

### Get Rating

```go
rating, err := client.Applications().GetRating(ctx, "564177498")
```

Get detailed rating information for an application.

### Get Overall Rating

```go
rating, err := client.Applications().GetOverallRating(ctx, "564177498")
```

Get overall rating from the public iTunes lookup API.

## Application Information

The returned `ApplicationDTO` contains:
- Adam ID and Bundle ID
- Application name and artist information
- Version and version ID
- Price and currency
- Rating and review count
- Release date
- Genre information
- Device families (iPhone, iPad)
- File size
- Minimum OS version
- Description
- Icon and screenshot URLs
- Helper flags (IsFree, IsUniversal)

