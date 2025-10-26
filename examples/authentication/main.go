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
	fmt.Printf("Authenticated before login: %t\n", client.IsAuthenticated())

	// Perform login
	fmt.Println("\n=== Logging in ===")
	authResp, err := client.Auth().Login(ctx, password)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	fmt.Printf("Login successful!\n")
	fmt.Printf("Apple ID: %s\n", authResp.AppleID)
	fmt.Printf("DSID: %s\n", authResp.DSID)
	fmt.Printf("Password Token: %s...\n", authResp.PasswordToken[:20])
	fmt.Printf("Authenticated: %t\n", authResp.Authenticated)

	// Check authentication status after login
	fmt.Printf("\nAuthenticated after login: %t\n", client.IsAuthenticated())
	fmt.Printf("Can purchase: %t\n", client.CanPurchase())

	// Now you can use authenticated endpoints
	// For example, you can now purchase applications if you have kbsync set up

	fmt.Println("\n=== Client Info ===")
	fmt.Printf("Region: %s\n", client.Region())
	fmt.Printf("Supported regions: %v\n", client.SupportedRegions())
}

