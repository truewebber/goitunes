# Multi-Region Example

This example demonstrates how to work with multiple App Store regions simultaneously.

## Features

- Query multiple regions
- Compare prices across regions
- View regional chart differences
- Handle regional availability
- List all supported regions

## Usage

```bash
cd examples/multi_region
go run main.go
```

## Code Explanation

### Creating Clients for Different Regions

```go
usClient, _ := goitunes.New("us")
ukClient, _ := goitunes.New("gb")
jpClient, _ := goitunes.New("jp")
```

Each client is configured for a specific App Store region.

### Comparing Prices

```go
for _, region := range regions {
    client, _ := goitunes.New(region)
    apps, _ := client.Applications().GetByBundleID(ctx, bundleID)
    fmt.Printf("%s: %.2f %s\n", region, apps[0].Price, apps[0].Currency)
}
```

Get pricing information for the same app across different regions.

### Regional Charts

```go
charts, _ := client.Charts().GetTop200(ctx, goitunes.GenreAll, goitunes.ChartTypeTopFree)
```

Chart rankings differ significantly between regions based on local preferences.

### Listing Supported Regions

```go
regions := client.SupportedRegions()
```

Get a list of all available App Store regions.

## Supported Regions

The library supports 29 App Store regions:

- **Americas**: us, ca, mx, br, ar
- **Europe**: gb, fr, de, it, es, pt, nl, ru, tr
- **Asia-Pacific**: jp, kr, cn, hk, tw, au, nz, sg, my, id, th, vn, in
- **Middle East/Africa**: ae, za

## Regional Considerations

### 1. Pricing
- Prices vary by region
- Different currencies
- Local tax considerations
- Regional pricing strategies

### 2. Availability
- Not all apps available in all regions
- Regional restrictions
- Content licensing issues
- Local regulations

### 3. Charts
- Rankings differ by region
- Local preferences
- Cultural differences
- Seasonal variations

### 4. Ratings
- Rating counts vary
- Different review systems
- Language-specific reviews
- Regional user bases

### 5. Metadata
- Localized names and descriptions
- Different screenshots
- Regional marketing
- Language support

## Use Cases

### Price Comparison Service
Compare app prices across regions to find the best deals.

### Market Research
Analyze app performance in different markets:
- Regional popularity
- Pricing strategies
- Competitive analysis
- Market trends

### Availability Checker
Check if apps are available in specific regions.

### Multi-Region Publishing
Monitor your app's performance across all regions:
- Track rankings
- Monitor reviews
- Analyze pricing
- Compare metrics

## Example Output

```
=== Comparing App Prices Across Regions ===

us: Facebook - Free (Rating: 4.5)
ru: Facebook - Free (Rating: 4.3)
gb: Facebook - Free (Rating: 4.4)
jp: Facebook - Free (Rating: 4.2)
de: Facebook - Free (Rating: 4.4)

=== Top Apps in Different Regions ===

Top 3 Free Apps in us:
  1. TikTok
  2. YouTube
  3. Instagram

Top 3 Free Apps in jp:
  1. LINE
  2. Yahoo! JAPAN
  3. YouTube

Top 3 Free Apps in br:
  1. WhatsApp
  2. Instagram
  3. TikTok
```

## Best Practices

1. **Caching**: Cache region-specific data to reduce API calls
2. **Error Handling**: Some regions may have temporary issues
3. **Rate Limiting**: Be mindful when querying many regions
4. **Currency Conversion**: Consider exchange rates for price comparison
5. **Localization**: Handle different languages and formats

