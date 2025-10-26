# Get Top Charts Example

This example demonstrates how to retrieve top application charts from the App Store.

## Features

- Get top 200 applications (free, paid, or grossing)
- Get top 1500 applications with pagination
- Filter by genre
- Custom ranges and options

## Usage

```bash
cd examples/get_top_charts
go run main.go
```

## Code Explanation

### Creating a Client

```go
client, err := goitunes.New("us")
```

Create a client for a specific App Store region (e.g., "us", "ru", "gb").

### Getting Top 200 Charts

```go
charts, err := client.Charts().GetTop200(
    ctx,
    goitunes.GenreAll,
    goitunes.ChartTypeTopFree,
)
```

Retrieve the top 200 applications. Available chart types:
- `ChartTypeTopFree` - Top free applications
- `ChartTypeTopPaid` - Top paid applications
- `ChartTypeTopGrossing` - Top grossing applications

### Getting Top 1500 Charts

```go
charts, err := client.Charts().GetTop1500(
    ctx,
    goitunes.GenreAll,
    goitunes.ChartTypeTopPaid,
    0,  // page number (0-based)
    50, // page size
)
```

Retrieve up to 1500 applications using pagination.

### Custom Options

```go
charts, err := client.Charts().GetTop200(
    ctx,
    goitunes.GenreAll,
    goitunes.ChartTypeTopFree,
    goitunes.WithRange(50, 10), // Get positions 50-60
)
```

Use options to customize the request.

