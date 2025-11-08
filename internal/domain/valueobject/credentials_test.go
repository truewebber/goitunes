package valueobject_test

import (
	"testing"

	"github.com/truewebber/goitunes/v2/internal/domain/valueobject"
)

const (
	testAppleID = "test@example.com"
)

func TestNewCredentials(t *testing.T) {
	t.Parallel()

	appleID := testAppleID

	creds, err := valueobject.NewCredentials(appleID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if creds.AppleID() != appleID {
		t.Errorf("Expected appleID %s, got %s", appleID, creds.AppleID())
	}

	if creds.IsAuthenticated() {
		t.Error("New credentials should not be authenticated")
	}
}

func TestNewCredentials_EmptyAppleID(t *testing.T) {
	t.Parallel()

	_, err := valueobject.NewCredentials("")
	if err == nil {
		t.Error("Expected error for empty appleID")
	}
}

func TestNewCredentialsWithTokens(t *testing.T) {
	t.Parallel()

	appleID := testAppleID
	token := "test_token"
	dsid := "123456"

	creds, err := valueobject.NewCredentialsWithTokens(appleID, token, dsid)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if creds.AppleID() != appleID {
		t.Errorf("Expected appleID %s, got %s", appleID, creds.AppleID())
	}

	if creds.PasswordToken() != token {
		t.Errorf("Expected token %s, got %s", token, creds.PasswordToken())
	}

	if creds.DSID() != dsid {
		t.Errorf("Expected DSID %s, got %s", dsid, creds.DSID())
	}

	if !creds.IsAuthenticated() {
		t.Error("Credentials with tokens should be authenticated")
	}
}

func TestNewCredentialsWithTokens_Validation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		appleID     string
		token       string
		dsid        string
		expectError bool
	}{
		{"Valid", testAppleID, "token", "123", false},
		{"Empty appleID", "", "token", "123", true},
		{"Empty token", testAppleID, "", "123", true},
		{"Empty DSID", testAppleID, "token", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := valueobject.NewCredentialsWithTokens(tt.appleID, tt.token, tt.dsid)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCredentials_CanPurchase(t *testing.T) {
	t.Parallel()

	// Without kbsync
	creds, err := valueobject.NewCredentialsWithTokens(testAppleID, "token", "123")
	if err != nil {
		t.Fatalf("Failed to create credentials: %v", err)
	}

	if creds.CanPurchase() {
		t.Error("Credentials without kbsync should not be able to purchase")
	}

	// With kbsync
	creds.SetKbsync("test_kbsync")

	if !creds.CanPurchase() {
		t.Error("Credentials with kbsync should be able to purchase")
	}

	// Not authenticated
	creds2, err := valueobject.NewCredentials(testAppleID)
	if err != nil {
		t.Fatalf("Failed to create credentials: %v", err)
	}

	creds2.SetKbsync("test_kbsync")

	if creds2.CanPurchase() {
		t.Error("Unauthenticated credentials should not be able to purchase")
	}
}

func TestCredentials_Equals(t *testing.T) {
	t.Parallel()

	creds1, err := valueobject.NewCredentialsWithTokens(testAppleID, "token", "123")
	if err != nil {
		t.Fatalf("Failed to create credentials1: %v", err)
	}

	creds2, err := valueobject.NewCredentialsWithTokens(testAppleID, "token", "123")
	if err != nil {
		t.Fatalf("Failed to create credentials2: %v", err)
	}

	creds3, err := valueobject.NewCredentialsWithTokens("other@example.com", "token", "123")
	if err != nil {
		t.Fatalf("Failed to create credentials3: %v", err)
	}

	if !creds1.Equals(creds2) {
		t.Error("Same credentials should be equal")
	}

	if creds1.Equals(creds3) {
		t.Error("Different credentials should not be equal")
	}

	if creds1.Equals(nil) {
		t.Error("Credentials should not equal nil")
	}
}
