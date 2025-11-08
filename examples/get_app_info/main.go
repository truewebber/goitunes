package main

import (
	"context"
	"log"
	"reflect"

	"github.com/truewebber/goitunes/v2/pkg/goitunes"
)

const (
	bytesPerKB = 1024
	bytesPerMB = bytesPerKB * bytesPerKB
)

func main() {
	client := createClient()
	ctx := context.Background()

	demonstrateGetByAdamID(ctx, client)
	demonstrateGetByBundleID(ctx, client)
	demonstrateGetRating(ctx, client)
	demonstrateGetOverallRating(ctx, client)
}

func createClient() *goitunes.Client {
	client, err := goitunes.New("us")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return client
}

func demonstrateGetByAdamID(ctx context.Context, client *goitunes.Client) {
	log.Println("=== Get Application by Adam ID ===")

	apps, err := client.Applications().GetByAdamID(ctx, "564177498") // Angry Birds
	if err != nil {
		log.Fatalf("Failed to get app info: %v", err)
	}

	if len(apps) > 0 {
		app := apps[0]
		logAppDetails(app)
	}
}

func logAppDetails(app interface{}) {
	// Use reflection to access exported fields
	v := reflect.ValueOf(app)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		log.Printf("App details: %+v", app)

		return
	}

	log.Printf("Name: %s", getStringField(v, "Name"))
	log.Printf("Bundle ID: %s", getStringField(v, "BundleID"))
	log.Printf("Artist: %s", getStringField(v, "ArtistName"))
	log.Printf("Version: %s", getStringField(v, "Version"))
	log.Printf("Price: $%.2f %s", getFloatField(v, "Price"), getStringField(v, "Currency"))
	log.Printf("Rating: %.1f (%d reviews)", getFloatField(v, "Rating"), getIntField(v, "RatingCount"))
	log.Printf("Genre: %s", getStringField(v, "GenreName"))
	log.Printf("File Size: %.2f MB", float64(getInt64Field(v, "FileSize"))/bytesPerMB)
	log.Printf("Minimum OS: %s", getStringField(v, "MinimumOSVersion"))
	log.Printf("Free: %t", getBoolField(v, "IsFree"))
	log.Printf("Universal: %t", getBoolField(v, "IsUniversal"))
	log.Printf("Icon URL: %s", getStringField(v, "IconURL"))
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
		return 0.0
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

func getInt64Field(v reflect.Value, field string) int64 {
	f := v.FieldByName(field)
	if !f.IsValid() || f.Kind() != reflect.Int64 {
		return 0
	}

	return f.Int()
}

func getBoolField(v reflect.Value, field string) bool {
	f := v.FieldByName(field)
	if !f.IsValid() || f.Kind() != reflect.Bool {
		return false
	}

	return f.Bool()
}

func demonstrateGetByBundleID(ctx context.Context, client *goitunes.Client) {
	log.Println("\n=== Get Multiple Applications by Bundle ID ===")

	apps, err := client.Applications().GetByBundleID(
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
}

func demonstrateGetRating(ctx context.Context, client *goitunes.Client) {
	log.Println("\n=== Get Rating Information ===")

	rating, err := client.Applications().GetRating(ctx, "564177498")
	if err != nil {
		log.Fatalf("Failed to get rating: %v", err)
	}

	log.Printf("Rating: %.2f", rating.Rating)
	log.Printf("Rating Count: %d", rating.RatingCount)
}

func demonstrateGetOverallRating(ctx context.Context, client *goitunes.Client) {
	log.Println("\n=== Get Overall Rating ===")

	overallRating, err := client.Applications().GetOverallRating(ctx, "564177498")
	if err != nil {
		log.Fatalf("Failed to get overall rating: %v", err)
	}

	log.Printf("Overall Rating: %.2f", overallRating.Rating)
	log.Printf("Overall Rating Count: %d", overallRating.RatingCount)
}
