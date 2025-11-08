package main

import (
	"context"
	"log"

	"github.com/truewebber/goitunes/v2/pkg/goitunes"
)

const (
	bytesPerKB = 1024
	bytesPerMB = bytesPerKB * bytesPerKB
)

func main() {
	// Create a new client for the US App Store
	client, err := goitunes.New("us")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Get application info by Adam ID
	log.Println("=== Get Application by Adam ID ===")

	apps, err := client.Applications().GetByAdamID(ctx, "564177498") // Angry Birds
	if err != nil {
		log.Fatalf("Failed to get app info: %v", err)
	}

	if len(apps) > 0 {
		app := apps[0]
		log.Printf("Name: %s", app.Name)
		log.Printf("Bundle ID: %s", app.BundleID)
		log.Printf("Artist: %s", app.ArtistName)
		log.Printf("Version: %s", app.Version)
		log.Printf("Price: $%.2f %s", app.Price, app.Currency)
		log.Printf("Rating: %.1f (%d reviews)", app.Rating, app.RatingCount)
		log.Printf("Genre: %s", app.GenreName)
		log.Printf("File Size: %.2f MB", float64(app.FileSize)/bytesPerMB)
		log.Printf("Minimum OS: %s", app.MinimumOSVersion)
		log.Printf("Free: %t", app.IsFree)
		log.Printf("Universal: %t", app.IsUniversal)
		log.Printf("Icon URL: %s", app.IconURL)
	}

	// Get multiple applications by Bundle ID
	log.Println("\n=== Get Multiple Applications by Bundle ID ===")

	apps, err = client.Applications().GetByBundleID(
		ctx,
		"com.facebook.Facebook",
		"com.instagram.android",
		"com.twitter.android",
	)
	if err != nil {
		log.Fatalf("Failed to get apps: %v", err)
	}

	for i := range apps {
		log.Printf("- %s (%s) - Rating: %.1f", apps[i].Name, apps[i].BundleID, apps[i].Rating)
	}

	// Get rating information
	log.Println("\n=== Get Rating Information ===")

	rating, err := client.Applications().GetRating(ctx, "564177498")
	if err != nil {
		log.Fatalf("Failed to get rating: %v", err)
	}

	log.Printf("Rating: %.2f", rating.Rating)
	log.Printf("Rating Count: %d", rating.RatingCount)

	// Get overall rating from open API
	log.Println("\n=== Get Overall Rating ===")

	overallRating, err := client.Applications().GetOverallRating(ctx, "564177498")
	if err != nil {
		log.Fatalf("Failed to get overall rating: %v", err)
	}

	log.Printf("Overall Rating: %.2f", overallRating.Rating)
	log.Printf("Overall Rating Count: %d", overallRating.RatingCount)
}
