package main

import (
	"context"
	"log"
	"reflect"

	"github.com/truewebber/goitunes/v2/pkg/goitunes"
)

const (
	topFreeDisplayLimit  = 10
	topPaidDisplayLimit  = 5
	topPaidPageSize      = 50
	topGrossingRangeFrom = 50
	topGrossingRangeTo   = 10
	defaultFloatValue    = 0.0
)

func main() {
	client := createClient()
	ctx := context.Background()

	demonstrateTop200Free(ctx, client)
	demonstrateTop1500Paid(ctx, client)
	demonstrateTopGrossing(ctx, client)
}

func createClient() *goitunes.Client {
	client, err := goitunes.New("us")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return client
}

func demonstrateTop200Free(ctx context.Context, client *goitunes.Client) {
	log.Println("=== Top 200 Free Applications ===")

	charts, err := client.Charts().GetTop200(
		ctx,
		goitunes.GenreAll,
		goitunes.ChartTypeTopFree,
	)
	if err != nil {
		log.Fatalf("Failed to get top charts: %v", err)
	}

	displayCharts(charts, topFreeDisplayLimit, func(chart interface{}) {
		logChartItem(chart, func(pos int, appName, bundleID string, rating float64) {
			log.Printf("%d. %s (%s) - Rating: %.1f", pos, appName, bundleID, rating)
		})
	})
}

func demonstrateTop1500Paid(ctx context.Context, client *goitunes.Client) {
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

	displayCharts(charts1500, topPaidDisplayLimit, func(chart interface{}) {
		logChartItem(chart, func(pos int, _, bundleID string, _ float64) {
			log.Printf("%d. %s - $%.2f %s", pos, bundleID, getPrice(chart), getCurrency(chart))
		})
	})
}

func demonstrateTopGrossing(ctx context.Context, client *goitunes.Client) {
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

	displayCharts(topGrossing, len(topGrossing), func(chart interface{}) {
		logChartItem(chart, func(pos int, appName, bundleID string, _ float64) {
			log.Printf("%d. %s - %s", pos, appName, bundleID)
		})
	})
}

func displayCharts(charts interface{}, limit int, logFn func(interface{})) {
	// Use reflection to iterate over slice
	v := reflect.ValueOf(charts)
	if v.Kind() != reflect.Slice {
		return
	}

	for i := 0; i < v.Len() && i < limit; i++ {
		logFn(v.Index(i).Interface())
	}
}

func logChartItem(chart interface{}, logFn func(int, string, string, float64)) {
	v := reflect.ValueOf(chart)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return
	}

	pos := getIntField(v, "Position")
	appField := v.FieldByName("App")

	var appName, bundleID string

	var rating float64

	if appField.IsValid() {
		appName = getStringField(appField, "Name")
		bundleID = getStringField(appField, "BundleID")
		rating = getFloatField(appField, "Rating")
	}

	logFn(pos, appName, bundleID, rating)
}

func getStringField(v reflect.Value, field string) string {
	f := v.FieldByName(field)
	if !f.IsValid() || f.Kind() != reflect.String {
		return ""
	}

	return f.String()
}

func getFloatField(v reflect.Value, field string) float64 {
	f := v.FieldByName(field)
	if !f.IsValid() || f.Kind() != reflect.Float64 {
		return defaultFloatValue
	}

	return f.Float()
}

func getIntField(v reflect.Value, field string) int {
	f := v.FieldByName(field)
	if !f.IsValid() || f.Kind() != reflect.Int {
		return 0
	}

	return int(f.Int())
}

func getPrice(chart interface{}) float64 {
	v := reflect.ValueOf(chart)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	appField := v.FieldByName("App")
	if appField.IsValid() {
		return getFloatField(appField, "Price")
	}

	return defaultFloatValue
}

func getCurrency(chart interface{}) string {
	v := reflect.ValueOf(chart)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	appField := v.FieldByName("App")
	if appField.IsValid() {
		return getStringField(appField, "Currency")
	}

	return ""
}
