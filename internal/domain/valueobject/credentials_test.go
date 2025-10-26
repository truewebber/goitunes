package valueobject

import "testing"

func TestNewCredentials(t *testing.T) {
	appleID := "test@example.com"

	creds, err := NewCredentials(appleID)
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
	_, err := NewCredentials("")
	if err == nil {
		t.Error("Expected error for empty appleID")
	}
}

func TestNewCredentialsWithTokens(t *testing.T) {
	appleID := "test@example.com"
	token := "test_token"
	dsid := "123456"

	creds, err := NewCredentialsWithTokens(appleID, token, dsid)
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
	tests := []struct {
		name        string
		appleID     string
		token       string
		dsid        string
		expectError bool
	}{
		{"Valid", "test@example.com", "token", "123", false},
		{"Empty appleID", "", "token", "123", true},
		{"Empty token", "test@example.com", "", "123", true},
		{"Empty DSID", "test@example.com", "token", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCredentialsWithTokens(tt.appleID, tt.token, tt.dsid)
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
	// Without kbsync
	creds, _ := NewCredentialsWithTokens("test@example.com", "token", "123")
	if creds.CanPurchase() {
		t.Error("Credentials without kbsync should not be able to purchase")
	}

	// With kbsync
	creds.SetKbsync("test_kbsync")
	if !creds.CanPurchase() {
		t.Error("Credentials with kbsync should be able to purchase")
	}

	// Not authenticated
	creds2, _ := NewCredentials("test@example.com")
	creds2.SetKbsync("test_kbsync")
	if creds2.CanPurchase() {
		t.Error("Unauthenticated credentials should not be able to purchase")
	}
}

func TestCredentials_Equals(t *testing.T) {
	creds1, _ := NewCredentialsWithTokens("test@example.com", "token", "123")
	creds2, _ := NewCredentialsWithTokens("test@example.com", "token", "123")
	creds3, _ := NewCredentialsWithTokens("other@example.com", "token", "123")

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
