package main

import (
	"context"
	"log"
	"os"

	"github.com/truewebber/goitunes/v2/pkg/goitunes"
)

const (
	bytesPerMB = 1024 * 1024
)

func main() {
	// Get credentials from environment variables
	appleID := os.Getenv("APPLE_ID")
	password := os.Getenv("APPLE_PASSWORD")
	kbsync := os.Getenv("KBSYNC")
	guid := os.Getenv("DEVICE_GUID")
	machineName := os.Getenv("MACHINE_NAME")

	if appleID == "" || password == "" || kbsync == "" {
		log.Fatal("APPLE_ID, APPLE_PASSWORD, and KBSYNC environment variables are required")
	}

	if guid == "" {
		guid = "00000000-0000-0000-0000-000000000000"
	}

	if machineName == "" {
		machineName = "MyMachine"
	}

	// Create a client with credentials, kbsync, and device info
	client, err := goitunes.New("us",
		goitunes.WithCredentials(appleID, "", ""),
		goitunes.WithKbsync(kbsync),
		goitunes.WithDevice(guid, machineName, goitunes.UserAgentWindows),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Perform login
	log.Println("=== Logging in ===")

	authResp, err := client.Auth().Login(ctx, password)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	log.Printf("Login successful! DSID: %s", authResp.DSID)

	// Check if we can purchase
	if !client.CanPurchase() {
		log.Fatal("Cannot purchase: missing credentials or kbsync")
	}

	// Get application info first
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
	log.Printf("Name: %s", app.Name)
	log.Printf("Adam ID: %s", app.AdamID)
	log.Printf("Version: %s (ID: %d)", app.Version, app.VersionID)
	log.Printf("Price: $%.2f", app.Price)

	// Purchase the application
	log.Println("\n=== Purchasing application ===")

	downloadInfo, err := client.Purchase().Buy(ctx, app.AdamID, app.VersionID)
	if err != nil {
		log.Fatalf("Purchase failed: %v", err)
	}

	log.Println("Purchase successful!")
	log.Printf("Bundle ID: %s", downloadInfo.BundleID)
	log.Printf("Download URL: %s", downloadInfo.URL)
	log.Printf("Download Key: %s", downloadInfo.DownloadKey)
	log.Printf("Download ID: %s", downloadInfo.DownloadID)
	log.Printf("Version ID: %d", downloadInfo.VersionID)
	log.Printf("File Size: %.2f MB", float64(downloadInfo.FileSize)/bytesPerMB)

	// Display download headers
	log.Println("\nDownload Headers:")

	for key, value := range downloadInfo.Headers {
		log.Printf("  %s: %s", key, value)
	}

	// SINF and Metadata are base64 encoded
	log.Printf("\nSINF (base64): %s...", downloadInfo.Sinf[:50])
	log.Printf("Metadata (base64): %s...", downloadInfo.Metadata[:50])

	log.Println("\n=== Download Instructions ===")
	log.Println("1. Use the download URL with provided headers")
	log.Println("2. Save the IPA file")
	log.Println("3. Inject SINF into the IPA (SC_Info/Sinf.sinf)")
	log.Println("4. Inject Metadata into the IPA (iTunesMetadata.plist)")
	log.Println("5. The application can now be installed")
}
