package main

import (
	"context"
	"fmt"
	"log"

	"github.com/truewebber/goitunes/pkg/goitunes"
)

func main() {
	// Create a new client for the US App Store
	client, err := goitunes.New("us")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Get application info by Adam ID
	fmt.Println("=== Get Application by Adam ID ===")
	apps, err := client.Applications().GetByAdamID(ctx, "564177498") // Angry Birds
	if err != nil {
		log.Fatalf("Failed to get app info: %v", err)
	}

	if len(apps) > 0 {
		app := apps[0]
		fmt.Printf("Name: %s\n", app.Name)
		fmt.Printf("Bundle ID: %s\n", app.BundleID)
		fmt.Printf("Artist: %s\n", app.ArtistName)
		fmt.Printf("Version: %s\n", app.Version)
		fmt.Printf("Price: $%.2f %s\n", app.Price, app.Currency)
		fmt.Printf("Rating: %.1f (%d reviews)\n", app.Rating, app.RatingCount)
		fmt.Printf("Genre: %s\n", app.GenreName)
		fmt.Printf("File Size: %.2f MB\n", float64(app.FileSize)/(1024*1024))
		fmt.Printf("Minimum OS: %s\n", app.MinimumOSVersion)
		fmt.Printf("Free: %t\n", app.IsFree)
		fmt.Printf("Universal: %t\n", app.IsUniversal)
		fmt.Printf("Icon URL: %s\n", app.IconURL)
	}

	// Get multiple applications by Bundle ID
	fmt.Println("\n=== Get Multiple Applications by Bundle ID ===")
	apps, err = client.Applications().GetByBundleID(
		ctx,
		"com.facebook.Facebook",
		"com.instagram.android",
		"com.twitter.android",
	)
	if err != nil {
		log.Fatalf("Failed to get apps: %v", err)
	}

	for _, app := range apps {
		fmt.Printf("- %s (%s) - Rating: %.1f\n", app.Name, app.BundleID, app.Rating)
	}

	// Get rating information
	fmt.Println("\n=== Get Rating Information ===")
	rating, err := client.Applications().GetRating(ctx, "564177498")
	if err != nil {
		log.Fatalf("Failed to get rating: %v", err)
	}

	fmt.Printf("Rating: %.2f\n", rating.Rating)
	fmt.Printf("Rating Count: %d\n", rating.RatingCount)

	// Get overall rating from open API
	fmt.Println("\n=== Get Overall Rating ===")
	overallRating, err := client.Applications().GetOverallRating(ctx, "564177498")
	if err != nil {
		log.Fatalf("Failed to get overall rating: %v", err)
	}

	fmt.Printf("Overall Rating: %.2f\n", overallRating.Rating)
	fmt.Printf("Overall Rating Count: %d\n", overallRating.RatingCount)
}

