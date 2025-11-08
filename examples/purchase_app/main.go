package main

import (
	"context"
	"log"
	"os"
	"reflect"

	"github.com/truewebber/goitunes/v2/pkg/goitunes"
)

const (
	bytesPerMB       = 1024 * 1024
	maxStringDisplay = 50
)

func main() {
	appleID, password, kbsync, guid, machineName := getEnvVars()
	client := createClient(appleID, kbsync, guid, machineName)
	ctx := context.Background()

	performLogin(ctx, client, password)
	checkPurchaseCapability(client)

	app := getAppInfo(ctx, client)
	purchaseApp(ctx, client, app)
}

//nolint:nonamedreturns // gocritic requires named returns
func getEnvVars() (appleID, password, kbsync, guid, machineName string) {
	appleID = os.Getenv("APPLE_ID")
	password = os.Getenv("APPLE_PASSWORD")
	kbsync = os.Getenv("KBSYNC")
	guid = os.Getenv("DEVICE_GUID")
	machineName = os.Getenv("MACHINE_NAME")

	if appleID == "" || password == "" || kbsync == "" {
		log.Fatal("APPLE_ID, APPLE_PASSWORD, and KBSYNC environment variables are required")
	}

	if guid == "" {
		guid = "00000000-0000-0000-0000-000000000000"
	}

	if machineName == "" {
		machineName = "MyMachine"
	}

	return appleID, password, kbsync, guid, machineName
}

func createClient(appleID, kbsync, guid, machineName string) *goitunes.Client {
	client, err := goitunes.New("us",
		goitunes.WithCredentials(appleID, "", ""),
		goitunes.WithKbsync(kbsync),
		goitunes.WithDevice(guid, machineName, goitunes.UserAgentWindows),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return client
}

func performLogin(ctx context.Context, client *goitunes.Client, password string) {
	log.Println("=== Logging in ===")

	authResp, err := client.Auth().Login(ctx, password)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	log.Printf("Login successful! DSID: %s", authResp.DSID)
}

func checkPurchaseCapability(client *goitunes.Client) {
	if !client.CanPurchase() {
		log.Fatal("Cannot purchase: missing credentials or kbsync")
	}
}

func getAppInfo(ctx context.Context, client *goitunes.Client) interface{} {
	bundleID := "com.example.app" // Replace with actual bundle ID
	log.Printf("\n=== Getting application info for %s ===", bundleID)

	apps, err := client.Applications().GetByBundleID(ctx, bundleID)
	if err != nil {
		log.Fatalf("Failed to get app info: %v", err)
	}

	if len(apps) == 0 {
		log.Fatal("Application not found")
	}

	app := apps[0]

	// Use reflection to access fields
	v := reflect.ValueOf(app)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	log.Printf("Name: %s", getStringField(v, "Name"))
	log.Printf("Adam ID: %s", getStringField(v, "AdamID"))
	log.Printf("Version: %s (ID: %d)", getStringField(v, "Version"), getInt64Field(v, "VersionID"))
	log.Printf("Price: $%.2f", getFloatField(v, "Price"))

	return app
}

func purchaseApp(ctx context.Context, client *goitunes.Client, app interface{}) {
	log.Println("\n=== Purchasing application ===")

	// Use reflection to get AdamID and VersionID
	v := reflect.ValueOf(app)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	adamID := getStringField(v, "AdamID")
	versionID := getInt64Field(v, "VersionID")

	downloadInfo, err := client.Purchase().Buy(ctx, adamID, versionID)
	if err != nil {
		log.Fatalf("Purchase failed: %v", err)
	}

	logPurchaseResult(downloadInfo)
	logDownloadInstructions()
}

func logPurchaseResult(downloadInfo interface{}) {
	log.Println("Purchase successful!")

	// Use reflection to access exported fields
	v := reflect.ValueOf(downloadInfo)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		log.Printf("Download info: %+v", downloadInfo)

		return
	}

	log.Printf("Bundle ID: %s", getStringField(v, "BundleID"))
	log.Printf("Download URL: %s", getStringField(v, "URL"))
	log.Printf("Download Key: %s", getStringField(v, "DownloadKey"))
	log.Printf("Download ID: %s", getStringField(v, "DownloadID"))
	log.Printf("Version ID: %d", getInt64Field(v, "VersionID"))
	log.Printf("File Size: %.2f MB", float64(getInt64Field(v, "FileSize"))/bytesPerMB)

	log.Println("\nDownload Headers:")

	headersField := v.FieldByName("Headers")
	if headersField.IsValid() && headersField.Kind() == reflect.Map {
		for _, key := range headersField.MapKeys() {
			value := headersField.MapIndex(key)
			log.Printf("  %s: %s", key.String(), value.String())
		}
	}

	sinf := getStringField(v, "Sinf")
	if len(sinf) > maxStringDisplay {
		log.Printf("\nSINF (base64): %s...", sinf[:maxStringDisplay])
	} else {
		log.Printf("\nSINF (base64): %s", sinf)
	}

	metadata := getStringField(v, "Metadata")
	if len(metadata) > maxStringDisplay {
		log.Printf("Metadata (base64): %s...", metadata[:maxStringDisplay])
	} else {
		log.Printf("Metadata (base64): %s", metadata)
	}
}

func getStringField(v reflect.Value, field string) string {
	f := v.FieldByName(field)
	if !f.IsValid() || f.Kind() != reflect.String {
		return ""
	}

	return f.String()
}

func getInt64Field(v reflect.Value, field string) int64 {
	f := v.FieldByName(field)
	if !f.IsValid() || f.Kind() != reflect.Int64 {
		return 0
	}

	return f.Int()
}

func getFloatField(v reflect.Value, field string) float64 {
	f := v.FieldByName(field)
	if !f.IsValid() || f.Kind() != reflect.Float64 {
		return 0.0
	}

	return f.Float()
}

func logDownloadInstructions() {
	log.Println("\n=== Download Instructions ===")
	log.Println("1. Use the download URL with provided headers")
	log.Println("2. Save the IPA file")
	log.Println("3. Inject SINF into the IPA (SC_Info/Sinf.sinf)")
	log.Println("4. Inject Metadata into the IPA (iTunesMetadata.plist)")
	log.Println("5. The application can now be installed")
}
