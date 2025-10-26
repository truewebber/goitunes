package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/truewebber/goitunes/pkg/goitunes"
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
	fmt.Println("=== Logging in ===")
	authResp, err := client.Auth().Login(ctx, password)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	fmt.Printf("Login successful! DSID: %s\n", authResp.DSID)

	// Check if we can purchase
	if !client.CanPurchase() {
		log.Fatal("Cannot purchase: missing credentials or kbsync")
	}

	// Get application info first
	bundleID := "com.example.app" // Replace with actual bundle ID
	fmt.Printf("\n=== Getting application info for %s ===\n", bundleID)
	
	apps, err := client.Applications().GetByBundleID(ctx, bundleID)
	if err != nil {
		log.Fatalf("Failed to get app info: %v", err)
	}

	if len(apps) == 0 {
		log.Fatal("Application not found")
	}

	app := apps[0]
	fmt.Printf("Name: %s\n", app.Name)
	fmt.Printf("Adam ID: %s\n", app.AdamID)
	fmt.Printf("Version: %s (ID: %d)\n", app.Version, app.VersionID)
	fmt.Printf("Price: $%.2f\n", app.Price)

	// Purchase the application
	fmt.Println("\n=== Purchasing application ===")
	downloadInfo, err := client.Purchase().Buy(ctx, app.AdamID, app.VersionID)
	if err != nil {
		log.Fatalf("Purchase failed: %v", err)
	}

	fmt.Println("Purchase successful!")
	fmt.Printf("Bundle ID: %s\n", downloadInfo.BundleID)
	fmt.Printf("Download URL: %s\n", downloadInfo.URL)
	fmt.Printf("Download Key: %s\n", downloadInfo.DownloadKey)
	fmt.Printf("Download ID: %s\n", downloadInfo.DownloadID)
	fmt.Printf("Version ID: %d\n", downloadInfo.VersionID)
	fmt.Printf("File Size: %.2f MB\n", float64(downloadInfo.FileSize)/(1024*1024))

	// Display download headers
	fmt.Println("\nDownload Headers:")
	for key, value := range downloadInfo.Headers {
		fmt.Printf("  %s: %s\n", key, value)
	}

	// SINF and Metadata are base64 encoded
	fmt.Printf("\nSINF (base64): %s...\n", downloadInfo.Sinf[:50])
	fmt.Printf("Metadata (base64): %s...\n", downloadInfo.Metadata[:50])

	fmt.Println("\n=== Download Instructions ===")
	fmt.Println("1. Use the download URL with provided headers")
	fmt.Println("2. Save the IPA file")
	fmt.Println("3. Inject SINF into the IPA (SC_Info/Sinf.sinf)")
	fmt.Println("4. Inject Metadata into the IPA (iTunesMetadata.plist)")
	fmt.Println("5. The application can now be installed")
}

