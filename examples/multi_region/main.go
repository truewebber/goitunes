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

	// Define regions to query
	regions := []string{"us", "ru", "gb", "jp", "de"}

	log.Println("=== Comparing App Prices Across Regions ===")

	bundleID := "com.facebook.Facebook"

	for _, region := range regions {
		// Create a client for each region
		client, err := goitunes.New(region)
		if err != nil {
			log.Printf("Failed to create client for %s: %v", region, err)

			continue
		}

		// Get app info
		apps, err := client.Applications().GetByBundleID(ctx, bundleID)
		if err != nil {
			log.Printf("Failed to get app for %s: %v", region, err)

			continue
		}

		if len(apps) == 0 {
			log.Printf("%s: App not available", region)

			continue
		}

		app := apps[0]

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

	log.Println("\n=== Top Apps in Different Regions ===")

	regionsTop := []string{"us", "jp", "br"}

	for _, region := range regionsTop {
		client, err := goitunes.New(region)
		if err != nil {
			log.Printf("Failed to create client for %s: %v", region, err)

			continue
		}

		charts, err := client.Charts().GetTop200(
			ctx,
			goitunes.GenreAll,
			goitunes.ChartTypeTopFree,
			goitunes.WithRange(1, topAppsLimit),
		)
		if err != nil {
			log.Printf("Failed to get charts for %s: %v", region, err)

			continue
		}

		log.Printf("Top 3 Free Apps in %s:", region)

		for i := range charts {
			log.Printf("  %d. %s", charts[i].Position, charts[i].App.Name)
		}

		log.Println()
	}

	log.Println("=== Supported Regions ===")

	// Get list of all supported regions
	client, err := goitunes.New("us")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	regions = client.SupportedRegions()

	log.Printf("Total supported regions: %d", len(regions))
	log.Printf("Regions: %v", regions)

	log.Println("\n=== Regional Differences ===")
	log.Println("✓ Pricing varies by region")
	log.Println("✓ Chart rankings differ")
	log.Println("✓ App availability varies")
	log.Println("✓ Currency and localization")
	log.Println("✓ Rating counts can differ")
}
