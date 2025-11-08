package main

import (
	"context"
	"log"
	"os"

	"github.com/truewebber/goitunes/v2/pkg/goitunes"
)

func main() {
	// Get credentials from environment variables
	appleID := os.Getenv("APPLE_ID")
	password := os.Getenv("APPLE_PASSWORD")
	guid := os.Getenv("DEVICE_GUID")
	machineName := os.Getenv("MACHINE_NAME")

	if appleID == "" || password == "" {
		log.Fatal("APPLE_ID and APPLE_PASSWORD environment variables are required")
	}

	// Set default device info if not provided
	if guid == "" {
		guid = "00000000-0000-0000-0000-000000000000"
	}

	if machineName == "" {
		machineName = "MyMachine"
	}

	// Create a client with credentials and device info
	client, err := goitunes.New("us",
		goitunes.WithCredentials(appleID, "", ""), // Empty tokens will be filled after login
		goitunes.WithDevice(guid, machineName, goitunes.UserAgentWindows),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Check authentication status before login
	log.Printf("Authenticated before login: %t", client.IsAuthenticated())

	// Perform login
	log.Println("\n=== Logging in ===")

	authResp, err := client.Auth().Login(ctx, password)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	log.Println("Login successful!")
	log.Printf("Apple ID: %s", authResp.AppleID)
	log.Printf("DSID: %s", authResp.DSID)
	log.Printf("Password Token: %s...", authResp.PasswordToken[:20])
	log.Printf("Authenticated: %t", authResp.Authenticated)

	// Check authentication status after login
	log.Printf("\nAuthenticated after login: %t", client.IsAuthenticated())
	log.Printf("Can purchase: %t", client.CanPurchase())

	// Now you can use authenticated endpoints
	// For example, you can now purchase applications if you have kbsync set up

	log.Println("\n=== Client Info ===")
	log.Printf("Region: %s", client.Region())
	log.Printf("Supported regions: %v", client.SupportedRegions())
}
