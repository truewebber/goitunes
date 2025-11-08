package main

import (
	"context"
	"fmt"
	"log"

	"github.com/truewebber/goitunes/v2/pkg/goitunes"
)

const (
	topAppsLimit = 3
)

func main() {
	ctx := context.Background()

	compareAppPrices(ctx)
	compareTopApps(ctx)
	showSupportedRegions()
	logRegionalDifferences()
}

func compareAppPrices(ctx context.Context) {
	regions := []string{"us", "ru", "gb", "jp", "de"}
	bundleID := "com.facebook.Facebook"

	log.Println("=== Comparing App Prices Across Regions ===")

	for _, region := range regions {
		compareAppPriceInRegion(ctx, region, bundleID)
	}
}

func compareAppPriceInRegion(ctx context.Context, region, bundleID string) {
	client, err := goitunes.New(region)
	if err != nil {
		log.Printf("Failed to create client for %s: %v", region, err)

		return
	}

	apps, err := client.Applications().GetByBundleID(ctx, bundleID)
	if err != nil {
		log.Printf("Failed to get app for %s: %v", region, err)

		return
	}

	if len(apps) == 0 {
		log.Printf("%s: App not available", region)

		return
	}

	app := apps[0]

	// Format price directly - fields are exported and accessible
	price := "Free"
	if !app.IsFree {
		price = fmt.Sprintf("%.2f %s", app.Price, app.Currency)
	}

	log.Printf("%s: %s - %s (Rating: %.1f)",
		region,
		app.Name,
		price,
		app.Rating,
	)
}

func compareTopApps(ctx context.Context) {
	regionsTop := []string{"us", "jp", "br"}

	log.Println("\n=== Top Apps in Different Regions ===")

	for _, region := range regionsTop {
		showTopAppsInRegion(ctx, region)
	}
}

func showTopAppsInRegion(ctx context.Context, region string) {
	client, err := goitunes.New(region)
	if err != nil {
		log.Printf("Failed to create client for %s: %v", region, err)

		return
	}

	charts, err := client.Charts().GetTop200(
		ctx,
		goitunes.GenreAll,
		goitunes.ChartTypeTopFree,
		goitunes.WithRange(1, topAppsLimit),
	)
	if err != nil {
		log.Printf("Failed to get charts for %s: %v", region, err)

		return
	}

	log.Printf("Top 3 Free Apps in %s:", region)

	for i := range charts {
		log.Printf("  %d. %s", charts[i].Position, charts[i].App.Name)
	}

	log.Println()
}

func showSupportedRegions() {
	log.Println("=== Supported Regions ===")

	client, err := goitunes.New("us")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	regions := client.SupportedRegions()

	log.Printf("Total supported regions: %d", len(regions))
	log.Printf("Regions: %v", regions)
}

func logRegionalDifferences() {
	log.Println("\n=== Regional Differences ===")
	log.Println("✓ Pricing varies by region")
	log.Println("✓ Chart rankings differ")
	log.Println("✓ App availability varies")
	log.Println("✓ Currency and localization")
	log.Println("✓ Rating counts can differ")
}
