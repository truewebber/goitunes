package main

import (
	"context"
	"fmt"
	"log"

	"github.com/truewebber/goitunes/pkg/goitunes"
)

func main() {
	ctx := context.Background()

	// Define regions to query
	regions := []string{"us", "ru", "gb", "jp", "de"}

	fmt.Println("=== Comparing App Prices Across Regions ===")

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
			fmt.Printf("%s: App not available\n", region)
			continue
		}

		app := apps[0]
		price := "Free"
		if !app.IsFree {
			price = fmt.Sprintf("%.2f %s", app.Price, app.Currency)
		}

		fmt.Printf("%s: %s - %s (Rating: %.1f)\n",
			region,
			app.Name,
			price,
			app.Rating,
		)
	}

	fmt.Println("\n=== Top Apps in Different Regions ===")

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
			goitunes.WithRange(1, 3),
		)
		if err != nil {
			log.Printf("Failed to get charts for %s: %v", region, err)
			continue
		}

		fmt.Printf("Top 3 Free Apps in %s:\n", region)
		for _, item := range charts {
			fmt.Printf("  %d. %s\n", item.Position, item.App.Name)
		}
		fmt.Println()
	}

	fmt.Println("=== Supported Regions ===")

	// Get list of all supported regions
	client, _ := goitunes.New("us")
	regions = client.SupportedRegions()

	fmt.Printf("Total supported regions: %d\n", len(regions))
	fmt.Printf("Regions: %v\n", regions)

	fmt.Println("\n=== Regional Differences ===")
	fmt.Println("✓ Pricing varies by region")
	fmt.Println("✓ Chart rankings differ")
	fmt.Println("✓ App availability varies")
	fmt.Println("✓ Currency and localization")
	fmt.Println("✓ Rating counts can differ")
}
