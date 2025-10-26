package main

import (
	"context"
	"fmt"
	"log"

	"github.com/truewebber/goitunes/v2/pkg/goitunes"
)

func main() {
	// Create a new client for the US App Store
	client, err := goitunes.New("us")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Get top 200 free applications in all categories
	fmt.Println("=== Top 200 Free Applications ===")
	charts, err := client.Charts().GetTop200(
		ctx,
		goitunes.GenreAll,
		goitunes.ChartTypeTopFree,
	)
	if err != nil {
		log.Fatalf("Failed to get top charts: %v", err)
	}

	// Display first 10 results
	for i, item := range charts {
		if i >= 10 {
			break
		}
		fmt.Printf("%d. %s (%s) - Rating: %.1f\n",
			item.Position,
			item.App.Name,
			item.App.BundleID,
			item.App.Rating,
		)
	}

	// Get top 1500 paid applications with pagination
	fmt.Println("\n=== Top 1500 Paid Applications (Page 0) ===")
	charts1500, err := client.Charts().GetTop1500(
		ctx,
		goitunes.GenreAll,
		goitunes.ChartTypeTopPaid,
		0,  // page
		50, // page size
	)
	if err != nil {
		log.Fatalf("Failed to get top 1500: %v", err)
	}

	// Display first 5 results
	for i, item := range charts1500 {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s - $%.2f %s\n",
			item.Position,
			item.App.BundleID,
			item.App.Price,
			item.App.Currency,
		)
	}

	// Get top grossing applications with custom range
	fmt.Println("\n=== Top Grossing Applications (positions 50-60) ===")
	topGrossing, err := client.Charts().GetTop200(
		ctx,
		goitunes.GenreAll,
		goitunes.ChartTypeTopGrossing,
		goitunes.WithRange(50, 10),
	)
	if err != nil {
		log.Fatalf("Failed to get top grossing: %v", err)
	}

	for _, item := range topGrossing {
		fmt.Printf("%d. %s - %s\n",
			item.Position,
			item.App.Name,
			item.App.BundleID,
		)
	}
}
