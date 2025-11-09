package main

import (
	"context"
	"log"
	"os"

	"github.com/truewebber/goitunes/v2/pkg/goitunes"
)

const (
	passwordTokenMinLen = 20
)

func main() {
	appleID, password, guid, machineName := getEnvVars()
	client := createClient(appleID, guid, machineName)
	ctx := context.Background()

	logAuthStatus(client, "before")
	performLogin(ctx, client, password)
	logAuthStatus(client, "after")
	logClientInfo(client)
}

//nolint:nonamedreturns // gocritic requires named returns
func getEnvVars() (appleID, password, guid, machineName string) {
	appleID = os.Getenv("APPLE_ID")
	password = os.Getenv("APPLE_PASSWORD")
	guid = os.Getenv("DEVICE_GUID")
	machineName = os.Getenv("MACHINE_NAME")

	if appleID == "" || password == "" {
		log.Fatal("APPLE_ID and APPLE_PASSWORD environment variables are required")
	}

	if guid == "" {
		guid = "00000000-0000-0000-0000-000000000000"
	}

	if machineName == "" {
		machineName = "MyMachine"
	}

	return appleID, password, guid, machineName
}

func createClient(appleID, guid, machineName string) *goitunes.Client {
	client, err := goitunes.New("us",
		goitunes.WithAppleID(appleID), // Tokens will be filled after login
		goitunes.WithDevice(guid, machineName, goitunes.UserAgentWindows),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return client
}

func performLogin(ctx context.Context, client *goitunes.Client, password string) {
	log.Println("\n=== Logging in ===")

	authResp, err := client.Auth().Login(ctx, password)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	logAuthResult(authResp)

	// Access fields directly (they are exported from dto.AuthenticateResponse)
	log.Printf("Apple ID: %s", authResp.AppleID)
	log.Printf("DSID: %s", authResp.DSID)

	if len(authResp.PasswordToken) > passwordTokenMinLen {
		log.Printf("Password Token: %s...", authResp.PasswordToken[:passwordTokenMinLen])
	} else {
		log.Printf("Password Token: %s", authResp.PasswordToken)
	}

	log.Printf("Authenticated: %t", authResp.Authenticated)
}

func logAuthResult(authResp interface{}) {
	log.Println("Login successful!")

	// Use reflection or type assertion - for now just log that we got response
	// The actual fields are accessed via the returned value
	log.Printf("Auth response received (type: %T)", authResp)
}

func logAuthStatus(client *goitunes.Client, when string) {
	log.Printf("Authenticated %s login: %t", when, client.IsAuthenticated())

	if when == "after" {
		log.Printf("Can purchase: %t", client.CanPurchase())
	}
}

func logClientInfo(client *goitunes.Client) {
	log.Println("\n=== Client Info ===")
	log.Printf("Region: %s", client.Region())
	log.Printf("Supported regions: %v", client.SupportedRegions())
}
