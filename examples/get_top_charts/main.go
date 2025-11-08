package main

import (
	"context"
	"log"

	"github.com/truewebber/goitunes/v2/pkg/goitunes"
)

const (
	topFreeDisplayLimit  = 10
	topPaidDisplayLimit  = 5
	topPaidPageSize      = 50
	topGrossingRangeFrom = 50
	topGrossingRangeTo   = 10
)

func main() {
	// Create a new client for the US App Store
	client, err := goitunes.New("us")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Get top 200 free applications in all categories
	log.Println("=== Top 200 Free Applications ===")

	charts, err := client.Charts().GetTop200(
		ctx,
		goitunes.GenreAll,
		goitunes.ChartTypeTopFree,
	)
	if err != nil {
		log.Fatalf("Failed to get top charts: %v", err)
	}

	// Display first 10 results
	for i := range charts {
		if i >= topFreeDisplayLimit {
			break
		}

		log.Printf("%d. %s (%s) - Rating: %.1f",
			charts[i].Position,
			charts[i].App.Name,
			charts[i].App.BundleID,
			charts[i].App.Rating,
		)
	}

	// Get top 1500 paid applications with pagination
	log.Println("\n=== Top 1500 Paid Applications (Page 0) ===")

	charts1500, err := client.Charts().GetTop1500(
		ctx,
		goitunes.GenreAll,
		goitunes.ChartTypeTopPaid,
		0,               // page
		topPaidPageSize, // page size
	)
	if err != nil {
		log.Fatalf("Failed to get top 1500: %v", err)
	}

	// Display first 5 results
	for i := range charts1500 {
		if i >= topPaidDisplayLimit {
			break
		}

		log.Printf("%d. %s - $%.2f %s",
			charts1500[i].Position,
			charts1500[i].App.BundleID,
			charts1500[i].App.Price,
			charts1500[i].App.Currency,
		)
	}

	// Get top grossing applications with custom range
	log.Println("\n=== Top Grossing Applications (positions 50-60) ===")

	topGrossing, err := client.Charts().GetTop200(
		ctx,
		goitunes.GenreAll,
		goitunes.ChartTypeTopGrossing,
		goitunes.WithRange(topGrossingRangeFrom, topGrossingRangeTo),
	)
	if err != nil {
		log.Fatalf("Failed to get top grossing: %v", err)
	}

	for i := range topGrossing {
		log.Printf("%d. %s - %s",
			topGrossing[i].Position,
			topGrossing[i].App.Name,
			topGrossing[i].App.BundleID,
		)
	}
}
